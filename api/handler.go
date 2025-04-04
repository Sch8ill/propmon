package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sch8ill/propmon/proposal"
)

const maxResponseCount = 1000

type handler struct {
	repository *proposal.Repository
}

func newHandler(repository *proposal.Repository) *handler {
	return &handler{repository: repository}
}

func (h *handler) getProposals(c *gin.Context) {
	max, err := strconv.Atoi(c.Query("max"))
	if err != nil {
		max = 100
	}
	if max > maxResponseCount {
		max = maxResponseCount
	}

	filter := &proposal.Proposal{
		ProviderID:  c.Query("id"),
		ServiceType: c.Query("service"),
		Location: proposal.Location{
			Country: c.Query("country"),
			IpType:  c.Query("type"),
		},
	}

	proposals := h.repository.Match(filter, max)
	if len(proposals) == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, proposals)
}
