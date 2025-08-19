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

func (m *Message) Message() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return json
}
