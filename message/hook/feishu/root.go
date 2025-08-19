package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

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

const MessageChannelFeishu message.MessageChannel = "webhook-feishu"

type Content struct {
	Text  *Text  `json:"text,omitempty"`
	Post  *Post  `json:"post,omitempty"`
	Card  *Card  `json:"card,omitempty"`
	Image string `json:"image_key,omitempty"`
}

type Message struct {
	MsgType   MessageType `json:"msg_type"`
	Content   *Content    `json:"content,omitempty"`
	Timestamp string      `json:"timestamp"`
	Sign      string      `json:"sign"`
}

func (m *Message) Message(channel message.MessageChannel) ([]byte, error) {
	if err := MessageChannelFeishu.Check(channel); err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (m *Message) Signature(secret string) error {
	if m.Timestamp != "" && m.Sign != "" {
		return nil
	}
	m.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	// timestamp + key sha256, then base64 encode
	signString := m.Timestamp + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(signString))
	_, err := h.Write(data)
	if err != nil {
		return err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	m.Sign = signature
	return nil
}
