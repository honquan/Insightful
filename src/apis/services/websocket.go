package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/facebookgo/muster"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/jrallison/go-workers"
	"insightful/model"
	"insightful/src/apis/conf"
	"insightful/src/apis/pkg/enum"
	repository "insightful/src/apis/repositories"
	"log"
	"time"
)

type WebsocketService interface {
	ReaderWithGoWorker(ctx context.Context, conn *websocket.Conn) error
	ReaderWithGoCraft(ctx context.Context, conn *websocket.Conn) error
	CoordinateWorkerGo(message *workers.Msg)
	CoordinateWorkerCraft(job *work.Job) error
}

type websocketService struct {
	insightfullRepo repository.InsightfullRepository

	muster muster.Client

	Client *websocketService
	Items  []interface{}
}

func NewWebsocketService(insightfullRepo repository.InsightfullRepository) WebsocketService {
	wss := &websocketService{
		insightfullRepo: insightfullRepo,
	}
	wss.muster.MaxBatchSize = 100
	wss.muster.BatchTimeout = 5000 * time.Millisecond
	wss.muster.PendingWorkCapacity = 100
	wss.muster.BatchMaker = func() muster.Batch {
		return &websocketService{
			Client:          wss,
			insightfullRepo: insightfullRepo,
		}
	}
	err := wss.muster.Start()
	if err != nil {

	}

	return wss
}

func (s *websocketService) ReaderWithGoWorker(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		s.AddJob(enum.JobNameCoordinate, time.Now(), p)
		//_, err = workers.EnqueueWithOptions(enum.JobNameCoordinate, "Add", p, workers.EnqueueOptions{Retry: true, RetryCount: 4, At: float64(time.Now().UTC().Unix())})
		//if err != nil {
		//
		//}

		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	return err
		//}

	}
}

func (s *websocketService) CoordinateWorkerGo(message *workers.Msg) {
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
	s.Push(model.Insightful{
		Mongo: model.Mongo{
			CreatedAt: time.Now(),
		},
		Done:      0,
		Coodiates: data,
	})
	return
}

func (s *websocketService) AddJob(queue string, at time.Time, args ...interface{}) string {
	ts := float64(at.UTC().Unix())
	jid, err := workers.EnqueueWithOptions(queue, "Add", args, workers.EnqueueOptions{Retry: true, RetryCount: 4, At: ts})
	if err != nil {

	}
	return jid
}

// Similarly the ShoppingClient has to be stopped in order to ensure we flush
// pending items and wait for in progress batches.
func (s *websocketService) Stop() error {
	return s.muster.Stop()
}

// The CoordinateClient provides a typed Add method which enqueues the work.
func (s *websocketService) Push(item interface{}) {
	s.muster.Work <- item
}

// The batch provides an untyped Add to satisfy the muster.Batch interface. As
// is the case here, the Batch implementation is internal to the user of muster
// and not exposed to the users of ShoppingClient.
func (s *websocketService) Add(item interface{}) {
	s.Items = append(s.Items, item)
}

// Once a Batch is ready, it will be Fired. It must call notifier.Done once the
// batch has been processed.
func (s *websocketService) Fire(notifier muster.Notifier) {
	defer notifier.Done()
	//log.Println(" ==============================================================================================")
	//log.Println(" ==============================================================================================")
	//log.Println(" ==============================================================================================")
	//log.Println(" ==============================================================================================")
	//log.Println(" ==============================================================================================")
	//log.Println("Delivery websocket ===================", s.Items)
	err := s.insightfullRepo.CreateMany(context.Background(), s.Items)
	if err != nil {

	}
	//os.Stdout.Sync()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *websocketService) ReaderWithGoCraft(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		// enqueue go craft
		s.enqueueJobCraft(
			enum.JobNameCoordinate,
			work.Q{enum.GoCraftMessage: p},
		)

		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	return err
		//}

	}
}

var enqueuer = work.NewEnqueuer(enum.CoordinateNameSpace, &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", fmt.Sprintf(":%v", conf.EnvConfig.RedisPort))
	},
})

func (s *websocketService) enqueueJobCraft(job string, payload work.Q) {
	_, err := enqueuer.Enqueue(job, payload)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *websocketService) CoordinateWorkerCraft(job *work.Job) error {
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
	s.Push(model.Insightful{
		Mongo: model.Mongo{
			CreatedAt: time.Now(),
		},
		Done:      0,
		Coodiates: data,
	})

	return nil
}
