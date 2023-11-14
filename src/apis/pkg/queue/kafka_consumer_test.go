package queue

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

type HandlerTest struct {
}

func (h *HandlerTest) HandlerFunc(msg *ConsumerMessage) {
	log.Printf("Receive msg, key: %v, value: %v, topic: %v \n", string(msg.Key), string(msg.Value), string(msg.Topic))
}

func (h *HandlerTest) Close() {

}

func ConsumerGroupExample(t *testing.T) {
	config := ConsumerConfig{
		Topics:          []string{"test"},
		SeedBrokers:     []string{"vep-kafka-1.int.vinid.dev:9093", "vep-kafka-2.int.vinid.dev:9093", "vep-kafka-3.int.vinid.dev:9093"},
		ConsumerGroupID: "cg-id-1",
		EnableTLS:       true,
		CaCertFile:      "./certs/kafka-int-ca.pem",
		ClientCertFile:  "./certs/test-clients.int.vinid.dev-cert.pem",
		ClientKeyFile:   "./certs/test-clients.int.vinid.dev-key.pem",
	}
	consumerGroup, err := NewKafkaConsumer(config)
	if err != nil {
		t.Fatal(err)
	}

	handlerTest := &HandlerTest{}

	consumerGroup.SetHandler(handlerTest)

	consumerGroup.Start()

	<-time.After(10 * time.Second)

	consumerGroup.Close()

}

func TestKafkaConsumer_BindingMessage_ShouldSuccess(t *testing.T) {
	if os.Getenv("ENVIRONMENT") != "LOCAL" {
		fmt.Println("Only run this test on local environment!")
		return
	}
	config := ConsumerConfig{
		Topics:          []string{"test-1"},
		SeedBrokers:     []string{"localhost:9092"},
		ConsumerGroupID: "cg-id-1",
	}
	consumerGroup, err := NewKafkaConsumer(config)
	if err != nil {
		t.Fatal(err)
	}

	handlerTest := &HandlerTest{}

	consumerGroup.SetHandler(handlerTest)

	consumerGroup.Start()

	<-time.After(10 * time.Second)

	consumerGroup.Close()
}
