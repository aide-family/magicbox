package other

import "github.com/aide-family/magicbox/httpx"

type Config interface {
	GetURL() string
	GetHeaders() map[string][]string
	GetBasicAuth() *httpx.BasicAuth
}
