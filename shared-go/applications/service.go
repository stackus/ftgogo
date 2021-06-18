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
	"github.com/stackus/edat-kafka-go"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat-pgx"
	"github.com/stackus/edat-stan"
	"golang.org/x/sync/errgroup"

	"github.com/stackus/edat/es"
	"github.com/stackus/edat/inmem"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/saga"

	"shared-go/egress"
	"shared-go/instrumentation"
	"shared-go/logging"
	"shared-go/logging/zerologto"
	"shared-go/rpc"
	"shared-go/web"
)

type svcConfig struct {
	Environment     string        `envconfig:"ENVIRONMENT" default:"production"`
	ServiceID       string        `envconfig:"SERVICE_ID" required:"true"`
	LogLevel        logging.Level `envconfig:"LOG_LEVEL" default:"WARN" desc:"options: [TRACE,DEBUG,INFO,WARN,ERROR,PANIC]"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" desc:"time to allow services to gracefully stop"`
	Web             webCfg        `envconfig:"WEB"`                                                             // Web Config
	Rpc             rpc.ServerCfg `envconfig:"RPC"`                                                             // RPC Config
	Postgres        pgCfg         `envconfig:"PG"`                                                              // DataDriver / Postgres
	EventDriver     string        `envconfig:"EVENT_DRIVER" default:"inmem" desc:"options: [inmem,nats,kafka]"` // "inmem", "nats", "kafka"
	Nats            natsCfg       `envconfig:"NATS"`                                                            // EventDriver / Nats Streaming Config
	Kafka           kafkaCfg      `envconfig:"KAFKA"`                                                           // EventDriver / Kafka Config
}

type Service struct {
	appFn             func(*Service) error
	app               *cobra.Command
	Cfg               *svcConfig
	Logger            zerolog.Logger
	PgConn            edatpgx.Client
	AggregateStore    es.AggregateRootStore
	SagaInstanceStore saga.InstanceStore
	Publisher         *msg.Publisher
	Subscriber        *msg.Subscriber
	WebServer         web.Server
	RpcServer         rpc.Server
}

func NewService(appFn func(*Service) error) *Service {
	svc := &Service{
		appFn: appFn,
		Cfg:   &svcConfig{},
	}

	svc.app = &cobra.Command{
		Use:           "service",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          svc.run,
	}

	svc.app.Flags().StringSliceVarP(&envFiles, "envfiles", "f", []string{}, "environment variable override files")

	appUsage := svc.app.UsageString()

	svc.app.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println(appUsage)
		fmt.Println("Environment Variables:")
		return envconfig.Usage("", svc.Cfg)
	})

	return svc
}

func (s Service) Execute() error {
	cobra.OnInitialize(s.initConfig)

	return s.app.Execute()
}

func (s *Service) initConfig() {
	var err error

	appPrefix := os.Getenv("APP_PREFIX")

	if len(envFiles) >= 1 {
		err = godotenv.Load(envFiles...)
		if err != nil {
			fmt.Println("error reading environment variable overrides", err, "files:", envFiles)
			os.Exit(1)
		}
	}

	err = envconfig.Process(appPrefix, s.Cfg)
	if err != nil {
		fmt.Println("error reading environment variables", err)
		os.Exit(1)
	}
}

