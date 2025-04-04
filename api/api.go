package api

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sch8ill/propmon/metrics"
	"github.com/sch8ill/propmon/proposal"
)

type API struct {
	address    string
	repository *proposal.Repository
}

func New(address string, repository *proposal.Repository) *API {
	return &API{
		address:    address,
		repository: repository,
	}
}

func (a *API) Run() error {
	gin.DefaultWriter = io.Discard
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})))

	api := r.Group("/api/v1")
	rateLimiter := newRateLimiter(20, time.Minute)
	api.Use(rateLimiter.middleware())

	handler := newHandler(a.repository)
	api.GET("/proposals", handler.getProposals)

	return r.Run(a.address)
}
