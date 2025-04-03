package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// custom registry to discard default go metrics
var registry = prometheus.NewRegistry()

func Listen(addr string) error {
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}

	return nil
}
