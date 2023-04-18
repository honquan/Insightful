package connection

import (
	"fmt"
	"github.com/go-redis/redis"
	"insightful/src/apis/conf"
	"insightful/src/apis/pkg/redisutil"
)

func NewRedisConnection() (redisutil.Cache, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", conf.EnvConfig.RedisHost, conf.EnvConfig.RedisPort),
			Password: conf.EnvConfig.RedisPassword,
			DB:       conf.EnvConfig.RedisDatabase,
		},
	)

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}

	redisStore := redisutil.NewCache(redisClient)
	return redisStore, nil
}
