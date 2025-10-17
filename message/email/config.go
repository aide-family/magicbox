package email

type Config interface {
	GetHost() string
	GetPort() int32
	GetUsername() string
	GetPassword() string
}
