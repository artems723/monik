package client

import "github.com/go-resty/resty/v2"

type HttpClient struct {
	client *resty.Client
}

func NewHttpClient() HttpClient {
	client := resty.New()
	return HttpClient{client: client}
}
