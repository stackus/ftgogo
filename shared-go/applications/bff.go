package applications

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	_ "github.com/stackus/edat-msgpack"
	http2 "github.com/stackus/edat/http"
	"github.com/stackus/edat/log"
	"golang.org/x/sync/errgroup"

	"shared-go/egress"
	"shared-go/instrumentation"
	"shared-go/logging"
	"shared-go/logging/zerologto"
	"shared-go/rpc"
	"shared-go/web"
)

type bffConfig struct {
	Environment     string        `envconfig:"ENVIRONMENT" default:"production"`
	ServiceID       string        `envconfig:"SERVICE_ID" required:"true"`
	LogLevel        logging.Level `envconfig:"LOG_LEVEL" default:"WARN" desc:"options: [TRACE,DEBUG,INFO,WARN,ERROR,PANIC]"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" desc:"time to allow services to gracefully stop"`
	Web             webCfg        `envconfig:"WEB"` // Web Config
	Rpc             rpc.ClientCfg `envconfig:"RPC"` // RPC Config
}

type BFF struct {
	appFn     func(*BFF) error
	app       *cobra.Command
	Cfg       *bffConfig
	Logger    zerolog.Logger
	WebServer web.Server
	Cleanup   CleanupFunc
}

func NewBFF(appFn func(*BFF) error) *BFF {
	bff := &BFF{
		appFn: appFn,
		Cfg:   &bffConfig{},
	}

	bff.app = &cobra.Command{
		Use:           "service",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          bff.run,
	}

	bff.app.Flags().StringSliceVarP(&envFiles, "envfiles", "f", []string{}, "environment variable override files")

	appUsage := bff.app.UsageString()

	bff.app.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println(appUsage)
		fmt.Println("Environment Variables:")
		return envconfig.Usage("", bff.Cfg)
	})

	return bff
}

func (s BFF) Execute() error {
	cobra.OnInitialize(s.initConfig)

	return s.app.Execute()
}

func (s *BFF) initConfig() {
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

func (s *BFF) run(*cobra.Command, []string) error {
	var err error

	s.Logger, err = logging.NewZeroLogger(logging.Config{
		Environment: s.Cfg.Environment,
		LogLevel:    s.Cfg.LogLevel,
	})
	if err != nil {
		return err
	}

	log.DefaultLogger = zerologto.Logger(s.Logger)

	s.WebServer = web.NewServer(s.Cfg.Web.Http, web.WithHealthCheck(s.Cfg.Web.PingPath))

	s.WebServer.Options(s.Cfg.Web.ApiPath,
		web.WithSecure(),
		web.WithCors(s.Cfg.Web.Cors),
		web.WithMiddleware(
			instrumentation.WebInstrumentation(),
			web.ZeroLogger(s.Logger),
			http2.RequestContext,
		))

	s.WebServer.Mount(s.Cfg.Web.MetricsPath, func(r chi.Router) http.Handler {
		return promhttp.Handler()
	})

	err = s.appFn(s)
	if err != nil {
		return err
	}

	waiter := egress.NewWaiter()

	waiter.Add(s.waitForWebServer)
	waiter.Add(s.waitForCleanup)

	s.Logger.Debug().Msg("backend-for-frontend starting")

	return waiter.Wait()
}

func (s BFF) waitForWebServer(ctx context.Context) (err error) {
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

func (s BFF) waitForCleanup(ctx context.Context) (err error) {
	defer s.Logger.Debug().Msg("cleanup has been completed")

	<-ctx.Done()

	if s.Cleanup == nil {
		return nil
	}

	sCtx, cancel := context.WithTimeout(context.Background(), s.Cfg.ShutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		err = s.Cleanup(sCtx)
		close(done)
	}()

	select {
	case <-sCtx.Done():
		return fmt.Errorf("cleanup failed to complete within allowed time")
	case <-done:
		return
	}
}
