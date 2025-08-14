// Package redis is a cache driver for redis.
package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/aide-family/magicbox/plugin/cache"
)

var _ cache.Driver = (*initializer)(nil)
var _ cache.Interface = (*redisCache)(nil)

// CacheDriver returns a new cache driver.
func CacheDriver(cli *redis.Client) cache.Driver {
	return &initializer{cli: cli}
}

type initializer struct {
	cli *redis.Client
}

// New implements cache.Driver.
func (i *initializer) New(ctx context.Context) (cache.Interface, error) {
	return &redisCache{
		cli: i.cli,
	}, nil
}

type redisCache struct {
	cli *redis.Client
}

// Close implements cache.Interface.
func (r *redisCache) Close() error {
	return r.cli.Close()
}

// Del implements cache.Interface.
func (r *redisCache) Del(ctx context.Context, key string) error {
	return r.cli.Del(ctx, key).Err()
}

// Exists implements cache.Interface.
func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	res, err := r.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

// Get implements cache.Interface.
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.cli.Get(ctx, key).Result()
}

// HDel implements cache.Interface.
func (r *redisCache) HDel(ctx context.Context, key string, field string) error {
	return r.cli.HDel(ctx, key, field).Err()
}

// HExists implements cache.Interface.
func (r *redisCache) HExists(ctx context.Context, key string, field string) (bool, error) {
	res, err := r.cli.HExists(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return res, nil
}

// HGet implements cache.Interface.
func (r *redisCache) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.cli.HGet(ctx, key, field).Result()
}

// HMGet implements cache.Interface.
func (r *redisCache) HMGet(ctx context.Context, key string, fields ...string) ([][]byte, error) {
	res, err := r.cli.HMGet(ctx, key, fields...).Result()
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
func (r *redisCache) HMSet(ctx context.Context, key string, fields map[string]string) error {
	return r.cli.HMSet(ctx, key, fields).Err()
}

// HSet implements cache.Interface.
func (r *redisCache) HSet(ctx context.Context, key string, field string, value string) error {
	return r.cli.HSet(ctx, key, field, value).Err()
}

// IncMax implements cache.Interface.
func (r *redisCache) IncMax(ctx context.Context, key string, max int, ttl time.Duration) (bool, error) {
	res, err := r.cli.Eval(ctx, `
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
func (r *redisCache) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.cli.SetNX(ctx, key, 1, ttl).Result()
}

// Set implements cache.Interface.
func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.cli.Set(ctx, key, value, ttl).Err()
}

// Unlock implements cache.Interface.
func (r *redisCache) Unlock(ctx context.Context, key string) error {
	return r.cli.Del(ctx, key).Err()
}

// ZAdd implements cache.Interface.
func (r *redisCache) ZAdd(ctx context.Context, key string, score float64, member string) error {
	return r.cli.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// ZRange implements cache.Interface.
func (r *redisCache) ZRange(ctx context.Context, key string, start int, stop int) ([]string, error) {
	return r.cli.ZRange(ctx, key, int64(start), int64(stop)).Result()
}

// ZRangeByScore implements cache.Interface.
func (r *redisCache) ZRangeByScore(ctx context.Context, key string, min float64, max float64) ([]string, error) {
	return r.cli.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    strconv.FormatFloat(min, 'f', -1, 64),
		Max:    strconv.FormatFloat(max, 'f', -1, 64),
		Offset: 0,
		Count:  0,
	}).Result()
}

// ZRem implements cache.Interface.
func (r *redisCache) ZRem(ctx context.Context, key string, member string) error {
	return r.cli.ZRem(ctx, key, member).Err()
}

// ZRemRangeByScore implements cache.Interface.
func (r *redisCache) ZRemRangeByScore(ctx context.Context, key string, min float64, max float64) error {
	return r.cli.ZRemRangeByScore(ctx, key, strconv.FormatFloat(min, 'f', -1, 64), strconv.FormatFloat(max, 'f', -1, 64)).Err()
}
