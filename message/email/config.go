package email

type Config interface {
	GetHost() string
	GetPort() int
	GetUsername() string
	GetPassword() string
}
