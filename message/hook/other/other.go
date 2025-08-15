// Package other is a simple package that provides a other hook.
package other

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aide-family/magicbox/httpx"
	"github.com/aide-family/magicbox/message"
)

var _ message.Sender = (*otherHookSender)(nil)
var _ message.Driver = (*initializer)(nil)

func SenderDriver(config Config) message.Driver {
	return &initializer{config: config}
}

type initializer struct {
	config Config
}

// New implements message.Driver.
func (i *initializer) New() (message.Sender, error) {
	return &otherHookSender{
		cli:    httpx.NewClient(httpx.GetHTTPClient()),
		config: i.config,
	}, nil
}

type otherHookSender struct {
	cli    *httpx.Client
	config Config
}

// Send implements message.Sender.
func (o *otherHookSender) Send(ctx context.Context, message message.Message) error {
	opts := []httpx.Option{
		httpx.WithHeaders(o.config.GetHeaders()),
		httpx.WithBasicAuth(o.config.GetBasicAuth()),
	}

	resp, err := o.cli.Post(ctx, o.config.GetURL(), message, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
}
