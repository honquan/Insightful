package queue

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

type KafkaConfig struct {
	SeedBrokers                               []string
	NumFlushMessages                          int
	TopicMap                                  map[string]string
	EnableTLS                                 bool
	InsecureSkipVerify                        bool
	ClientCertFile, ClientKeyFile, CaCertFile string
}

func (this *KafkaConfig) CreateTopicIfNeeded(numPartition int32) error {
	return this.CreateTopicIfNeededWithReplica(numPartition, 1)
}

func (this *KafkaConfig) CreateTopicIfNeededWithReplica(numPartition int32, numReplica int16) error {
	var errs error
	for _, brokerConfig := range this.SeedBrokers {
		if brokerConfig == "" {
			errs = errors.New("kafka host port is required")
			break
		}
		kafkaConfig := strings.Split(brokerConfig, ":")
		kafkaHost := kafkaConfig[0]
		kafkaPort := kafkaConfig[1]
		if kafkaHost == "" || kafkaPort == "" {
			errs = errors.New("kafka host port invalid")
			break
		}

		broker := sarama.NewBroker(fmt.Sprintf("%v:%v", kafkaHost, kafkaPort))
		config := sarama.NewConfig()
		config.Version = sarama.V1_1_0_0

		if this.EnableTLS {
			tlsConfig, err := newTLSConfig(this.ClientCertFile,
				this.ClientKeyFile,
				this.CaCertFile,
			)
			if err != nil {
				return err
			}
			config.Net.TLS.Enable = true
			config.Net.TLS.Config = tlsConfig
			tlsConfig.InsecureSkipVerify = this.InsecureSkipVerify
		}

		err := broker.Open(config)

		if err != nil {
			errs = err
			break
		}

		// Setup the Topic details in CreateTopicRequest struct
		topicDetail := &sarama.TopicDetail{}
		topicDetail.NumPartitions = numPartition
		topicDetail.ReplicationFactor = numReplica
		topicDetail.ConfigEntries = make(map[string]*string)

		topicDetails := make(map[string]*sarama.TopicDetail)
		for _, topicItem := range this.TopicMap {
			topicDetails[topicItem] = topicDetail
		}

		request := sarama.CreateTopicsRequest{
			Timeout:      time.Second * 15,
			TopicDetails: topicDetails,
		}

		// Send request to Broker
		_, err = broker.CreateTopics(&request)

		// handle errors if any
		if err != nil {
			errs = err
			break
		}

		// close connection to broker
		broker.Close()
	}

	if errs != nil {
		return errs
	}

	return nil
}
