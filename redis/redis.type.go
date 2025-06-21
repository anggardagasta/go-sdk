package redis

import (
	"context"
	"time"

	rdsv9 "github.com/redis/go-redis/v9"
)

type Module struct {
	client rdsv9.UniversalClient
}

type Config struct {
	Hosts        []string      `yaml:"hosts" json:"hosts"`
	Mode         string        `yaml:"mode" json:"mode"`
	Password     string        `yaml:"password" json:"password"`
	PoolSize     int           `yaml:"pool_size" json:"pool_size"`
	PoolTimeout  time.Duration `yaml:"pool_timeout" json:"pool_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	Expiration   time.Duration `yaml:"expiration" json:"expiration"`
}

// Method interface list of functional to allow
type Method interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HSetX(ctx context.Context, key string, field string, value interface{}, expiration time.Duration) error
	HGet(ctx context.Context, key, field string) (string, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	LPush(ctx context.Context, key string, values ...interface{}) error
	ZRemRangeByScore(ctx context.Context, key string, min string, max string) error
	ZRevRangeByScore(ctx context.Context, key string, max, min string, offset, count int64) ([]string, error)
	ZRevRangeByScoreWithScore(ctx context.Context, key string, max, min string, offset, count int64) ([]rdsv9.Z, error)
	ZRem(ctx context.Context, key string, members ...interface{}) error
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	HDel(ctx context.Context, key string, fields ...string) error
	ZAdd(ctx context.Context, key string, members ...*rdsv9.Z) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Ping(ctx context.Context) (bool, error)
	Del(ctx context.Context, key string) error
	Pipeline(ctx context.Context) rdsv9.Pipeliner
	TTL(ctx context.Context, key string) (time.Duration, error)
}
