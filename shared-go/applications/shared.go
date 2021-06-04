package applications

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	edatpgx "github.com/stackus/edat-pgx"
	http2 "github.com/stackus/edat/http"
	"google.golang.org/grpc"

	"shared-go/instrumentation"
	"shared-go/logging/zerologto"
	"shared-go/rpc"
	"shared-go/web"
)

const OrderService = "orderservice"
const ConsumerService = "consumerservice"
const OrderHistoryService = "orderhistoryserver"

type CleanupFunc func(ctx context.Context) error

var envFiles []string

type pgCfg struct {
	Conn string `envconfig:"CONN" desc:"a postgres DATABASE_URL or CONNECTION_STRING"`
}

type natsCfg struct {
	URL            string        `envconfig:"URL"`
	ClusterID      string        `envconfig:"CLUSTER_ID"`
	AckWaitTimeout time.Duration `envconfig:"ACK_WAIT_TIMEOUT" default:"30s"`
}

type kafkaCfg struct {
	Brokers []string `envconfig:"BROKERS"`
}

type webCfg struct {
	ApiPath     string        `envconfig:"API_PATH" default:"/api"`
	PingPath    string        `envconfig:"PING_PATH" default:"/ping"`
	MetricsPath string        `envconfig:"METRICS_PATH" default:"/metrics"`
	Http        web.ServerCfg `envconfig:"HTTP"`
	Cors        web.CorsCfg   `envconfig:"CORS"`
}

type rpcCfg struct {
	Server rpc.ServerCfg `envconfig:"SERVER"`
	Client rpc.ClientCfg `envconfig:"CLIENT"`
}

func initWebServer(cfg webCfg, poolConn *pgxpool.Pool, logger zerolog.Logger) web.Server {
	webServer := web.NewServer(cfg.Http, web.WithHealthCheck(cfg.PingPath))

	webServer.Options(cfg.ApiPath,
		web.WithSecure(),
		web.WithCors(cfg.Cors),
		web.WithMiddleware(
			instrumentation.WebInstrumentation(),
			web.ZeroLogger(logger),
			http2.RequestContext,
			// 4. Outbox: Use a WEB request middleware to start a new transaction for each incoming request
			edatpgx.WebSessionMiddleware(poolConn, zerologto.Logger(logger)),
		))

	webServer.Mount(cfg.MetricsPath, func(r chi.Router) http.Handler {
		return promhttp.Handler()
	})

	return webServer
}

func initRpcServer(cfg rpc.ServerCfg, poolConn *pgxpool.Pool, logger zerolog.Logger) rpc.Server {
	return rpc.NewServer(cfg,
		rpc.WithServerUnaryInterceptors(
			// 4. Outbox: Use a RPC request middleware to start a new transaction for each incoming request
			edatpgx.RpcSessionUnaryInterceptor(poolConn, zerologto.Logger(logger)),
		),
		rpc.WithServerUnaryEnsureStatus(),
	)
}

func initRpcClients(cfg rpc.ClientCfg) (map[string]*grpc.ClientConn, error) {
	var err error

	clients := map[string]*grpc.ClientConn{}

	clients[OrderService], err = rpc.NewClientConn(cfg, "orderservice:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return nil, err
	}

	clients[ConsumerService], err = rpc.NewClientConn(cfg, "consumerservice:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return nil, err
	}

	clients[OrderHistoryService], err = rpc.NewClientConn(cfg, "orderhistoryservice:8000", rpc.WithClientUnaryConvertStatus())
	if err != nil {
		return nil, err
	}

	return clients, nil
}
