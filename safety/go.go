package safety

import (
	"context"
	"time"

	"github.com/aide-family/magicbox/log"
	klog "github.com/go-kratos/kratos/v2/log"
)

func Go(ctx context.Context, name string, f func(context.Context) error, logger log.Interface) {
	helper := klog.NewHelper(klog.With(logger, "func", name))

	start := time.Now()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				helper.Errorf("panic: %v", r)
			}
			helper.Infof("completed, cost: %v", time.Since(start))
		}()

		if err := f(ctx); err != nil {
			helper.Errorf("run error: %v", err)
		}
	}()
}
