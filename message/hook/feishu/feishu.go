// Package feishu is a simple package that provides a feishu hook.
package feishu

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/aide-family/magicbox/httpx"
	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/hook"
)

var _ message.Sender = (*feishuHookSender)(nil)
var _ message.Driver = (*initializer)(nil)

func SenderDriver(config Config) message.Driver {
	return &initializer{config: config}
}

type initializer struct {
	config Config
}

// New implements message.Driver.
func (i *initializer) New() (message.Sender, error) {
	return &feishuHookSender{
		cli:    httpx.NewClient(httpx.GetHTTPClient()),
		config: i.config,
	}, nil
}

type feishuHookSender struct {
	cli    *httpx.Client
	config Config
}

// Send implements message.Sender.
func (f *feishuHookSender) Send(ctx context.Context, message message.Message) error {
	opts := []httpx.Option{
		httpx.WithHeaders(map[string][]string{
			"Content-Type": {"application/json"},
		}),
	}

	feishuMessage := &Message{}
	var ok bool
	if feishuMessage, ok = message.(*Message); !ok {
		jsonBytes, err := message.Message(MessageChannelFeishu)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(jsonBytes, feishuMessage); err != nil {
			return err
		}
	}
	if err := feishuMessage.Signature(f.config.GetSecret()); err != nil {
		return err
	}

	u, err := url.Parse(f.config.GetURL())
	if err != nil {
		return err
	}
	u.Path += "/" + f.config.GetKey()
	jsonBytes, err := feishuMessage.Message(MessageChannelFeishu)
	if err != nil {
		return err
	}
	resp, err := f.cli.Post(ctx, u.String(), jsonBytes, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return hook.RequestAssert(resp, unmarshalResponse)
}
