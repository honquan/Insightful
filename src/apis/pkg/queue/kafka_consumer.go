package queue

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type InitOffsetType int64

const (
	InitOffsetNewest InitOffsetType = -1
	InitOffsetOldest InitOffsetType = -2
)

type ConsumerConfig struct {
	SeedBrokers                               []string
	ConsumerGroupID                           string
	Topics                                    []string
	InitialOffset                             InitOffsetType
	EnableTLS                                 bool
	InsecureSkipVerify                        bool
	ClientCertFile, ClientKeyFile, CaCertFile string
}

type ConsumerMessage struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
	Timestamp time.Time
}

type KafkaConsumerHandlerFunc func(message *ConsumerMessage)

type consumerGroupHandler struct {
	handleFunc KafkaConsumerHandlerFunc
}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (cg consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		cg.handleFunc(&ConsumerMessage{
			Key:       msg.Key,
			Value:     msg.Value,
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Timestamp: msg.Timestamp,
		})
		sess.MarkMessage(msg, "")
	}
	return nil
}

type ConsumerHandler interface {
	HandlerFunc(*ConsumerMessage)
	Close()
}

type KafkaConsumer struct {
	client        sarama.Client
	consumerGroup sarama.ConsumerGroup
	topics        []string
	handler       ConsumerHandler

	running bool
}

func NewKafkaConsumer(cf ConsumerConfig) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_0_0
	config.Consumer.Return.Errors = true

	if cf.EnableTLS {
		tlsConfig, err := newTLSConfig(cf.ClientCertFile,
			cf.ClientKeyFile,
			cf.CaCertFile,
		)
		if err != nil {
			return nil, err
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
		tlsConfig.InsecureSkipVerify = cf.InsecureSkipVerify
	}

	switch cf.InitialOffset {
	case InitOffsetNewest:
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	case InitOffsetOldest:
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	kkClient, err := sarama.NewClient(cf.SeedBrokers, config)
	if err != nil {
		return nil, err
	}

	kkConsumerGroup, err := sarama.NewConsumerGroupFromClient(cf.ConsumerGroupID, kkClient)
	if err != nil {
		return nil, err
	}

	consumer := &KafkaConsumer{
		client:        kkClient,
		consumerGroup: kkConsumerGroup,
		topics:        cf.Topics,
	}

	return consumer, nil
}

func (c *KafkaConsumer) SetHandler(fn ConsumerHandler) {
	c.handler = fn
}

func (c *KafkaConsumer) Start() {
	// Track errors
	go func() {
		for err := range c.consumerGroup.Errors() {
			fmt.Println("ConsumerGroup error, detail: ", err)
		}
	}()

	// Iterate over consumer sessions.
	c.running = true

	ctx := context.Background()
	for c.running {
		handler := consumerGroupHandler{handleFunc: c.handler.HandlerFunc}

		err := c.consumerGroup.Consume(ctx, c.topics, handler)
		if err != nil {
			if !c.running {
				fmt.Println("Consumer closed")
			} else {
				fmt.Println("Error: ", err)
			}
		}
	}
}

func (c *KafkaConsumer) Close() {
	c.running = false
	c.consumerGroup.Close()
	c.handler.Close()
}
