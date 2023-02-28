package httpclient

import (
	"encoding/json"
	"github.com/artems723/monik/internal/client/agent"
	"github.com/go-resty/resty/v2"
	"log"
)

type HTTPClient struct {
	client    *resty.Client
	rateLimit int
	jobs      chan struct{}
}

func New(rateLimit int) HTTPClient {
	client := resty.New()
	return HTTPClient{
		client:    client,
		rateLimit: rateLimit,
		jobs:      make(chan struct{}, rateLimit),
	}
}

func (c HTTPClient) SendData(metrics []*agent.Metric, URL string) ([]agent.Metric, error) {
	defer func() {
		<-c.jobs // release worker
	}()
	c.jobs <- struct{}{} // acquire worker
	m, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("httpclient.SendData: unable to marshal. Error: %v. Metric: %v", err, metrics)
		return nil, err
	}
	var result []agent.Metric
	_, err = c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(m).
		SetResult(&result).
		Post(URL)
	if err != nil {
		log.Printf("httpclient.SendData: error sending request: %s", err)
		return nil, err
	}
	return result, nil
}