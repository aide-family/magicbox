package sugared

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Config interface {
	GetLevel() log.Level
	GetFormat() Formatter
	GetOutput() string
	GetEnableCaller() bool
	GetEnableColor() bool
	GetEnableStack() bool
	IsDev() bool
}

type Formatter string

const (
	FormatterConsole Formatter = "console"
	FormatterJSON              = "json"
)
