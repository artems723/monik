package httpClient

import "github.com/go-resty/resty/v2"

type HTTPClient struct {
	Client *resty.Client
}

func NewHTTPClient() HTTPClient {
	client := resty.New()
	return HTTPClient{Client: client}
}
