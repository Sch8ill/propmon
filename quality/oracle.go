package quality

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sch8ill/propmon/proposal"
)

const path = "/api/v3/providers/detailed?country=DE"

type Oracle struct {
	url    string
	client *http.Client
}

func NewOracle(url string) *Oracle {
	return &Oracle{
		url:    url,
		client: http.DefaultClient,
	}
}

func (o *Oracle) Quality() (map[string]*proposal.Quality, error) {
	res, err := o.client.Get(o.url + path)
	if err != nil {
		return nil, fmt.Errorf("quality oracle request failed: %w", err)
	}

	defer res.Body.Close()

	var qualityRes map[string]*proposal.Quality
	if err := json.NewDecoder(res.Body).Decode(&qualityRes); err != nil {
		return nil, fmt.Errorf("failed to decode quality oracle response : %w", err)
	}

	return qualityRes, nil
}
