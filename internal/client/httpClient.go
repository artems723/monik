package client

import "github.com/go-resty/resty/v2"

type HTTPClient struct {
	client *resty.Client
}

func NewHTTPClient() HTTPClient {
	client := resty.New()
	return HTTPClient{client: client}
}
