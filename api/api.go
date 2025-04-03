package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sch8ill/propmon/metrics"
)

type API struct {
	address string
}

func New(address string) *API {
	return &API{address: address}
}

func (a *API) Run() error {
	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})))

	return r.Run(a.address)
}
