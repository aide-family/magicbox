package builder

import "encoding/json"

type MessageType string

const (
	MessageTypeText       MessageType = "text"
	MessageTypeMarkdown   MessageType = "markdown"
	MessageTypeMarkdownV2 MessageType = "markdown_v2"
	MessageTypeImage      MessageType = "image"
)

type Message struct {
	MsgType    MessageType        `json:"msgtype"`
	Text       *TextMessage       `json:"text,omitempty"`
	Markdown   *MarkdownMessage   `json:"markdown,omitempty"`
	MarkdownV2 *MarkdownV2Message `json:"markdown_v2,omitempty"`
	Image      *ImageMessage      `json:"image,omitempty"`
}

func (m *Message) Message() []byte {
	json, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	return json
}
