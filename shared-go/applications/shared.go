package applications

import (
	"time"

	"shared-go/web"
)

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
	ApiPath  string        `envconfig:"API_PATH" default:"/api"`
	PingPath string        `envconfig:"PING_PATH" default:"/ping"`
	Http     web.ServerCfg `envconfig:"HTTP"`
	Cors     web.CorsCfg   `envconfig:"CORS"`
}
