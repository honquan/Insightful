package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/facebookgo/muster"
	"github.com/gocraft/work"
	"github.com/gorilla/websocket"
	"github.com/jrallison/go-workers"
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"insightful/model"
	"insightful/src/apis/kit/custom_worker"
	"insightful/src/apis/pkg/enum"
	repository "insightful/src/apis/repositories"
	"log"
	"os"
	"time"
	"unsafe"
)

type WebsocketService interface {
	ReaderWithGoWorker(ctx context.Context, conn *websocket.Conn) error
	ReaderWithGoCraft(ctx context.Context, conn *websocket.Conn) error
	ReaderWithCustomWorkerPool(ctx context.Context, conn *websocket.Conn) error
	ReaderWithAntsWorkerPool(ctx context.Context, conn *websocket.Conn) error
	CoordinateWorkerGo(message *workers.Msg)
	CoordinateWorkerCraft(job *work.Job) error
	CoordinateWorkerPool(job custom_worker.Job) error
	CoordinateAntsWorkerPool(i interface{})
}

type websocketService struct {
	insightfullRepo repository.InsightfullRepository

	muster   muster.Client
	enqueuer *work.Enqueuer

	Client *websocketService
	Items  []mongo.WriteModel

	antPool *ants.PoolWithFunc
}

func NewWebsocketService(dispatcher *custom_worker.Dispatcher, insightfullRepo repository.InsightfullRepository, enqueuer *work.Enqueuer) WebsocketService {
	wss := &websocketService{
		insightfullRepo: insightfullRepo,
		enqueuer:        enqueuer,
	}

	// init muster
	wss.muster.MaxBatchSize = 1000
	wss.muster.MaxConcurrentBatches = 10000
	wss.muster.BatchTimeout = 5000 * time.Millisecond
	wss.muster.PendingWorkCapacity = 5000
	wss.muster.BatchMaker = func() muster.Batch {
		return &websocketService{
			Client:          wss,
			insightfullRepo: insightfullRepo,
		}
	}
	err := wss.muster.Start()
	if err != nil {
		panic(err)
	}

	// init custom worker pool
	dispatcher.AppendCallbackWorker(wss.CoordinateWorkerPool)
	dispatcher.Run()

	///////////////////////////////////////
	wss.antPool = wss.initAntsWorker()

	return wss
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////	GO WORKER	//////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *websocketService) ReaderWithGoWorker(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error when ReaderWithGoWorker ReadMessage:", err)
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
		fmt.Println("error when CoordinateWorkerGo message.Args:", err)
		return
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(arr[0].(string))
	var data interface{}
	err = json.Unmarshal(rawDecodedText, &data)
	if err != nil {
		fmt.Println("error when unmarshal:", err)
	}

	// Go ahead and proccess
	s.Push(model.Insightful{
		Mongo: model.Mongo{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: 0,
		},
		Coordinates: data,
	})
	return
}

func (s *websocketService) AddJob(queue string, at time.Time, args ...interface{}) string {
	ts := float64(at.UTC().Unix())
	jid, err := workers.EnqueueWithOptions(queue, "Add", args, workers.EnqueueOptions{Retry: true, RetryCount: 4, At: ts})
	if err != nil {
		fmt.Println("error when EnqueueWithOptions:", err)
	}
	return jid
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////		GO CRAFT	////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *websocketService) ReaderWithGoCraft(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error when ReaderWithGoCraft ReadMessage:", err)
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

//var enqueuer = work.NewEnqueuer(enum.CoordinateNameSpace, &redis.Pool{
//	MaxActive: 5,
//	MaxIdle:   5,
//	Wait:      true,
//	Dial: func() (redis.Conn, error) {
//		return redis.Dial("tcp", fmt.Sprintf(":%v", conf.EnvConfig.RedisPort))
//	},
//})

func (s *websocketService) enqueueJobCraft(job string, payload work.Q) {
	_, err := s.enqueuer.Enqueue(job, payload)
	if err != nil {
		fmt.Println("error when enqueueJobCraft Enqueue:", err)
		log.Fatal(err)
	}
}

func (s *websocketService) CoordinateWorkerCraft(job *work.Job) error {
	// Extract arguments:
	if err := job.ArgError(); err != nil {
		fmt.Println("error when job get arg:", err)
		return err
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(job.Args[enum.GoCraftMessage].(string))
	if err != nil {
		fmt.Println("error when DecodeString:", err)
	}

	var data interface{}
	err = json.Unmarshal(rawDecodedText, &data)
	if err != nil {
		fmt.Println("error when unmarshal:", err)
	}

	// Go ahead and proccess
	s.Push(model.Insightful{
		Mongo: model.Mongo{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: 0,
		},
		Coordinates: data,
	})

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////// CUSTOM WORKER POOL ////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *websocketService) ReaderWithCustomWorkerPool(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error when read message:", err)
			return err
		}

		// Push the work onto the queue.
		custom_worker.Submit(p)
	}
}

func (s *websocketService) CoordinateWorkerPool(job custom_worker.Job) error {
	var data interface{}
	err := json.Unmarshal(job.Payload, &data)
	if err != nil {
		fmt.Println("error when unmarshal:", err)
	}

	// Go ahead and proccess
	s.Push(model.Insightful{
		Mongo: model.Mongo{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: 0,
		},
		Coordinates: data,
	})
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////// ANTS WORKER POOL  /////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *websocketService) initAntsWorker() *ants.PoolWithFunc {
	// Use the common pool.
	//var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(100000, func(i interface{}) {
		s.CoordinateAntsWorkerPool(i)
		//wg.Done()
	})
	//defer p.Release()
	// Submit tasks one by one.
	//for i := 0; i < runTimes; i++ {
	//	wg.Add(1)
	//	_ = p.Invoke(int32(i))
	//}
	//wg.Wait()
	//fmt.Printf("running goroutines: %d\n", p.Running())
	p.Running()

	return p
}

func (s *websocketService) ReaderWithAntsWorkerPool(ctx context.Context, conn *websocket.Conn) error {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error when read message:", err)
			return err
		}

		// Push the work onto the queue.
		err = s.antPool.Invoke(ByteSlice2String(p))
	}
}

func (s *websocketService) CoordinateAntsWorkerPool(i interface{}) {
	//var data interface{}
	//err := json.Unmarshal(i.([]byte), &data)
	//if err != nil {
	//	fmt.Println("error when unmarshal:", err)
	//}

	// Go ahead and proccess
	s.Push(model.Insightful{
		Mongo: model.Mongo{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: 0,
		},
		Coordinates: i,
	})
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////	 MUSTER	 ///////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
	s.Items = append(s.Items, mongo.NewInsertOneModel().SetDocument(item))
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

	err := s.insightfullRepo.BulkWrite(context.Background(), s.Items)
	if err != nil {
		fmt.Println("error when create many mongo:", err)
	}
	os.Stdout.Sync()
}

func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func convert(myBytes []byte) string {
	return string(myBytes[:])
}
