package applications

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	edatkafkago "github.com/stackus/edat-kafka-go"
	edatpgx "github.com/stackus/edat-pgx"
	edatstan "github.com/stackus/edat-stan"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/stackus/edat/inmem"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"

	"shared-go/egress"
	"shared-go/instrumentation"
	"shared-go/logging"
	"shared-go/logging/zerologto"
	"shared-go/rpc"
	"shared-go/web"
)

type monoConfig struct {
	Environment     string        `envconfig:"ENVIRONMENT" default:"production"`
	ServiceID       string        `envconfig:"SERVICE_ID" required:"true"`
	LogLevel        logging.Level `envconfig:"LOG_LEVEL" default:"WARN" desc:"options: [TRACE,DEBUG,INFO,WARN,ERROR,PANIC]"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" desc:"time to allow services to gracefully stop"`
	Web             webCfg        `envconfig:"WEB"`                                                             // Web Config
	Rpc             rpcCfg        `envconfig:"RPC"`                                                             // RPC Server & Client Config
	Postgres        pgCfg         `envconfig:"PG"`                                                              // DataDriver / Postgres
	EventDriver     string        `envconfig:"EVENT_DRIVER" default:"inmem" desc:"options: [inmem,nats,kafka]"` // "inmem", "nats", "kafka"
	Nats            natsCfg       `envconfig:"NATS"`                                                            // EventDriver / Nats Streaming Config
	Kafka           kafkaCfg      `envconfig:"KAFKA"`                                                           // EventDriver / Kafka Config
}

type Monolith struct {
	appFn        func(*Monolith) error
	app          *cobra.Command
	Cfg          *monoConfig
	Logger       zerolog.Logger
	PgConn       edatpgx.Client
	CDCPgConn    edatpgx.Client
	Publishers   []*msg.Publisher
	CDCPublisher *msg.Publisher
	Subscriber   *msg.Subscriber
	WebServer    web.Server
	RpcServer    rpc.Server
	Clients      map[string]*grpc.ClientConn
	Processors   []outbox.MessageProcessor
}

func NewMonolith(appFn func(*Monolith) error) *Monolith {
	mono := &Monolith{
		appFn: appFn,
		Cfg:   &monoConfig{},
	}

	mono.app = &cobra.Command{
		Use:           "monolith",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          mono.run,
	}

	mono.app.Flags().StringSliceVarP(&envFiles, "envfiles", "f", []string{}, "environment variable override files")

	appUsage := mono.app.UsageString()

	mono.app.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println(appUsage)
		fmt.Println("Environment Variables:")
		return envconfig.Usage("", mono.Cfg)
	})

	return mono
}

func (m Monolith) Execute() error {
	cobra.OnInitialize(m.initConfig)

	return m.app.Execute()
}

func (m *Monolith) initConfig() {
	var err error

	appPrefix := os.Getenv("APP_PREFIX")

	if len(envFiles) >= 1 {
		err = godotenv.Load(envFiles...)
		if err != nil {
			fmt.Println("error reading environment variable overrides", err, "files:", envFiles)
			os.Exit(1)
		}
	}

	err = envconfig.Process(appPrefix, m.Cfg)
	if err != nil {
		fmt.Println("error reading environment variables", err)
		os.Exit(1)
	}
}

func (m *Monolith) run(*cobra.Command, []string) error {
	var err error

	m.Logger, err = logging.NewZeroLogger(logging.Config{
		Environment: m.Cfg.Environment,
		LogLevel:    m.Cfg.LogLevel,
	})
	if err != nil {
		return err
	}

	log.DefaultLogger = zerologto.Logger(m.Logger)

	var poolConn *pgxpool.Pool
	poolConn, err = pgxpool.Connect(context.Background(), m.Cfg.Postgres.Conn)
	if err != nil {
		panic(err)
	}

	m.CDCPgConn = poolConn
	m.PgConn = edatpgx.NewSessionClient()

	defer func() {
		if poolConn != nil {
			poolConn.Close()
		}
	}()

	var msgConsumer msg.Consumer
	var msgProducer msg.Producer

	switch {
	case m.Cfg.EventDriver == "nats":
		var conn stan.Conn
		conn, err = stan.Connect(m.Cfg.Nats.ClusterID, m.Cfg.ServiceID, stan.NatsURL(m.Cfg.Nats.URL))
		if err != nil {
			panic(err)
		}
		msgConsumer = edatstan.NewConsumer(conn, m.Cfg.ServiceID,
			edatstan.WithConsumerActWait(m.Cfg.Nats.AckWaitTimeout),
		)
		msgProducer = edatstan.NewProducer(conn)
	case m.Cfg.EventDriver == "kafka":
		msgConsumer = edatkafkago.NewConsumer(m.Cfg.Kafka.Brokers, m.Cfg.ServiceID)
		msgProducer = edatkafkago.NewProducer(m.Cfg.Kafka.Brokers)
	default:
		msgConsumer = inmem.NewConsumer()
		msgProducer = inmem.NewProducer()
	}

	m.CDCPublisher = msg.NewPublisher(msgProducer)

	m.Subscriber = msg.NewSubscriber(msgConsumer)
	m.Subscriber.Use(
		instrumentation.MessageInstrumentation(),
		edatpgx.ReceiverSessionMiddleware(poolConn, zerologto.Logger(m.Logger)),
	)

	m.Clients, err = initRpcClients(m.Cfg.Rpc, m.Logger)
	if err != nil {
		return err
	}

	m.RpcServer = initRpcServer(m.Cfg.Rpc.Server, poolConn, m.Logger)

	m.WebServer = initWebServer(m.Cfg.Web, poolConn, m.Logger)

	err = m.appFn(m)
	if err != nil {
		return err
	}

	waiter := egress.NewWaiter()

	waiter.Add(m.waitForRpcServer)
	waiter.Add(m.waitForWebServer)
	waiter.Add(m.waitForMessaging)
	waiter.Add(m.waitForProcessors)
	waiter.Add(m.waitForConnections)

	m.Logger.Debug().Msg("service starting")

	return waiter.Wait()
}

