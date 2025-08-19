package wechat

import (
	"encoding/json"

	"github.com/aide-family/magicbox/message"
)

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

type TextMessage struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

func NewTextMessage(content string) *TextMessage {
	return &TextMessage{
		Content: content,
	}
}

func (t *TextMessage) WithMentionedList(list []string) *TextMessage {
	t.MentionedList = list
	return t
}

func (t *TextMessage) WithMentionedMobileList(list []string) *TextMessage {
	t.MentionedMobileList = list
	return t
}

func (t *TextMessage) Message() message.Message {
	return &Message{
		MsgType: MessageTypeText,
		Text:    t,
	}
}

type MarkdownMessage struct {
	Content string `json:"content"`
}

func NewMarkdownMessage(content string) *MarkdownMessage {
	return &MarkdownMessage{
		Content: content,
	}
}

func (m *MarkdownMessage) Message() message.Message {
	return &Message{
		MsgType:  MessageTypeMarkdown,
		Markdown: m,
	}
}

type MarkdownV2Message struct {
	Content string `json:"content"`
}

func NewMarkdownV2Message(content string) *MarkdownV2Message {
	return &MarkdownV2Message{
		Content: content,
	}
}

func (m *MarkdownV2Message) Message() message.Message {
	return &Message{
		MsgType:    MessageTypeMarkdownV2,
		MarkdownV2: m,
	}
}

type ImageMessage struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

func NewImageMessage(base64 string, md5 string) *ImageMessage {
	return &ImageMessage{
		Base64: base64,
		MD5:    md5,
	}
}

func (m *ImageMessage) Message() message.Message {
	return &Message{
		MsgType: MessageTypeImage,
		Image:   m,
	}
}
