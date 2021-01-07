package egress

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// WaiterFn is given a cancellable context and will bubble up any error to the waiting process
type WaiterFn func(ctx context.Context) error

type waiterCfg struct {
	signals []os.Signal
}

type waiter struct {
	ctx    context.Context
	group  *errgroup.Group
	cancel context.CancelFunc
}

type Waiter interface {
	Add(fn WaiterFn)
	Wait() error
}

func NewWaiter(options ...WaiterOption) Waiter {
	ctx, cancel := context.WithCancel(context.Background())
	group, groupCtx := errgroup.WithContext(ctx)

	cfg := &waiterCfg{
		signals: []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM},
	}

	for _, option := range options {
		option(cfg)
	}

	w := &waiter{
		ctx:    groupCtx,
		group:  group,
		cancel: cancel,
	}

	w.group.Go(func() error {
		defer w.cancel()

		s := make(chan os.Signal, 1)
		signal.Notify(s, cfg.signals...)

		select {
		case <-s:
		case <-w.ctx.Done():
		}

		return nil
	})

	return w
}

func (w *waiter) Add(fn WaiterFn) {
	w.group.Go(func() error {
		return fn(w.ctx)
	})
}

// Wait for either a stop signal or any process to end normally or with an error
func (w waiter) Wait() error {
	if err := w.group.Wait(); err != nil && err != context.Canceled {
		return err
	}

	return nil
}
