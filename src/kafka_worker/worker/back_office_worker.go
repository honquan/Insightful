package worker

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"tcbmerchantsite/pkg/common/queue"
	"tcbmerchantsite/pkg/common/stringutil"
	"tcbmerchantsite/pkg/dto"
	"tcbmerchantsite/pkg/log"
	"tcbmerchantsite/pkg/service"
)

type BackOfficeWorker struct {
	id                int
	q                 chan *queue.ConsumerMessage
	backOfficeService service.BackOfficeService
}

func NewBackOfficeWorker(id int, q chan *queue.ConsumerMessage, boService service.BackOfficeService) *BackOfficeWorker {
	return &BackOfficeWorker{
		id:                id,
		q:                 q,
		backOfficeService: boService,
	}
}

func (b *BackOfficeWorker) Start() error {
	go func() {
		for {
			msg, ok := <-b.q
			if ok {
				var err error
				kafkaMsg := &dto.BackOfficeMessage{}
				if err = json.Unmarshal(msg.Value, kafkaMsg); err != nil {
					log.Errorw(nil, "Error receive message", "err", err)
					continue
				}
				ctx := b.initContext(kafkaMsg)
				defer handlePanic(ctx)
				log.Infow(ctx, "Receive message back office", "msg", kafkaMsg.ToString())

				if kafkaMsg.Payload == nil || kafkaMsg.Payload.ID == 0 || kafkaMsg.Payload.JobType == "" {
					log.Errorw(nil, "Msg kafka is invalid", "kafka", kafkaMsg)
					continue
				}

				jobType := kafkaMsg.Payload.JobType
				// handle import
				if stringutil.Contains(jobType.ToArrayImportType(), string(jobType)) {
					b.backOfficeService.Import(ctx, kafkaMsg)
				}
			} else {
				// Channel closed, exit
				return
			}
		}
	}()

	return nil
}

func (b *BackOfficeWorker) initContext(kafkaMsg *dto.BackOfficeMessage) context.Context {
	ctx := context.Background()
	msgID := kafkaMsg.ID
	if msgID == "" {
		msgID = uuid.New().String()
	}
	ctx = context.WithValue(ctx, log.MessageIDKey, msgID)
	return ctx
}

func handlePanic(ctx context.Context) {
	// detect if panic occurs or not
	err := recover()
	if err != nil {
		log.Errorw(ctx, "[PANIC RECOVERED]", err)
	}
}
