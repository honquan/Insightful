package cmd

import (
	"github.com/samber/do"
	"github.com/spf13/cobra"
	"tcbmerchantsite/cmd/kafka_worker/handler"
	"tcbmerchantsite/pkg/adapter/mtp"
	"tcbmerchantsite/pkg/adapter/tpay_ops"
	"tcbmerchantsite/pkg/common/language_message"
	"tcbmerchantsite/pkg/common/queue"
	"tcbmerchantsite/pkg/config"
	"tcbmerchantsite/pkg/connection"
	"tcbmerchantsite/pkg/log"
	"tcbmerchantsite/pkg/repository"
	"tcbmerchantsite/pkg/service"
)

var BackOfficeCmd = &cobra.Command{
	Use:   "backoffice-kafka-worker",
	Short: "backoffice kafka worker",
	Long:  "Kafka kafka_worker process backoffice",
	Run: func(cmd *cobra.Command, args []string) {
		if err := InitKafkaConsumer(); err != nil {
			panic(err)
		}
	},
}

func InitKafkaConsumer() error {
	di := do.New()
	defer func() {
		di.Shutdown()
	}()
	config.Inject(di)
	connection.Inject(di)
	repository.Inject(di)
	service.Inject(di)
	do.Provide(di, mtp.NewAdapter)
	do.Provide(di, tpay_ops.NewAdapter)

	// set language msg
	language_message.InitLanguage()

	// get handler backoffice
	handlerBackOffice, err := handler.NewBackOfficeHandler(di)
	if err != nil {
		log.Errorw(nil, "Error when create handler backoffice", "error: ", err)
		return err
	}

	// set handler and start
	kafkaConsumer := do.MustInvoke[*queue.KafkaConsumer](di)
	kafkaConsumer.SetHandler(handlerBackOffice)
	log.Infow(nil, "Start consume backoffice kafka...")
	kafkaConsumer.Start()

	return nil
}
