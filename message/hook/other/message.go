package other

import (
	"encoding/json"

	"github.com/aide-family/magicbox/message"
)

var _ message.Message = (*Message)(nil)

type Message map[string]any

func (m *Message) Message() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return json
}
