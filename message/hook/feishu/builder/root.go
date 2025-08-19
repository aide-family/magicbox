package builder

import (
	"encoding/json"

	"github.com/aide-family/magicbox/message"
)

var _ message.Message = (*Message)(nil)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypePost  MessageType = "post"
	MessageTypeImage MessageType = "image"
	MessageTypeCard  MessageType = "interactive"
)

type Content struct {
	Text  *Text  `json:"text,omitempty"`
	Post  *Post  `json:"post,omitempty"`
	Card  *Card  `json:"card,omitempty"`
	Image string `json:"image_key,omitempty"`
}

type Message struct {
	MsgType MessageType `json:"msg_type"`
	Content *Content    `json:"content,omitempty"`
}

func (m *Message) Message() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return json
}
