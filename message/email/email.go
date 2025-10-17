// Package email is a simple package that provides a email sender.
package email

import (
	"context"
	"io"

	"gopkg.in/gomail.v2"

	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/serialize"
)

var _ message.Sender = (*emailSender)(nil)
var _ message.Driver = (*initializer)(nil)

func SenderDriver(config Config) message.Driver {
	return &initializer{config: config}
}

type initializer struct {
	config Config
}

func (i *initializer) New() (message.Sender, error) {
	host, port := i.config.GetHost(), i.config.GetPort()
	username, password := i.config.GetUsername(), i.config.GetPassword()
	dialer := gomail.NewDialer(host, int(port), username, password)
	return &emailSender{dialer: dialer, config: i.config}, nil
}

type emailSender struct {
	dialer *gomail.Dialer
	config Config
}

func (e *emailSender) Send(ctx context.Context, m message.Message) error {
	emailMessage := &Message{}
	var ok bool
	if emailMessage, ok = m.(*Message); !ok {
		jsonBytes, err := m.Message(MessageChannelEmail)
		if err != nil {
			return err
		}
		if err := serialize.JSONUnmarshal(jsonBytes, emailMessage); err != nil {
			return err
		}
	}
	msg := gomail.NewMessage(gomail.SetCharset("UTF-8"), gomail.SetEncoding(gomail.Base64))
	msg.SetHeader("From", e.config.GetUsername())
	msg.SetHeader("To", emailMessage.To...)
	msg.SetHeader("Cc", emailMessage.Cc...)
	msg.SetHeader("Subject", emailMessage.Subject)
	msg.SetBody(emailMessage.ContentType, emailMessage.Body)
	for _, attachment := range emailMessage.Attachments {
		msg.Attach(attachment.Filename, gomail.SetHeader(map[string][]string{
			"Content-Disposition": {"attachment"},
		}), gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Data)
			return err
		}))
	}
	for key, values := range emailMessage.Headers {
		msg.SetHeader(key, values...)
	}
	return e.dialer.DialAndSend(msg)
}
