package worker

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"insightful/src/apis/conf"
	"insightful/src/apis/pkg/enum"
	"insightful/src/apis/services"
	"os"
	"os/signal"
)

type Context struct {
	customerID int64
}

func RunGoCraft() {
	var websocketService services.WebsocketService
	_ = services.GetServiceContainer().Invoke(func(s services.WebsocketService) {
		websocketService = s
	})

	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "CoordinateNameSpace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 100, enum.CoordinateNameSpace, &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf(":%v", conf.EnvConfig.RedisPort))
		},
	})

	// Map the name of jobs to handler functions
	pool.Job(enum.JobNameCoordinate, websocketService.CoordinateWorkerCraft)

	// Customize options:
	//pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).Export)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) saveCoordinate(job *work.Job) error {
	// Extract arguments:
	if err := job.ArgError(); err != nil {
		return err
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(job.Args[enum.GoCraftMessage].(string))

	var data interface{}
	err = json.Unmarshal(rawDecodedText, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Go ahead and proccess
	//SM.Add(data)

	return nil
}
