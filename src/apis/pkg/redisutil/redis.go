package redisutil

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

// Cache provides an access to Redis store.
type Cache interface {
	Set(key string, value interface{}, expireTime int64) error
	SetNX(key string, value interface{}, expireTime int64) (bool, error)
	Get(key string, value interface{}) (interface{}, error)
	Del(keys ...string) (int64, error)
	Expire(key string, expire int64) (bool, error)
	Ping() error
	Exist(key string) (bool, error)
	Incr(key string) (int64, error)
	TTL(key string) (int64, error)
	Keys(pattern string) ([]string, error)
}

type cache struct {
	client *redis.Client
}

// NewCache returns a new instance of Store.
func NewCache(client *redis.Client) Cache {
	return &cache{client}
}

// Del deletes the elements with the specified key.
func (s *cache) Del(keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, nil
	}
	return s.client.Del(keys...).Result()
}

// Set sets the new element with specified to Store.
func (s *cache) Set(key string, value interface{}, expired int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = s.client.Set(key, data, time.Duration(expired)*time.Second).Result()
	return err
}

// SetNX sets if exist key the new element with specified to Store.
func (s *cache) SetNX(key string, value interface{}, expired int64) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	keyNotExisted, err := s.client.SetNX(key, data, time.Duration(expired)*time.Second).Result()
	return keyNotExisted, err
}

// Get get element with specified to Cache.
func (s *cache) Get(key string, value interface{}) (interface{}, error) {
	data, err := s.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data), value)
	return value, err
}

// Expire sets expiration time for the element with specified key.
func (s *cache) Expire(key string, expire int64) (bool, error) {
	value, err := s.client.Expire(key, time.Duration(expire)*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return value, nil
}

// Ping checks the Redis connection.
func (s *cache) Ping() error {
	if _, err := s.client.Ping().Result(); err != nil {
		return err
	}
	return nil
}

// Exist check key exist.
func (s *cache) Exist(key string) (bool, error) {
	value, err := s.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return value == 1, nil
}

// Incr value of key
func (s *cache) Incr(key string) (int64, error) {
	result, err := s.client.Incr(key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// TTL returns the time to live in seconds
func (s *cache) TTL(key string) (int64, error) {
	result, err := s.client.TTL(key).Result()
	if err != nil {
		return 0, err
	}
	return int64(result.Seconds()), nil
}

// Keys return all keys with the specified pattern
func (s *cache) Keys(pattern string) ([]string, error) {
	result, err := s.client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
