package services

import (
	"context"
	"fmt"
	"github.com/facebookgo/muster"
	"github.com/fasthttp/websocket"
	"github.com/panjf2000/ants/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"insightful/model"
	repository "insightful/src/apis/repositories"
	"log"
	"os"
	"time"
	"unsafe"
)

type WebsocketService interface {
	ReaderFastHttpWithAnts(ctx context.Context, conn *websocket.Conn) error
	CoordinateFastHttpWithAnts(i interface{})
}

type websocketService struct {
	insightfullRepo repository.InsightfullRepository

	muster muster.Client

	Client *websocketService
	Items  []mongo.WriteModel

	antPool *ants.PoolWithFunc
}

func NewWebsocketService(insightfullRepo repository.InsightfullRepository) WebsocketService {
	wss := &websocketService{
		insightfullRepo: insightfullRepo,
	}

	// init muster
	wss.muster.MaxBatchSize = 1000
	wss.muster.MaxConcurrentBatches = 5000
	wss.muster.BatchTimeout = 5000 * time.Millisecond
	wss.muster.PendingWorkCapacity = 3000
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

	///////////////////////////////////////
	wss.antPool = wss.initAntsWorker()

	return wss
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////		FASTHTTP	////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *websocketService) initAntsWorker() *ants.PoolWithFunc {
	// Use the common pool.
	//var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(100000, func(i interface{}) {
		s.CoordinateFastHttpWithAnts(i)
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

func (s *websocketService) ReaderFastHttpWithAnts(ctx context.Context, conn *websocket.Conn) error {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return err
		}

		// Push the work onto the queue.
		err = s.antPool.Invoke(ByteSlice2String(p))
	}
}

func (s *websocketService) CoordinateFastHttpWithAnts(i interface{}) {
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
