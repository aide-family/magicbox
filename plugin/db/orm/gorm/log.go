package gorm

import (
	"context"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"

	"github.com/aide-family/magicbox/log"
)

// NewLogger creates a new gorm logger.
func NewLogger(logger log.Interface) logger.Interface {
	return &gormLogger{helper: klog.NewHelper(logger)}
}

type gormLogger struct {
	helper *klog.Helper
}

// Error implements logger.Interface.
func (g *gormLogger) Error(ctx context.Context, msg string, args ...any) {
	g.helper.WithContext(ctx).Error(append([]any{"msg", msg}, args...)...)
}

// Info implements logger.Interface.
func (g *gormLogger) Info(ctx context.Context, msg string, args ...any) {
	g.helper.WithContext(ctx).Info(append([]any{"msg", msg}, args...)...)
}

// LogMode implements logger.Interface.
func (g *gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return g
}

// Trace implements logger.Interface.
func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	duration := time.Since(begin)
	if err != nil {
		g.helper.WithContext(ctx).Errorw("begin", begin, "sql", sql, "rowsAffected", rowsAffected, "err", err, "duration", duration)
	} else {
		g.helper.WithContext(ctx).Debugw("begin", begin, "sql", sql, "rowsAffected", rowsAffected, "duration", duration)
	}
}

// Warn implements logger.Interface.
func (g *gormLogger) Warn(ctx context.Context, msg string, args ...any) {
	g.helper.WithContext(ctx).Warn(append([]any{"msg", msg}, args...)...)
}
