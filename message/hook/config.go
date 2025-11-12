package hook

type Config interface {
	GetURL() string
	GetSecret() string
}
