package alicloud

import (
	"encoding/json"

	"github.com/aide-family/magicbox/message"
)

var _ message.Message = (*Message)(nil)

type Message struct {
	TemplateParam string   `json:"templateParam"`
	TemplateCode  string   `json:"templateCode"`
	PhoneNumbers  []string `json:"phoneNumbers"`
}

func (m *Message) Message(channel message.MessageChannel) ([]byte, error) {
	if err := MessageChannelSMSAliCloud.Check(channel); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
