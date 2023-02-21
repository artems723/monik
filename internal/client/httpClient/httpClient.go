package httpClient

import (
	"encoding/json"
	"github.com/artems723/monik/internal/client"
	"github.com/go-resty/resty/v2"
	"log"
)

type Client interface {
	SendData([]*client.Metric, string) ([]client.Metric, error)
}

type HTTPClient struct {
	client *resty.Client
}

func NewHTTPClient() HTTPClient {
	client := resty.New()
	return HTTPClient{client: client}
}

func (c HTTPClient) SendData(metrics []*client.Metric, URL string) ([]client.Metric, error) {
	m, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("httpClient.SendData: unable to marshal. Error: %v. Metric: %v", err, metrics)
		return nil, err
	}
	var result []client.Metric
	_, err = c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(m).
		SetResult(&result).
		Post(URL)
	if err != nil {
		log.Printf("httpClient.SendData: error sending request: %s", err)
		return nil, err
	}
	return result, nil
}
