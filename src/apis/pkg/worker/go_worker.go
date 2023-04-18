package worker

import (
	"fmt"
	"github.com/jrallison/go-workers"
	"insightful/src/apis/conf"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func RunGoWorker() {
	// init go worker
	workers.Configure(map[string]string{
		// location of redis instance
		"server": fmt.Sprintf("%v:%v", conf.EnvConfig.RedisHost, conf.EnvConfig.RedisPort),
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "100",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": strconv.Itoa(rand.Intn(10000)),
	})

	// register job types and the function to execute them
	workers.Process("Sample", SampleWorker, 3)   // (queue name, Executor/Worker, concurrency
	workers.Process("Sample2", SampleWorker2, 3) // (queue name, Executor/Worker, concurrency
	go workers.StatsServer(8890)
	workers.Run()
}

func SampleWorker(message *workers.Msg) {
	//time.Sleep(3000 * time.Millisecond)
	_, _ = message.Args().Array()
	//log.Println("Working sample 1 on job, arg: %s, msg: %s", args, message.Jid())
	return
}

func SampleWorker2(message *workers.Msg) {
	//time.Sleep(3000 * time.Millisecond)
	args, _ := message.Args().Array()
	log.Println("Working sample 2 on job, arg: %s, msg: %s", args, message.Jid())
	return
}

func AddJob(queue string, at time.Time, args ...interface{}) string {
	ts := float64(at.UTC().Unix())
	jid, err := workers.EnqueueWithOptions(queue, "Add", args, workers.EnqueueOptions{Retry: true, RetryCount: 4, At: ts})
	if err != nil {

	}
	return jid
}
