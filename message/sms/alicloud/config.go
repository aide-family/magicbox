package alicloud

import (
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Config interface {
	GetSignName() string
	GetEndpoint() string
	GetAccessKeyID() string
	GetAccessKeySecret() string
}

var runtimeOptions = &util.RuntimeOptions{
	Autoretry:   tea.Bool(true),
	MaxAttempts: tea.Int(3),
	IgnoreSSL:   tea.Bool(true),
}
