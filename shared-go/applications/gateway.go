package applications

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
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
	"google.golang.org/grpc"

	"shared-go/egress"
	"shared-go/instrumentation"
	"shared-go/logging"
	"shared-go/logging/zerologto"
	"shared-go/rpc"
	"shared-go/web"
)

type gatewayConfig struct {
	Environment     string        `envconfig:"ENVIRONMENT" default:"production"`
	ServiceID       string        `envconfig:"SERVICE_ID" required:"true"`
	LogLevel        logging.Level `envconfig:"LOG_LEVEL" default:"WARN" desc:"options: [TRACE,DEBUG,INFO,WARN,ERROR,PANIC]"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" desc:"time to allow services to gracefully stop"`
	Web             webCfg        `envconfig:"WEB"` // Web Config
	Rpc             rpc.ClientCfg `envconfig:"RPC"` // RPC Config
}

type Gateway struct {
	appFn     func(*Gateway) error
	app       *cobra.Command
	Cfg       *gatewayConfig
	Logger    zerolog.Logger
	WebServer web.Server
	Clients   map[string]*grpc.ClientConn
}

func NewGateway(appFn func(*Gateway) error) *Gateway {
	gateway := &Gateway{
		appFn: appFn,
		Cfg:   &gatewayConfig{},
	}

	gateway.app = &cobra.Command{
		Use:           "service",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          gateway.run,
	}

	gateway.app.Flags().StringSliceVarP(&envFiles, "envfiles", "f", []string{}, "environment variable override files")

	appUsage := gateway.app.UsageString()

	gateway.app.SetUsageFunc(func(command *cobra.Command) error {
		fmt.Println(appUsage)
		fmt.Println("Environment Variables:")
		return envconfig.Usage("", gateway.Cfg)
	})

	return gateway
}

func (g Gateway) Execute() error {
	cobra.OnInitialize(g.initConfig)

	return g.app.Execute()
}

func (g *Gateway) initConfig() {
	var err error

	appPrefix := os.Getenv("APP_PREFIX")

	if len(envFiles) >= 1 {
		err = godotenv.Load(envFiles...)
		if err != nil {
			fmt.Println("error reading environment variable overrides", err, "files:", envFiles)
			os.Exit(1)
		}
	}

	err = envconfig.Process(appPrefix, g.Cfg)
	if err != nil {
		fmt.Println("error reading environment variables", err)
		os.Exit(1)
	}
}

func (g *Gateway) run(*cobra.Command, []string) error {
	var err error

	g.Logger, err = logging.NewZeroLogger(logging.Config{
		Environment: g.Cfg.Environment,
		LogLevel:    g.Cfg.LogLevel,
	})
	if err != nil {
		return err
	}

	log.DefaultLogger = zerologto.Logger(g.Logger)

	g.Clients, err = initRpcClients(g.Cfg.Rpc)
	if err != nil {
		return err
	}

	g.WebServer = web.NewServer(g.Cfg.Web.Http, web.WithHealthCheck(g.Cfg.Web.PingPath))

	g.WebServer.Options(g.Cfg.Web.ApiPath,
		web.WithSecure(),
		web.WithCors(g.Cfg.Web.Cors),
		web.WithMiddleware(
			instrumentation.WebInstrumentation(),
			web.ZeroLogger(g.Logger),
			http2.RequestContext,
		))

	g.WebServer.Mount(g.Cfg.Web.MetricsPath, func(r chi.Router) http.Handler {
		return promhttp.Handler()
	})

	err = g.appFn(g)
	if err != nil {
		return err
	}

	waiter := egress.NewWaiter()

	waiter.Add(g.waitForWebServer)
	waiter.Add(g.waitForConnections)

	g.Logger.Debug().Msg("backend-for-frontend starting")

	return waiter.Wait()
}

func (g Gateway) waitForWebServer(ctx context.Context) (err error) {
	defer g.Logger.Debug().Msg("web server has been shutdown")

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		defer g.Logger.Debug().Msg("web server exited")
		err = g.WebServer.Start()
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		g.Logger.Debug().Msg("shutting down the web server")
		ctx, cancel := context.WithTimeout(context.Background(), g.Cfg.ShutdownTimeout)
		defer cancel()

		if err = g.WebServer.Shutdown(ctx); err != nil {
			g.Logger.Warn().Msg("timed out while shutting down the web server")
		}

		return nil
	})

	return group.Wait()
}

func (g Gateway) waitForConnections(ctx context.Context) error {
	var err error
	defer g.Logger.Debug().Msg("cleanup has been completed")

	<-ctx.Done()
	wg := sync.WaitGroup{}

	for _, client := range g.Clients {
		wg.Add(1)
		go func(c *grpc.ClientConn) {
			sCtx, cancel := context.WithTimeout(context.Background(), g.Cfg.ShutdownTimeout)
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
