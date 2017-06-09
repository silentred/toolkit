package filter

import (
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "req_count",
		Help: "total count of request",
	}, []string{"service"})

	reqDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "req_duration_us",
		Help: "latency of request in microsecond",
	}, []string{"service"})
)

// func init() {
// 	prometheus.MustRegister(reqCount)
// 	prometheus.MustRegister(reqDuration)

// 	// Expose the registered metrics via HTTP.
// 	http.Handle("/metrics", promhttp.Handler())
// 	go http.ListenAndServe(":8088", nil)
// }

func GetPrometheusLogHandler() echo.HandlerFunc {
	prometheus.MustRegister(reqCount)
	prometheus.MustRegister(reqDuration)

	handler := func(e echo.Context) error {
		promhttp.Handler().ServeHTTP(e.Response().Writer, e.Request())
		return nil
	}

	return handler
}

func Metrics() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			t := time.Now()

			if err = next(c); err != nil {
				c.Error(err)
			}

			duration := float64(time.Now().Sub(t).Nanoseconds()) / 1000

			reqCount.WithLabelValues("http").Add(1)
			reqDuration.WithLabelValues("http").Observe(duration)

			return nil
		}
	}
}
