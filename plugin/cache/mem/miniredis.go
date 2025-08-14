// Package mem is a cache driver for memory.
package mem

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/aide-family/magicbox/plugin/cache"
)

var _ cache.Driver = (*initializer)(nil)
var _ cache.Interface = (*miniRedisCache)(nil)

// CacheDriver returns a new cache driver.
func CacheDriver() cache.Driver {
	return &initializer{}
}

type initializer struct{}

// New implements cache.Driver.
func (i *initializer) New(ctx context.Context) (cache.Interface, error) {
	cli, err := miniredis.Run()
	if err != nil {
		return nil, err
	}
	return &miniRedisCache{
		cli: redis.NewClient(&redis.Options{
			Addr: cli.Addr(),
		}),
	}, nil
}

type miniRedisCache struct {
	cli *redis.Client
}

// Close implements cache.Interface.
func (m *miniRedisCache) Close() error {
	return m.cli.Close()
}

// Del implements cache.Interface.
func (m *miniRedisCache) Del(ctx context.Context, key string) error {
	return m.cli.Del(ctx, key).Err()
}

// Exists implements cache.Interface.
func (m *miniRedisCache) Exists(ctx context.Context, key string) (bool, error) {
	res, err := m.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

// Get implements cache.Interface.
func (m *miniRedisCache) Get(ctx context.Context, key string) (string, error) {
	return m.cli.Get(ctx, key).Result()
}

// Set implements cache.Interface.
func (m *miniRedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return m.cli.Set(ctx, key, value, ttl).Err()
}

// HDel implements cache.Interface.
func (m *miniRedisCache) HDel(ctx context.Context, key string, field string) error {
	return m.cli.HDel(ctx, key, field).Err()
}

// HExists implements cache.Interface.
func (m *miniRedisCache) HExists(ctx context.Context, key string, field string) (bool, error) {
	res, err := m.cli.HExists(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return res, nil
}

// HGet implements cache.Interface.
func (m *miniRedisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	return m.cli.HGet(ctx, key, field).Result()
}

// HMGet implements cache.Interface.
func (m *miniRedisCache) HMGet(ctx context.Context, key string, fields ...string) ([][]byte, error) {
	res, err := m.cli.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	resStr := make([][]byte, 0, len(res))
	for _, v := range res {
		switch val := v.(type) {
		case string:
			resStr = append(resStr, []byte(val))
		case []byte:
			resStr = append(resStr, val)
		case nil:
			resStr = append(resStr, []byte{})
		default:
			return nil, fmt.Errorf("invalid type: %T", val)
		}
	}
	return resStr, nil
}

// HMSet implements cache.Interface.
func (m *miniRedisCache) HMSet(ctx context.Context, key string, fields map[string]string) error {
	return m.cli.HMSet(ctx, key, fields).Err()
}

// HSet implements cache.Interface.
func (m *miniRedisCache) HSet(ctx context.Context, key string, field string, value string) error {
	return m.cli.HSet(ctx, key, field, value).Err()
}

// IncMax implements cache.Interface.
func (m *miniRedisCache) IncMax(ctx context.Context, key string, max int, ttl time.Duration) (bool, error) {
	res, err := m.cli.Eval(ctx, `
		local key = KEYS[1]
		local max = tonumber(ARGV[1])
		local expire = tonumber(ARGV[2])
		local current = tonumber(redis.call("get", key))
		if current == nil then
			redis.call("set", key, 1)
			redis.call("expire", key, expire)
			return 1
		end
		if current >= max then
			return current
		end
		redis.call("incr", key)
		return current
	`, []string{key}, max, int(ttl.Seconds())).Int()
	if err != nil {
		return false, err
	}
	return res < max, nil
}

// Lock implements cache.Interface.
func (m *miniRedisCache) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return m.cli.SetNX(ctx, key, 1, ttl).Result()
}

// Unlock implements cache.Interface.
func (m *miniRedisCache) Unlock(ctx context.Context, key string) error {
	return m.cli.Del(ctx, key).Err()
}

// ZAdd implements cache.Interface.
func (m *miniRedisCache) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return m.cli.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// ZRange implements cache.Interface.
func (m *miniRedisCache) ZRange(ctx context.Context, key string, start int, stop int) ([]string, error) {
	return m.cli.ZRange(ctx, key, int64(start), int64(stop)).Result()
}

// ZRangeByScore implements cache.Interface.
func (m *miniRedisCache) ZRangeByScore(ctx context.Context, key string, min float64, max float64) ([]string, error) {
	return m.cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    strconv.FormatFloat(min, 'f', -1, 64),
		Max:    strconv.FormatFloat(max, 'f', -1, 64),
		Offset: 0,
		Count:  0,
	}).Result()
}

// ZRem implements cache.Interface.
func (m *miniRedisCache) ZRem(ctx context.Context, key string, member string) error {
	return m.cli.ZRem(ctx, key, member).Err()
}

// ZRemRangeByScore implements cache.Interface.
func (m *miniRedisCache) ZRemRangeByScore(ctx context.Context, key string, min float64, max float64) error {
	return m.cli.ZRemRangeByScore(ctx, key, strconv.FormatFloat(min, 'f', -1, 64), strconv.FormatFloat(max, 'f', -1, 64)).Err()
}
