package other

import (
	"encoding/json"

	"github.com/aide-family/magicbox/message"
)

var _ message.Message = (*Message)(nil)

type Message map[string]any

func (m *Message) Message(channel message.MessageChannel) ([]byte, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return json, nil
}
