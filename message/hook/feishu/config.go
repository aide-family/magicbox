package feishu

type Config interface {
	GetURL() string
	GetSecret() string
	GetKey() string
}
