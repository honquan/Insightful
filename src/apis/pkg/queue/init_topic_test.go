package queue

import (
	"fmt"
	"os"
	"testing"
)

func TestKafkaConfig_CreateTopicIfNeeded(t *testing.T) {
	if os.Getenv("ENVIRONMENT") != "LOCAL" {
		fmt.Println("Only run this test on local environment!")
		return
	}
	topicMap := make(map[string]string)
	topicMap["Struct1"] = "test-1"
	topicMap["Struct2"] = "test-2"

	kafkaClient := &KafkaConfig{
		SeedBrokers:      []string{"localhost:9092"},
		NumFlushMessages: 1,
		TopicMap:         topicMap,
	}

	var numPartition int32
	numPartition = 8

	err := kafkaClient.CreateTopicIfNeeded(numPartition)
	if err != nil {
		t.Fatal(err)
	}
}
