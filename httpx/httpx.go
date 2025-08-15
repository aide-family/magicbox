// Package httpx is a simple package that provides a http client.
package httpx

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
)

var (
	httpClient     = http.DefaultClient
	httpClientOnce sync.Once
)

func SetHTTPClient(cli *http.Client) {
	httpClientOnce.Do(func() {
		httpClient = cli
	})
}

func GetHTTPClient() *http.Client {
	return httpClient
}

type Option func(*http.Request)

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) Get(ctx context.Context, url string, options ...Option) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		option(req)
	}
	return c.client.Do(req.WithContext(ctx))
}

func (c *Client) Post(ctx context.Context, url string, body []byte, options ...Option) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for _, option := range options {
		option(req)
	}
	return c.client.Do(req.WithContext(ctx))
}

func WithBasicAuth(basicAuth *BasicAuth) Option {
	return func(req *http.Request) {
		if basicAuth == nil {
			return
		}
		req.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}
}

func WithHeaders(headers map[string][]string) Option {
	return func(req *http.Request) {
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
}

func WithQuery(query url.Values) Option {
	return func(req *http.Request) {
		oldQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				oldQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = oldQuery.Encode()
	}
}

func WithTLS(tlsConfig *tls.ConnectionState) Option {
	return func(req *http.Request) {
		req.TLS = tlsConfig
	}
}
