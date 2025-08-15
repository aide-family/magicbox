package wechat

type Config interface {
	GetURL() string
	GetSecret() string
}
