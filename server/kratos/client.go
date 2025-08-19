package kratos

import (
	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/wrr"
)

func init() {
	selector.SetGlobalSelector(wrr.NewBuilder())
}
