// Package wechat provides a set of message types for WeChat.
package wechat

import (
	"github.com/aide-family/magicbox/message"
	"github.com/aide-family/magicbox/serialize"
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

func (m *Message) Message(channel message.MessageChannel) ([]byte, error) {
	if err := MessageChannelWechat.Check(channel); err != nil {
		return nil, err
	}
	jsonBytes, err := serialize.JSONMarshal(m)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
