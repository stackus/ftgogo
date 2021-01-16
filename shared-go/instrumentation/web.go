package instrumentation

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func WebInstrumentation() func(http.Handler) http.Handler {
	responseTime := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "web_response_time",
		Help:    "Web response time in milliseconds",
		Buckets: []float64{300, 600, 900, 1_500, 5_000, 10_000, 20_000},
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			next.ServeHTTP(writer, request)
			responseTime.Observe(float64(time.Since(start).Milliseconds()))
		})
	}
}
