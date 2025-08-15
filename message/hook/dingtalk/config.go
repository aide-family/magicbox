package dingtalk

type Config interface {
	GetURL() string
	GetSecret() string
	GetKey() string
}
