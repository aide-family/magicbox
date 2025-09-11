package other

import (
	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/serialize"
)

var _ message.Message = (*Message)(nil)

type Message map[string]any

func (m *Message) Message(channel message.MessageChannel) ([]byte, error) {
	json, err := serialize.JSONMarshal(m)
	if err != nil {
		return nil, err
	}
	return json, nil
}
