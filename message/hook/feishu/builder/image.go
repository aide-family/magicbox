package builder

import "github.com/aide-family/magicbox/message"

func NewImageMessage(imageKey string) message.Message {
	return &Message{
		MsgType: MessageTypeImage,
		Content: &Content{
			Image: imageKey,
		},
	}
}
