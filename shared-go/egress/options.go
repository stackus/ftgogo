package egress

import (
	"os"
)

type WaiterOption func(*waiterCfg)

func WithSignals(signals ...os.Signal) WaiterOption {
	return func(cfg *waiterCfg) {
		cfg.signals = signals
	}
}
