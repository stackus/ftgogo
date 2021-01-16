package instrumentation

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stackus/edat/msg"
)

func MessageInstrumentation() func(msg.MessageReceiver) msg.MessageReceiver {
	responseTime := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "message_response_time",
		Help:    "Message response time in microseconds",
		Buckets: []float64{300, 600, 900, 1_500, 5_000, 10_000, 20_000},
	})

	return func(next msg.MessageReceiver) msg.MessageReceiver {
		return msg.ReceiveMessageFunc(func(ctx context.Context, message msg.Message) error {
			start := time.Now()
			err := next.ReceiveMessage(ctx, message)
			responseTime.Observe(float64(time.Since(start).Microseconds()))

			return err
		})
	}
}
