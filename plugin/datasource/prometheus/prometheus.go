// Package prometheus provides a client for Prometheus.
package prometheus

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aide-family/magicbox/httpx"
	"github.com/aide-family/magicbox/plugin/datasource"
)

type Client struct {
	c datasource.MetricConfig
}

func NewClient(c datasource.MetricConfig) *Client {
	return &Client{c: c}
}

func (p *Client) Proxy(ctx context.Context, w http.ResponseWriter, r *http.Request, target string) error {
	proxyClient := &httpx.ProxyClient{
		Host: p.c.GetEndpoint(),
	}
	return proxyClient.Proxy(ctx, w, r, target)
}

func (p *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*http.Response, error) {
	hx := httpx.NewClient(httpx.GetHTTPClient())

	api, err := url.JoinPath(p.c.GetEndpoint(), "/api/v1/query_range")
	if err != nil {
		return nil, err
	}
	toURL, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	opts := []httpx.Option{
		httpx.WithHeaders(http.Header{
			"Accept":          {"*/*"},
			"Accept-Language": {"zh-CN,zh;q=0.9"},
			"Connection":      {"keep-alive"},
		}),
		httpx.WithBasicAuth(p.c.GetBasicAuth()),
		httpx.WithTLS(p.c.GetTLS()),
		httpx.WithQuery(url.Values{
			"query": {query},
			"start": {strconv.FormatInt(start.Unix(), 10)},
			"end":   {strconv.FormatInt(end.Unix(), 10)},
			"step":  {strconv.FormatInt(int64(step), 10)},
		}),
	}
	return hx.Do(ctx, httpx.MethodGet, toURL.String(), opts...)
}

func (p *Client) Query(ctx context.Context, query string, time time.Time) (*http.Response, error) {
	hx := httpx.NewClient(httpx.GetHTTPClient())

	api, err := url.JoinPath(p.c.GetEndpoint(), "/api/v1/query")
	if err != nil {
		return nil, err
	}
	toURL, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	opts := []httpx.Option{
		httpx.WithHeaders(http.Header{
			"Accept":          {"*/*"},
			"Accept-Language": {"zh-CN,zh;q=0.9"},
			"Connection":      {"keep-alive"},
		}),
		httpx.WithBasicAuth(p.c.GetBasicAuth()),
		httpx.WithTLS(p.c.GetTLS()),
		httpx.WithQuery(url.Values{
			"query": {query},
			"time":  {strconv.FormatInt(time.Unix(), 10)},
		}),
	}
	return hx.Do(ctx, httpx.MethodGet, toURL.String(), opts...)
}

func (p *Client) Series(ctx context.Context, start, end time.Time, match []string) (*http.Response, error) {
	hx := httpx.NewClient(httpx.GetHTTPClient())

	api, err := url.JoinPath(p.c.GetEndpoint(), "/api/v1/series")
	if err != nil {
		return nil, err
	}
	toURL, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	opts := []httpx.Option{
		httpx.WithHeaders(http.Header{
			"Accept":          {"*/*"},
			"Accept-Language": {"zh-CN,zh;q=0.9"},
			"Connection":      {"keep-alive"},
			"Content-Type":    {"application/x-www-form-urlencoded;charset=UTF-8"},
		}),
		httpx.WithBasicAuth(p.c.GetBasicAuth()),
		httpx.WithTLS(p.c.GetTLS()),
		httpx.WithBody([]byte(url.Values{
			"start":   {start.Format(time.RFC3339)},
			"end":     {end.Format(time.RFC3339)},
			"match[]": match,
		}.Encode())),
	}
	return hx.Do(ctx, httpx.MethodPost, toURL.String(), opts...)
}

func (p *Client) Metadata(ctx context.Context, metric string) (*http.Response, error) {
	hx := httpx.NewClient(httpx.GetHTTPClient())

	api, err := url.JoinPath(p.c.GetEndpoint(), "/api/v1/metadata")
	if err != nil {
		return nil, err
	}
	toURL, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	opts := []httpx.Option{
		httpx.WithHeaders(http.Header{
			"Accept":          {"*/*"},
			"Accept-Language": {"zh-CN,zh;q=0.9"},
			"Connection":      {"keep-alive"},
		}),
		httpx.WithBasicAuth(p.c.GetBasicAuth()),
		httpx.WithTLS(p.c.GetTLS()),
	}
	return hx.Do(ctx, httpx.MethodGet, toURL.String(), opts...)
}
