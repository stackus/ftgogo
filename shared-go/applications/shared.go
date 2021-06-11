package applications

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	edatpgx "github.com/stackus/edat-pgx"
	"github.com/stackus/errors"
	"google.golang.org/grpc"

	http2 "github.com/stackus/edat/http"

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
			http2.RequestContext,
			web.ZeroLogger(logger),
			edatpgx.WebSessionMiddleware(poolConn, zerologto.Logger(logger)),
		))

	webServer.Mount(cfg.MetricsPath, func(r chi.Router) http.Handler {
		return promhttp.Handler()
	})

	return webServer
}

func initRpcServer(cfg rpc.ServerCfg, poolConn *pgxpool.Pool, logger zerolog.Logger) rpc.Server {
	return rpc.NewServer(cfg,
		rpc.WithUnaryServerInterceptors(
			rpc.RequestContextUnaryServerInterceptor,
			rpc.WithUnaryServerLogging(logger),
			edatpgx.RpcSessionUnaryInterceptor(poolConn, zerologto.Logger(logger)),
		),
		rpc.WithServerUnaryEnsureStatus(),
	)
}

func initRpcClients(cfg rpcCfg, logger zerolog.Logger) (map[string]*grpc.ClientConn, error) {
	var err error

	clients := map[string]*grpc.ClientConn{}

	if cfg.Server.Network == "unix" {
		conn, err := grpc.Dial(cfg.Server.Address, grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			addr, err := net.ResolveUnixAddr(cfg.Server.Network, cfg.Server.Address)
			if err != nil {
				return nil, err
			}
			return net.DialUnix(cfg.Server.Network, nil, addr)
		}), grpc.WithInsecure(),
			grpc.WithChainUnaryInterceptor(
				rpc.RequestContextUnaryClientInterceptor,
				rpc.WithUnaryClientLogging(logger),
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
					return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
				},
			),
		)
		if err != nil {
			return nil, err
		}

		clients[OrderService] = conn
		clients[ConsumerService] = conn
		clients[OrderHistoryService] = conn

		return clients, nil
	}

	clients[OrderService], err = rpc.NewClientConn(cfg.Client, "orderservice:8000",
		rpc.WithUnaryClientInterceptors(
			rpc.RequestContextUnaryClientInterceptor,
			rpc.WithUnaryClientLogging(logger),
		),
		rpc.WithClientUnaryConvertStatus(),
	)
	if err != nil {
		return nil, err
	}

	clients[ConsumerService], err = rpc.NewClientConn(cfg.Client, "consumerservice:8000",
		rpc.WithUnaryClientInterceptors(
			rpc.RequestContextUnaryClientInterceptor,
			rpc.WithUnaryClientLogging(logger),
		),
		rpc.WithClientUnaryConvertStatus(),
	)
	if err != nil {
		return nil, err
	}

	clients[OrderHistoryService], err = rpc.NewClientConn(cfg.Client, "orderhistoryservice:8000",
		rpc.WithUnaryClientInterceptors(
			rpc.RequestContextUnaryClientInterceptor,
			rpc.WithUnaryClientLogging(logger),
		),
		rpc.WithClientUnaryConvertStatus(),
	)
	if err != nil {
		return nil, err
	}

	return clients, nil
}