func (s *Service) run(*cobra.Command, []string) error {
	var err error

	s.Logger, err = logging.NewZeroLogger(logging.Config{
		Environment: s.Cfg.Environment,
		LogLevel:    s.Cfg.LogLevel,
	})
	if err != nil {
		return err
	}

	log.DefaultLogger = zerologto.Logger(s.Logger)

	var pgConn *pgxpool.Pool
	pgConn, err = pgxpool.Connect(context.Background(), s.Cfg.Postgres.Conn)
	if err != nil {
		panic(err)
	}

	// 1. Outbox: Use session client which will fetch a transaction from the context
	s.PgConn = edatpgx.NewSessionClient()

	defer func() {
		if pgConn != nil {
			pgConn.Close()
		}
	}()

	s.AggregateStore = edatpgx.NewSnapshotStore(s.PgConn)(edatpgx.NewEventStore(s.PgConn))
	s.SagaInstanceStore = edatpgx.NewSagaInstanceStore(s.PgConn)

	var msgConsumer msg.Consumer
	var msgProducer msg.Producer

	switch {
	case s.Cfg.EventDriver == "nats":
		var conn stan.Conn
		conn, err = stan.Connect(s.Cfg.Nats.ClusterID, s.Cfg.ServiceID, stan.NatsURL(s.Cfg.Nats.URL))
		if err != nil {
			panic(err)
		}
		msgConsumer = edatstan.NewConsumer(conn, s.Cfg.ServiceID,
			edatstan.WithConsumerActWait(s.Cfg.Nats.AckWaitTimeout),
		)
		msgProducer = edatstan.NewProducer(conn)
	case s.Cfg.EventDriver == "kafka":
		msgConsumer = edatkafkago.NewConsumer(s.Cfg.Kafka.Brokers, s.Cfg.ServiceID)
		msgProducer = edatkafkago.NewProducer(s.Cfg.Kafka.Brokers)
	default:
		msgConsumer = inmem.NewConsumer()
		msgProducer = inmem.NewProducer()
	}

	// 2. Outbox: Producer publishes into the db
	msgProducer = edatpgx.NewMessageStore(s.PgConn)
	s.Publisher = msg.NewPublisher(msgProducer)

	s.Subscriber = msg.NewSubscriber(msgConsumer)
	s.Subscriber.Use(
		instrumentation.MessageInstrumentation(),
		// 3. Outbox: Use a message receiver middleware to start a new transaction for each incoming message
		edatpgx.ReceiverSessionMiddleware(pgConn, zerologto.Logger(s.Logger)),
	)

	s.RpcServer = initRpcServer(s.Cfg.Rpc, pgConn, s.Logger)
	s.WebServer = initWebServer(s.Cfg.Web, pgConn, s.Logger)

	err = s.appFn(s)
	if err != nil {
		return err
	}

	waiter := egress.NewWaiter()

	waiter.Add(s.waitForRpcServer)
	waiter.Add(s.waitForWebServer)
	waiter.Add(s.waitForMessaging)

	s.Logger.Debug().Msg("service starting")

	return waiter.Wait()
}

func (s Service) waitForMessaging(ctx context.Context) error {
	defer s.Logger.Debug().Msg("messaging has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer s.Logger.Debug().Msg("messaging has exited")
		return s.Subscriber.Start(ctx)
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.Logger.Debug().Msg("shutting down messaging")
		sCtx, cancel := context.WithTimeout(context.Background(), s.Cfg.ShutdownTimeout)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			if s.Publisher != nil {
				if err := s.Publisher.Stop(sCtx); err != nil {
					s.Logger.Error().Err(err).Msg("error while shutting down publisher")
				}
			}
		}()
		go func() {
			defer wg.Done()
			if s.Subscriber != nil {
				if err := s.Subscriber.Stop(sCtx); err != nil {
					s.Logger.Error().Err(err).Msg("error while shutting down subscriber")
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
			s.Logger.Warn().Msg("timed out while shutting down messaging")
		}
		return nil
	})

	return group.Wait()
}

func (s Service) waitForWebServer(ctx context.Context) (err error) {
	defer s.Logger.Debug().Msg("web server has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer s.Logger.Debug().Msg("web server exited")
		err = s.WebServer.Start()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.Logger.Debug().Msg("shutting down the web server")
		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.ShutdownTimeout)
		defer cancel()

		if err = s.WebServer.Shutdown(ctx); err != nil {
			s.Logger.Warn().Msg("timed out while shutting down the web server")
		}

		return nil
	})

	return group.Wait()
}

func (s Service) waitForRpcServer(ctx context.Context) (err error) {
	defer s.Logger.Debug().Msg("rpc server has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer s.Logger.Debug().Msg("rpc server exited")
		err = s.RpcServer.Start()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.Logger.Debug().Msg("shutting down the rpc server")
		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.ShutdownTimeout)
		defer cancel()

		if err = s.RpcServer.Shutdown(ctx); err != nil {
			s.Logger.Warn().Msg("timed out while shutting down the rpc server")
		}

		return nil
	})

	return group.Wait()
}
