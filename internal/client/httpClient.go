package client

import "github.com/go-resty/resty/v2"

type Client struct {
	client *resty.Client
}

func New() Client {
	client := resty.New()
	return Client{client: client}
}
