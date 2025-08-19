package email

import (
	"encoding/json"

	"github.com/aide-family/magicbox/message"
)

var _ message.Message = (*Message)(nil)

type Attachment struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

type Message struct {
	To          []string            `json:"to"`
	Cc          []string            `json:"cc"`
	Subject     string              `json:"subject"`
	Body        string              `json:"body"`
	ContentType string              `json:"contentType"`
	Attachments []*Attachment       `json:"attachments"`
	Headers     map[string][]string `json:"headers"`
}

func (m *Message) Message() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return json
}
