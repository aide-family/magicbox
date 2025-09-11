// Package alicloud implements the alicloud sms driver.
package alicloud

import (
	"context"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapiv3 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"

	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/pointer"
	"github.com/aide-family/magicbox/serialize"
)

var _ message.Driver = (*initializer)(nil)
var _ message.Sender = (*alicloudSmsSender)(nil)

const MessageChannelSMSAliCloud message.MessageChannel = "sms-alicloud"

func SenderDriver(config Config) message.Driver {
	return &initializer{config: config}
}

type initializer struct {
	config Config
}

// New implements message.Driver.
func (i *initializer) New() (message.Sender, error) {
	clientV3, err := i.initV3()
	if err != nil {
		return nil, err
	}
	return &alicloudSmsSender{config: i.config, clientV3: clientV3}, nil
}

func (a *initializer) initV3() (*dysmsapiv3.Client, error) {
	accessKeyID := a.config.GetAccessKeyID()
	accessKeySecret := a.config.GetAccessKeySecret()
	if accessKeyID == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("SMS sending credential information is not configured")
	}
	config := &openapi.Config{
		AccessKeyId:     &accessKeyID,
		AccessKeySecret: &accessKeySecret,
	}

	if endpoint := a.config.GetEndpoint(); endpoint != "" {
		config.Endpoint = tea.String(endpoint)
	}
	client, err := dysmsapiv3.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SMS clientV3: %w", err)
	}
	return client, nil
}

type alicloudSmsSender struct {
	config   Config
	clientV3 *dysmsapiv3.Client
}

// Send implements message.Sender.
func (a *alicloudSmsSender) Send(ctx context.Context, msg message.Message) error {
	newMessage := &Message{}
	var ok bool
	if newMessage, ok = msg.(*Message); !ok {
		jsonBytes, err := msg.Message(MessageChannelSMSAliCloud)
		if err != nil {
			return err
		}
		if err := serialize.JSONUnmarshal(jsonBytes, newMessage); err != nil {
			return err
		}
	}
	if len(newMessage.PhoneNumbers) == 0 {
		return fmt.Errorf("phone numbers is empty")
	}

	if len(newMessage.PhoneNumbers) == 1 {
		return a.send(ctx, newMessage)
	}
	return a.batchSend(ctx, newMessage)
}

func (a *alicloudSmsSender) send(_ context.Context, message *Message) error {
	phoneNumber := message.PhoneNumbers[0]
	sendSmsRequest := &dysmsapiv3.SendSmsRequest{
		PhoneNumbers:  pointer.Of(phoneNumber),
		SignName:      pointer.Of(a.config.GetSignName()),
		TemplateCode:  pointer.Of(message.TemplateCode),
		TemplateParam: pointer.Of(message.TemplateParam),
	}

	response, err := a.clientV3.SendSmsWithOptions(sendSmsRequest, runtimeOptions)
	if err != nil {
		return err
	}
	body := pointer.Get(response.Body)
	if pointer.Get(body.Code) == "OK" {
		return nil
	}

	return fmt.Errorf("send sms failed: %v", body)
}

func (a *alicloudSmsSender) batchSend(_ context.Context, message *Message) error {
	phoneNumbers := message.PhoneNumbers
	signNames := make([]string, 0, len(phoneNumbers))
	templateParams := make([]string, 0, len(phoneNumbers))
	for range phoneNumbers {
		signNames = append(signNames, a.config.GetSignName())
		templateParams = append(templateParams, message.TemplateParam)
	}

	phoneNumberJson, err := serialize.JSONMarshal(phoneNumbers)
	if err != nil {
		return fmt.Errorf("failed to marshal phone numbers: %v", err)
	}
	signNameJson, err := serialize.JSONMarshal(signNames)
	if err != nil {
		return fmt.Errorf("failed to marshal sign names: %v", err)
	}
	templateParamJson, err := serialize.JSONMarshal(templateParams)
	if err != nil {
		return fmt.Errorf("failed to marshal template params: %v", err)
	}
	sendBatchSmsRequest := &dysmsapiv3.SendBatchSmsRequest{
		PhoneNumberJson:   pointer.Of(string(phoneNumberJson)),
		SignNameJson:      pointer.Of(string(signNameJson)),
		TemplateParamJson: pointer.Of(string(templateParamJson)),
		TemplateCode:      pointer.Of(message.TemplateCode),
	}

	response, err := a.clientV3.SendBatchSmsWithOptions(sendBatchSmsRequest, runtimeOptions)
	if err != nil {
		return err
	}
	body := pointer.Get(response.Body)
	if pointer.Get(body.Code) == "OK" {
		return nil
	}
	return fmt.Errorf("send batch sms failed: %v", body)
}
