// Package other is a simple package that provides a other hook.
package other

import (
	"context"
	"io"

	"github.com/aide-family/magicbox/httpx"
	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/hook"
)

var _ message.Sender = (*otherHookSender)(nil)
var _ message.Driver = (*initializer)(nil)

const MessageChannelOther message.MessageChannel = "webhook-other"

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

	jsonBytes, err := message.Message(MessageChannelOther)
	if err != nil {
		return err
	}

	resp, err := o.cli.Post(ctx, o.config.GetURL(), jsonBytes, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return hook.RequestAssert(resp, unmarshalResponse)
}

func unmarshalResponse(body io.ReadCloser) error {
	return nil
}
