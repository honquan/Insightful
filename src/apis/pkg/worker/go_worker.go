package worker

import (
	"fmt"
	"github.com/jrallison/go-workers"
	"insightful/src/apis/conf"
	"insightful/src/apis/pkg/enum"
	"insightful/src/apis/service"
	"math/rand"
	"strconv"
	"time"
)

type MyLogger struct {
}

func (l *MyLogger) Println(v ...interface{}) {
	// noop
}
func (l *MyLogger) Printf(fmt string, v ...interface{}) {
	// noop
}

func RunGoWorker() {
	var websocketService service.WebsocketService
	_ = service.GetServiceContainer().Invoke(func(s service.WebsocketService) {
		websocketService = s
	})

	// init go worker
	workers.Configure(map[string]string{
		// location of redis instance
		"server": fmt.Sprintf("%v:%v", conf.EnvConfig.RedisHost, conf.EnvConfig.RedisPort),
		// instance of the database
		"database": "10",
		// number of connections to keep open with redis
		"pool": "100",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": strconv.Itoa(rand.Intn(10000)),
	})

	//workers.Middleware.Append(&myMiddleware{})
	workers.Logger = &MyLogger{}

	// register job types and the function to execute them
	workers.Process(enum.JobNameCoordinate, websocketService.CoordinateWorker, 100) // (queue name, Executor/Worker, concurrency

	// stats will be available at http://localhost:8890/stats
	go workers.StatsServer(8890)

	// Blocks until process is told to exit via unix signal
	workers.Run()
}

//func CoordinateWorker(message *workers.Msg) {
//	arr, err := message.Args().Array()
//	if err != nil {
//		return
//	}
//
//	rawDecodedText, err := base64.StdEncoding.DecodeString(arr[0].(string))
//	var data interface{}
//	err = json.Unmarshal(rawDecodedText, &data)
//	if err != nil {
//		fmt.Println("error:", err)
//	}
//
//	// Go ahead and proccess
//	SM.Add(data)
//	return
//}

func AddJob(queue string, at time.Time, args ...interface{}) string {
	ts := float64(at.UTC().Unix())
	jid, err := workers.EnqueueWithOptions(queue, "Add", args, workers.EnqueueOptions{Retry: true, RetryCount: 4, At: ts})
	if err != nil {

	}
	return jid
}
