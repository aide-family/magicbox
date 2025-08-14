// Package cache is a cache plugin.
package cache

import (
	"context"
	"time"
)

// Interface is a cache interface.
type Interface interface {
	// Close closes the cache.
	Close() error

	// Get gets the value of the key.
	Get(ctx context.Context, key string) (string, error)

	// Set sets the value of the key.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Del deletes the key.
	Del(ctx context.Context, key string) error

	// Exists checks if the key exists.
	Exists(ctx context.Context, key string) (bool, error)

	// IncMax increments the value of the key by the given amount.
	IncMax(ctx context.Context, key string, max int) (int, error)

	// Lock locks the key.
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)

	// Unlock unlocks the key.
	Unlock(ctx context.Context, key string) error

	// HSet sets the value of the field in the hash.
	HSet(ctx context.Context, key string, field string, value string) error

	// HGet gets the value of the field in the hash.
	HGet(ctx context.Context, key string, field string) (string, error)

	// HDel deletes the field in the hash.
	HDel(ctx context.Context, key string, field string) error

	// HExists checks if the field exists in the hash.
	HExists(ctx context.Context, key string, field string) (bool, error)

	// HMSet sets the values of the fields in the hash.
	HMSet(ctx context.Context, key string, fields map[string]string) error

	// HMGet gets the values of the fields in the hash.
	HMGet(ctx context.Context, key string, fields ...string) ([]string, error)

	// HMDel deletes the fields in the hash.
	HMDel(ctx context.Context, key string, fields ...string) error

	// ZAdd adds a member to the sorted set.
	ZAdd(ctx context.Context, key string, score float64, member string) error

	// ZRange gets the members of the sorted set.
	ZRange(ctx context.Context, key string, start int, stop int) ([]string, error)

	// ZRangeByScore gets the members of the sorted set by score.
	ZRangeByScore(ctx context.Context, key string, min float64, max float64) ([]string, error)

	// ZRem removes a member from the sorted set.
	ZRem(ctx context.Context, key string, member string) error

	// ZRemRangeByScore removes the members of the sorted set by score.
	ZRemRangeByScore(ctx context.Context, key string, min float64, max float64) error
}

// Driver is a driver for the cache.
type Driver interface {
	// New creates a new cache.
	New(ctx context.Context) (Interface, error)
}

// New creates a new cache.
func New(ctx context.Context, driver Driver) (Interface, error) {
	return driver.New(ctx)
}
