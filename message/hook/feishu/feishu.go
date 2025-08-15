// Package feishu is a simple package that provides a feishu hook.
package feishu

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
	"time"

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

	msg := make(map[string]any)
	if err := json.Unmarshal(message, &msg); err != nil {
		return err
	}
	requestTime := time.Now().Unix()
	msg["timestamp"] = strconv.FormatInt(requestTime, 10)
	sign, err := f.sign(requestTime)
	if err != nil {
		return err
	}
	msg["sign"] = sign
	requestBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	u, err := url.Parse(f.config.GetURL())
	if err != nil {
		return err
	}
	u.Path += "/" + f.config.GetKey()
	resp, err := f.cli.Post(ctx, u.String(), requestBody, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return hook.RequestAssert(resp, unmarshalResponse)
}

func (f *feishuHookSender) sign(timestamp int64) (string, error) {
	// timestamp + key sha256, then base64 encode
	signString := strconv.FormatInt(timestamp, 10) + "\n" + f.config.GetSecret()

	var data []byte
	h := hmac.New(sha256.New, []byte(signString))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
