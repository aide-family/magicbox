// Package dingtalk implements the dingtalk hook driver.
package dingtalk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/aide-family/magicbox/httpx"
	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/message/hook"
)

var _ message.Sender = (*dingtalkHookSender)(nil)
var _ message.Driver = (*initializer)(nil)

func SenderDriver(config Config) message.Driver {
	return &initializer{config: config}
}

type initializer struct {
	config Config
}

// New implements message.Driver.
func (i *initializer) New() (message.Sender, error) {
	return &dingtalkHookSender{
		cli:    httpx.NewClient(httpx.GetHTTPClient()),
		config: i.config,
	}, nil
}

type dingtalkHookSender struct {
	cli    *httpx.Client
	config Config
}

func (d *dingtalkHookSender) Send(ctx context.Context, message message.Message) error {
	timestamp := time.Now().UnixMilli()
	opts := []httpx.Option{
		httpx.WithHeaders(map[string][]string{
			"Content-Type": {"application/json"},
		}),
		httpx.WithQuery(url.Values{
			"access_token": {d.config.GetKey()},
			"timestamp":    {strconv.FormatInt(timestamp, 10)},
			"sign":         {d.sign(timestamp)},
		}),
	}

	resp, err := d.cli.Post(ctx, d.config.GetURL(), message, opts...)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return hook.RequestAssert(resp, unmarshalResponse)
}

func (d *dingtalkHookSender) sign(timestamp int64) string {
	message := fmt.Sprintf("%d\n%s", timestamp, d.config.GetSecret())
	key := []byte(d.config.GetSecret())
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	signature := h.Sum(nil)
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	signatureURLEncoded := url.QueryEscape(signatureBase64)
	return signatureURLEncoded
}
