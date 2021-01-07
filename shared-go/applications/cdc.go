package applications

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	_ "github.com/stackus/edat-msgpack"
	"github.com/stackus/edat-pgx"
	"github.com/stackus/edat-stan"
	"github.com/stackus/edat/log"
	"github.com/stackus/edat/msg"
	"github.com/stackus/edat/outbox"
	"golang.org/x/sync/errgroup"

	"shared-go/eddsarama"
	"shared-go/egress"
	"shared-go/logging"
	"shared-go/logging/zerologto"
)

type cdcConfig struct {
	Environment     string        `envconfig:"ENVIRONMENT" default:"production"`
	ServiceID       string        `envconfig:"SERVICE_ID" required:"true"`
	LogLevel        logging.Level `envconfig:"LOG_LEVEL" default:"WARN" desc:"options: [TRACE,DEBUG,INFO,WARN,ERROR,PANIC]"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" desc:"time to allow services to gracefully stop"`
	Postgres        pgCfg         `envconfig:"PG"`                                                              // DataDriver / Postgres
	EventDriver     string        `envconfig:"EVENT_DRIVER" default:"inmem" desc:"options: [inmem,nats,kafka]"` // "inmem", "nats", "kafka"
	Nats            natsCfg       `envconfig:"NATS"`                                                            // EventDriver / Nats Streaming Config
	Kafka           kafkaCfg      `envconfig:"KAFKA"`                                                           // EventDriver / Kafka Config
}

type CDC struct {
	appFn        func(*CDC) error
	app          *cobra.Command
	cfg          *cdcConfig
	Logger       zerolog.Logger
	PgConn       edatpgx.Client
	MessageStore outbox.MessageStore
	Publisher    *msg.Publisher
	Processor    outbox.MessageProcessor
}

func NewCDC(appFn func(*CDC) error) *CDC {
	svc := &CDC{
		appFn: appFn,
		cfg:   &cdcConfig{},
	}

	svc.app = &cobra.Command{
		Use:           "cdc",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          svc.run,
	}

	svc.app.Flags().StringSliceVarP(&envFiles, "envfiles", "f", []string{}, "environment variable override files")

	appUsage := svc.app.UsageString()

	svc.app.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println(appUsage)
		fmt.Println("Environment Variables:")
		return envconfig.Usage("", svc.cfg)
	})

	return svc
}

func (s CDC) Execute() error {
	cobra.OnInitialize(s.initConfig)

	return s.app.Execute()
}

func (s *CDC) initConfig() {
	var err error

	appPrefix := os.Getenv("APP_PREFIX")

	if len(envFiles) >= 1 {
		err = godotenv.Load(envFiles...)
		if err != nil {
			fmt.Println("error reading environment variable overrides", err, "files:", envFiles)
			os.Exit(1)
		}
	}

	err = envconfig.Process(appPrefix, s.cfg)
	if err != nil {
		fmt.Println("error reading environment variables", err)
		os.Exit(1)
	}
}

func (s *CDC) run(*cobra.Command, []string) error {
	var err error

	s.Logger, err = logging.NewZeroLogger(logging.Config{
		Environment: s.cfg.Environment,
		LogLevel:    s.cfg.LogLevel,
	})
	if err != nil {
		return err
	}

	log.DefaultLogger = zerologto.Logger(s.Logger)

	var pgConn *pgxpool.Pool
	pgConn, err = pgxpool.Connect(context.Background(), s.cfg.Postgres.Conn)
	if err != nil {
		panic(err)
	}

	s.PgConn = pgConn

	s.MessageStore = edatpgx.NewMessageStore(s.PgConn)

	defer func() {
		if pgConn != nil {
			pgConn.Close()
		}
	}()

	var msgProducer msg.Producer

	switch {
	case s.cfg.EventDriver == "nats":
		var conn stan.Conn
		conn, err = stan.Connect(s.cfg.Nats.ClusterID, s.cfg.ServiceID, stan.NatsURL(s.cfg.Nats.URL))
		if err != nil {
			panic(err)
		}
		msgProducer = edatstan.NewProducer(conn)
	case s.cfg.EventDriver == "kafka":
		msgProducer, err = eddsarama.NewProducer(s.cfg.Kafka.Brokers, s.cfg.ServiceID)
		if err != nil {
			panic(err)
		}
	default:
		panic("cdc services cannot be started with an inmem destination")
	}

	s.Publisher = msg.NewPublisher(msgProducer)

	s.Processor = outbox.NewPollingProcessor(s.MessageStore, s.Publisher)

	err = s.appFn(s)
	if err != nil {
		return err
	}

	waiter := egress.NewWaiter()

	waiter.Add(s.waitForMessaging)
	waiter.Add(s.waitForProcessor)

	s.Logger.Debug().Msg("cdc starting")

	return waiter.Wait()
}

func (s CDC) waitForMessaging(ctx context.Context) error {
	defer s.Logger.Debug().Msg("messaging has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		<-gCtx.Done()
		s.Logger.Debug().Msg("shutting down messaging")
		sCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
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

func (s CDC) waitForProcessor(ctx context.Context) (err error) {
	defer s.Logger.Debug().Msg("message processor has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer s.Logger.Debug().Msg("message processor exited")
		return s.Processor.Start(ctx)
	})

	group.Go(func() error {
		<-gCtx.Done()
		s.Logger.Debug().Msg("stopping the message processor")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err = s.Processor.Stop(ctx); err != nil {
			s.Logger.Warn().Msg("timed out while stopping the message processor")
		}

		return nil
	})

	return group.Wait()
}
