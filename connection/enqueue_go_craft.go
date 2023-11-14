package connection

import (
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"insightful/common/config"
	"insightful/src/apis/pkg/enum"
)

func InitEnqueueGoCraft(conf *config.Config) *work.Enqueuer {
	var enqueuer = work.NewEnqueuer(enum.CoordinateNameSpace, &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf(":%v", conf.Redis.RedisPort))
		},
	})

	return enqueuer
}
