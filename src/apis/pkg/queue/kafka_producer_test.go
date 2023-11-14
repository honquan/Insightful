package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

type testKafkaMsg struct {
	Content string `json:"content"`
	Status  int64  `json:"status"`
}

func KafkaProducerExample(t *testing.T) {
	kafkaProducer, err := NewKafkaProducer(ProducerConfig{
		SeedBrokers:      []string{"vep-kafka-1.int.vinid.dev:9093", "vep-kafka-2.int.vinid.dev:9093", "vep-kafka-3.int.vinid.dev:9093"},
		NumFlushMessages: 1,
		EnableTLS:        true,
		CaCertFile:       "./certs/kafka-int-ca.pem",
		ClientCertFile:   "./certs/test-clients.int.vinid.dev-cert.pem",
		ClientKeyFile:    "./certs/test-clients.int.vinid.dev-key.pem",
	})

	if err != nil {
		t.Fatal(err)
	}

	defer kafkaProducer.Close()

	jsonMsg := &testKafkaMsg{
		Content: fmt.Sprintf("%v", time.Now().Unix()),
		Status:  200,
	}

	msgBytes, err := json.Marshal(jsonMsg)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		kafkaProducer.SendMessage(&KafkaMessage{
			Topic: "test",
			Key:   nil,
			Value: msgBytes,
		})
		<-time.After(time.Second * 2)
		fmt.Println(i)
	}

	<-time.After(time.Second * 5)
}

func TestKafkaProducer_SendAbstractMessage(t *testing.T) {
	if os.Getenv("ENVIRONMENT") != "LOCAL" {
		fmt.Println("Only run this test on local environment!")
		return
	}
	type Struct1 struct {
		ID int64 `json:"id"`
	}

	type Struct2 struct {
		Name string `json:"name"`
	}
	topicMap := make(map[string]string)
	topicMap["Struct1"] = "test-1"
	topicMap["Struct2"] = "test-2"
	kafkaProducer, err := NewKafkaProducer(ProducerConfig{
		SeedBrokers:      []string{"localhost:9092"},
		NumFlushMessages: 1,
		TopicMap:         topicMap,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = kafkaProducer.SendAbstractMessage(&Struct1{ID: 123})
	if err != nil {
		t.Fatal(err)
	}
	err = kafkaProducer.SendAbstractMessage(&Struct2{Name: "name-1"})
	if err != nil {
		t.Fatal(err)
	}
	kafkaProducer.Close()
}
