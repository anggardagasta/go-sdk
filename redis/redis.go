package redis

import (
	"context"
	"errors"
	"time"

	rdsv9 "github.com/redis/go-redis/v9"
)

// New initialize configuration redis
func New(c Config) (Method, error) {
	var cl rdsv9.UniversalClient

	clCfg := &rdsv9.UniversalOptions{
		Addrs:        c.Hosts,
		Password:     c.Password,
		PoolSize:     c.PoolSize,
		PoolTimeout:  c.PoolTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}

	if c.Mode == "cluster" {
		cl = rdsv9.NewClusterClient(clCfg.Cluster())
	} else {
		cl = rdsv9.NewClient(clCfg.Simple())
	}

	if err := cl.Ping(context.Background()).Err(); err != nil {
		return Module{}, errors.Join(err, errors.New("redis connection failed"))
	}

	return Module{
		client: cl,
	}, nil
}

// Set command redis
func (m Module) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.client.Set(ctx, key, value, expiration).Err()
}

// Get command redis
func (m Module) Get(ctx context.Context, key string) (string, error) {
	val, err := m.client.Get(ctx, key).Result()
	if err == rdsv9.Nil {
		return "", nil
	}
	return val, err
}

// HSet command redis
func (m Module) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return m.client.HSet(ctx, key, field, value).Err()
}

// HSet command redis with Expire
func (m Module) HSetX(ctx context.Context, key string, field string, value interface{}, expiration time.Duration) error {
	err := m.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return m.client.Expire(ctx, key, expiration).Err()
}

// HGet command redis
func (m Module) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := m.client.HGet(ctx, key, field).Result()
	if err == rdsv9.Nil {
		return "", nil
	}
	return val, err
}

// Expire command redis
func (m Module) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return m.client.Expire(ctx, key, expiration).Err()
}

// LPush command redis
func (m Module) LPush(ctx context.Context, key string, values ...interface{}) error {
	return m.client.LPush(ctx, key, values...).Err()
}

// ZRemRangeByScore command redis
func (m Module) ZRemRangeByScore(ctx context.Context, key string, min string, max string) error {
	return m.client.ZRemRangeByScore(ctx, key, min, max).Err()
}

func (m Module) ZRevRangeByScore(ctx context.Context, key string, max, min string, offset, count int64) ([]string, error) {
	return m.client.ZRevRangeByScore(ctx, key, &rdsv9.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (m Module) ZRevRangeByScoreWithScore(ctx context.Context, key string, max, min string, offset, count int64) ([]rdsv9.Z, error) {
	return m.client.ZRevRangeByScoreWithScores(ctx, key, &rdsv9.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	}).Result()
}

func (m Module) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return m.client.ZRem(ctx, key, members...).Err()
}

// ZRange command range
func (m Module) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return m.client.ZRange(ctx, key, start, stop).Result()
}

// HDel command redis
func (m Module) HDel(ctx context.Context, key string, fields ...string) error {
	return m.client.HDel(ctx, key, fields...).Err()
}

// ZAdd command redis
func (m Module) ZAdd(ctx context.Context, key string, members ...*rdsv9.Z) error {
	converted := make([]rdsv9.Z, len(members))
	for i, member := range members {
		if member == nil {
			return errors.New("nil member provided")
		}
		converted[i] = *member
	}
	return m.client.ZAdd(ctx, key, converted...).Err()
}

// HGetAll command redis
func (m Module) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := m.client.HGetAll(ctx, key).Result()
	if err == rdsv9.Nil {
		return map[string]string{}, nil
	}
	return val, err
}

// Ping command pong
func (m Module) Ping(ctx context.Context) (bool, error) {
	res, err := m.client.Ping(ctx).Result()
	if err != nil {
		return false, err
	}

	if res == "PONG" {
		return true, nil
	}

	return false, errors.New("Redis client not receive PONG")
}

// Del command delete key
func (m Module) Del(ctx context.Context, key string) error {
	_, err := m.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (m Module) Pipeline(ctx context.Context) rdsv9.Pipeliner {
	return m.client.Pipeline()
}

func (m Module) TTL(ctx context.Context, key string) (time.Duration, error) {
	return m.client.TTL(ctx, key).Result()
}
