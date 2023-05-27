package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/jrallison/go-workers"
	"insightful/src/apis/conf"
	"insightful/src/apis/kit/custom_worker"
	"insightful/src/apis/pkg/enum"
	"log"
	"time"
)

type WebsocketService interface {
	ReaderWithGoWorker(ctx context.Context, conn *websocket.Conn) error
	ReaderWithGoCraft(ctx context.Context, conn *websocket.Conn) error
	CoordinateWorker(message *workers.Msg)
}

type websocketService struct{}

var sm = &custom_worker.CoordinateClient{
	MaxBatchSize:        20,
	BatchTimeout:        2000 * time.Millisecond,
	PendingWorkCapacity: 20,
}

var RedisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", fmt.Sprintf(":%v", conf.EnvConfig.RedisPort))
	},
}

func NewWebsocketService() WebsocketService {
	if err := sm.Start(); err != nil {
		log.Printf("Error when start muster: ", err)
	}

	return &websocketService{}
}

var enqueuer = work.NewEnqueuer(enum.CoordinateNameSpace, RedisPool)

func enqueueJobCraft(job string, payload work.Q) {
	_, err := enqueuer.Enqueue(job, payload)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *websocketService) ReaderWithGoWorker(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		s.AddJob(enum.JobNameCoordinate, time.Now().UTC(), p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			return err
		}

	}
}

func (s *websocketService) ReaderWithGoCraft(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		// enqueue go craft
		enqueueJobCraft(
			enum.JobNameCoordinate,
			work.Q{enum.GoCraftMessage: p},
		)

		if err := conn.WriteMessage(messageType, p); err != nil {
			return err
		}

	}
}

func (s *websocketService) CoordinateWorker(message *workers.Msg) {
	arr, err := message.Args().Array()
	if err != nil {
		return
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(arr[0].(string))
	var data interface{}
	err = json.Unmarshal(rawDecodedText, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Go ahead and proccess
	sm.Add(data)
	return
}

func (s *websocketService) AddJob(queue string, at time.Time, args ...interface{}) string {
	ts := float64(at.UTC().Unix())
	jid, err := workers.EnqueueWithOptions(queue, "Add", args, workers.EnqueueOptions{Retry: true, RetryCount: 4, At: ts})
	if err != nil {

	}
	return jid
}
