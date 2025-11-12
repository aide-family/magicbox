// Package wechat is a simple package that provides a wechat hook.
package wechat

import (
	"context"
	"net/url"

	"github.com/aide-family/magicbox/httpx"
	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/hook"
	"github.com/aide-family/magicbox/serialize"
)

var _ message.Sender = (*wechatHookSender)(nil)
var _ message.Driver = (*initializer)(nil)

const MessageChannelWechat message.MessageChannel = "webhook-wechat"

func SenderDriver(config hook.Config) message.Driver {
	return &initializer{config: config}
}

type initializer struct {
	config hook.Config
}

// New implements message.Driver.
func (i *initializer) New() (message.Sender, error) {
	return &wechatHookSender{
		cli:    httpx.NewClient(httpx.GetHTTPClient()),
		config: i.config,
	}, nil
}

type wechatHookSender struct {
	cli    *httpx.Client
	config hook.Config
}

// Send implements message.Sender.
func (w *wechatHookSender) Send(ctx context.Context, msg message.Message) error {
	opts := []httpx.Option{
		httpx.WithHeaders(map[string][]string{
			"Content-Type": {"application/json"},
		}),
		httpx.WithQuery(url.Values{
			"key": {w.config.GetSecret()},
		}),
	}
	var newMessage *Message
	var ok bool
	if newMessage, ok = msg.(*Message); !ok {
		jsonBytes, err := msg.Message(MessageChannelWechat)
		if err != nil {
			return err
		}
		if err := serialize.JSONUnmarshal(jsonBytes, newMessage); err != nil {
			return err
		}
	}
	jsonBytes, err := newMessage.Message(MessageChannelWechat)
	if err != nil {
		return err
	}
	resp, err := w.cli.Post(ctx, w.config.GetURL(), jsonBytes, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return hook.RequestAssert(resp, unmarshalResponse)
}
