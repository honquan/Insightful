// Package rediscmd provides functionality related to Redis caching database.
package rediscmd

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"

	"github.com/pkg/errors"
)

// New returns a new instance of Client.
func New(addr, password string) (Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	if err := cli.Ping().Err(); err != nil {
		return nil, errors.Wrap(err, "pinging redis server")
	}

	return &client{cli}, nil
}

// Client represents redis client.
type Client interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiresTimes ...int) error
	Get(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Incr(ctx context.Context, key string) (int64, error)
	TTL(ctx context.Context, key string) (int64, error)
	ExpireAt(ctx context.Context, key string, tm time.Time) error
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Exists(ctx context.Context, keys ...string) (bool, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, value map[string]interface{}) error
}

type client struct {
	cli *redis.Client
}

func (c *client) Ping(ctx context.Context) error {
	return c.cli.Ping().Err()
}

func (c *client) Set(ctx context.Context, key string, value interface{}, expiresTimes ...int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "decoding data")
	}

	expiresTime := 0 // Zero expiration means the key has no expiration time.
	if len(expiresTimes) > 0 && expiresTimes[0] > 0 {
		expiresTime = expiresTimes[0]
	}

	if err := c.cli.Set(key, data, time.Duration(expiresTime)*time.Second).Err(); err != nil {
		return errors.Wrap(err, "setting redis key")
	}
	return nil
}

func (c *client) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.cli.Get(key).Result()
	if err != nil {
		return errors.Wrap(err, "getting redis key")
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return errors.Wrap(err, "decoding value")
	}

	return nil
}

func (c *client) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.cli.Del(keys...).Result()
}

func (c *client) Incr(ctx context.Context, key string) (int64, error) {
	result, err := c.cli.Incr(key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "increasing redis key")
	}
	return result, nil
}

func (c *client) TTL(ctx context.Context, key string) (int64, error) {
	result, err := c.cli.TTL(key).Result()
	if err != nil {
		return 0, errors.Wrap(err, "getting ttl redis key")
	}
	return int64(result.Seconds()), nil
}

func (c *client) ExpireAt(ctx context.Context, key string, tm time.Time) error {
	if err := c.cli.ExpireAt(key, tm).Err(); err != nil {
		return errors.Wrap(err, "setting expire at")
	}
	return nil
}

func (c *client) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.cli.Expire(key, expiration).Result()
}

func (c *client) Exists(ctx context.Context, keys ...string) (bool, error) {
	countExists, err := c.cli.Exists(keys...).Result()
	if err != nil {
		return false, errors.Wrap(err, "finding existing keys")
	}
	return countExists > 0, nil
}

func (c *client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	data, err := c.cli.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "HGetAll")
	}
	return data, nil
}

func (c *client) HSet(ctx context.Context, key string, value map[string]interface{}) error {
	return c.cli.HSet("hash", key, value).Err()
}
