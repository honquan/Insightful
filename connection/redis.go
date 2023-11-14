package connection

import (
	"context"
	"insightful/common/adapters/rediscmd"
	"insightful/common/config"
	"log"
	"time"
)

func InitRedis(conf *config.Config) rediscmd.Client {
	redisClient, err := rediscmd.New(conf.Redis.Address, conf.Redis.RedisPassword)
	if err != nil {
		log.Println("Error when initialize Redis, err: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = redisClient.Ping(ctx)
	if err != nil {
		log.Println("Failed to ping Redis", err)
	}

	return redisClient
}