func (m Monolith) waitForMessaging(ctx context.Context) error {
	defer m.Logger.Debug().Msg("messaging has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer m.Logger.Debug().Msg("messaging has exited")
		return m.Subscriber.Start(ctx)
	})

	group.Go(func() error {
		<-gCtx.Done()
		m.Logger.Debug().Msg("shutting down messaging")
		sCtx, cancel := context.WithTimeout(context.Background(), m.Cfg.ShutdownTimeout)
		defer cancel()

		wg := sync.WaitGroup{}
		for _, p := range m.Publishers {
			publisher := p
			wg.Add(1)
			go func() {
				defer wg.Done()
				if publisher != nil {
					if err := publisher.Stop(sCtx); err != nil {
						m.Logger.Error().Err(err).Msg("error while shutting down publisher")
					}
				}
			}()
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if m.Subscriber != nil {
				if err := m.Subscriber.Stop(sCtx); err != nil {
					m.Logger.Error().Err(err).Msg("error while shutting down subscriber")
				}
			}
		}()
		done := make(chan struct{})
		go func() {
			defer close(done)
			wg.Wait()
		}()
		select {
		case <-done:
		case <-sCtx.Done():
			m.Logger.Warn().Msg("timed out while shutting down messaging")
		}
		return nil
	})

	return group.Wait()
}

func (m Monolith) waitForWebServer(ctx context.Context) (err error) {
	defer m.Logger.Debug().Msg("web server has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer m.Logger.Debug().Msg("web server exited")
		err = m.WebServer.Start()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		m.Logger.Debug().Msg("shutting down the web server")
		ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.ShutdownTimeout)
		defer cancel()

		if err = m.WebServer.Shutdown(ctx); err != nil {
			m.Logger.Warn().Msg("timed out while shutting down the web server")
		}

		return nil
	})

	return group.Wait()
}

func (m Monolith) waitForRpcServer(ctx context.Context) (err error) {
	defer m.Logger.Debug().Msg("rpc server has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer m.Logger.Debug().Msg("rpc server exited")
		err = m.RpcServer.Start()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		m.Logger.Debug().Msg("shutting down the rpc server")
		ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.ShutdownTimeout)
		defer cancel()

		if err = m.RpcServer.Shutdown(ctx); err != nil {
			m.Logger.Warn().Msg("timed out while shutting down the rpc server")
		}

		return nil
	})

	return group.Wait()
}

func (m Monolith) waitForProcessors(ctx context.Context) (err error) {
	defer m.Logger.Debug().Msg("message processors have been shutdown")

	group, gCtx := errgroup.WithContext(ctx)

	for _, processor := range m.Processors {
		startFn := func(p outbox.MessageProcessor) func() error {
			return func() error {
				err := p.Start(ctx)
				if err != nil {
					m.Logger.Err(err).Msg("processor exited with an error")
				}
				return err
			}
		}(processor)

		stopFn := func(p outbox.MessageProcessor) func() error {
			return func() error {
				<-gCtx.Done()
				m.Logger.Debug().Msg("stopping a message processor")
				ctx, cancel := context.WithTimeout(context.Background(), m.Cfg.ShutdownTimeout)
				defer cancel()

				if err = p.Stop(ctx); err != nil {
					m.Logger.Warn().Msg("timed out while stopping the message processor")
				}

				return nil
			}
		}(processor)

		group.Go(startFn)
		group.Go(stopFn)
	}

	return group.Wait()
}

func (m Monolith) waitForConnections(ctx context.Context) error {
	var err error
	defer m.Logger.Debug().Msg("cleanup has been completed")

	<-ctx.Done()
	wg := sync.WaitGroup{}

	for _, client := range m.Clients {
		wg.Add(1)
		go func(c *grpc.ClientConn) {
			sCtx, cancel := context.WithTimeout(context.Background(), m.Cfg.ShutdownTimeout)
			defer cancel()

			done := make(chan struct{})
			go func() {
				cErr := c.Close()
				close(done)
				if cErr != nil {
					err = cErr
				}
			}()
			select {
			case <-sCtx.Done():
				err = fmt.Errorf("cleanup failed to complete within allowed time")
			case <-done:
			}
			wg.Done()
		}(client)
	}

	wg.Wait()

	return err
}
