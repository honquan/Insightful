package custom_worker

import (
	"github.com/facebookgo/muster"
	"log"
	"time"
)

// The CoordinateClient manages the Coordinate page and dispatches.
type CoordinateClient struct {
	MaxBatchSize        uint          // How much a shopper can carry at a time.
	BatchTimeout        time.Duration // How long we wait once we need to get something.
	PendingWorkCapacity uint          // How long our shopping list can be.
	muster              muster.Client
}

// The CoordinateClient has to be started in order to initialize the underlying
// work channel as well as the background goroutine that handles the work.
func (s *CoordinateClient) Start() error {
	s.muster.MaxBatchSize = s.MaxBatchSize
	s.muster.BatchTimeout = s.BatchTimeout
	s.muster.PendingWorkCapacity = s.PendingWorkCapacity
	s.muster.BatchMaker = func() muster.Batch { return &batch{Client: s} }
	return s.muster.Start()
}

// Similarly the ShoppingClient has to be stopped in order to ensure we flush
// pending items and wait for in progress batches.
func (s *CoordinateClient) Stop() error {
	return s.muster.Stop()
}

// The CoordinateClient provides a typed Add method which enqueues the work.
func (s *CoordinateClient) Add(item interface{}) {
	s.muster.Work <- item
}

// The batch is the collection of items that will be dispatched together.
type batch struct {
	Client *CoordinateClient
	Items  []interface{}
}

// The batch provides an untyped Add to satisfy the muster.Batch interface. As
// is the case here, the Batch implementation is internal to the user of muster
// and not exposed to the users of ShoppingClient.
func (b *batch) Add(item interface{}) {
	b.Items = append(b.Items, item)
}

// Once a Batch is ready, it will be Fired. It must call notifier.Done once the
// batch has been processed.
func (b *batch) Fire(notifier muster.Notifier) {
	//defer notifier.Done()
	log.Println(" ==============================================================================================")
	log.Println(" ==============================================================================================")
	log.Println(" ==============================================================================================")
	log.Println(" ==============================================================================================")
	log.Println(" ==============================================================================================")
	log.Println("Delivery ===================", b.Items)
	//os.Stdout.Sync()
}
