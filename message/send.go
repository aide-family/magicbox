package message

import "context"

type Sender interface {
	Send(ctx context.Context, message Message) error
}

type Driver interface {
	New() (Sender, error)
}

func NewSender(driver Driver) (Sender, error) {
	return driver.New()
}
