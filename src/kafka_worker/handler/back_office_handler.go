package handler

import (
	"github.com/samber/do"
	"tcbmerchantsite/cmd/kafka_worker/worker"
	"tcbmerchantsite/pkg/common/queue"
	"tcbmerchantsite/pkg/config"
	"tcbmerchantsite/pkg/log"
	"tcbmerchantsite/pkg/service"
)

type BackOfficeHandler struct {
	q chan *queue.ConsumerMessage
}

func NewBackOfficeHandler(di *do.Injector) (*BackOfficeHandler, error) {
	backOfficeHandler := &BackOfficeHandler{
		q: make(chan *queue.ConsumerMessage),
	}

	// invoke back office service
	backOfficeService := do.MustInvoke[service.BackOfficeService](di)

	// invoke config kafka
	confKafka := do.MustInvoke[config.Kafka](di)

	// run
	for i := 0; i < confKafka.NumberWorkerRoutine; i++ {
		log.Infow(nil, "Creating kafka worker", "i", i+1)
		backOfficeWorker := worker.NewBackOfficeWorker(i, backOfficeHandler.q, backOfficeService)
		backOfficeWorker.Start()
	}

	log.Infow(nil, "Created workers", "count", confKafka.NumberWorkerRoutine)
	return backOfficeHandler, nil
}

func (b *BackOfficeHandler) HandlerFunc(msg *queue.ConsumerMessage) {
	b.q <- msg
}

func (b *BackOfficeHandler) Close() {
	close(b.q)
}
