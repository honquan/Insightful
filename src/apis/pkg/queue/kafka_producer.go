package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"reflect"
	"sync"
	"time"
)

type ProducerType string

const (
	KafkaSyncProducerType  ProducerType = "sync"
	KafkaAsyncProducerType ProducerType = "async"
)

var (
	NoTopicDefined = errors.New("no topic defined")
)

type KafkaMessage struct {
	Topic  string
	Key    []byte
	Value  []byte
	Header []sarama.RecordHeader
}

type ProducerConfig struct {
	SeedBrokers                               []string
	NumFlushMessages                          int
	TopicMap                                  map[string]string
	EnableTLS                                 bool
	InsecureSkipVerify                        bool
	ClientCertFile, ClientKeyFile, CaCertFile string
}

// KafkaProducer --
type KafkaProducer struct {
	producer sarama.AsyncProducer
	topicMap map[string]string
}

func NewKafkaProducer(cf ProducerConfig) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_0_0
	config.Producer.Flush.Messages = cf.NumFlushMessages
	config.Producer.Flush.Frequency = 1 * time.Second
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	config.Producer.Partitioner = sarama.NewHashPartitioner

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

	asyncProducer, err := sarama.NewAsyncProducer(cf.SeedBrokers, config)
	if err != nil {
		return nil, err
	}

	kafkaProducer := &KafkaProducer{
		producer: asyncProducer,
		topicMap: cf.TopicMap,
	}

	return kafkaProducer, nil
}

func (p *KafkaProducer) SendMessage(m *KafkaMessage) {
	msg := &sarama.ProducerMessage{
		Topic:   m.Topic,
		Value:   sarama.ByteEncoder(m.Value),
		Headers: m.Header,
	}

	if m.Key != nil {
		msg.Key = sarama.ByteEncoder(m.Key)
	}

	p.producer.Input() <- msg
}

func (p *KafkaProducer) SendAbstractMessage(msg interface{}) error {
	msgStructName := p.getTypeOfMessage(msg)
	topic := p.topicMap[msgStructName]
	if topic == "" {
		return NoTopicDefined
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.ByteEncoder(msgBytes),
	}

	p.producer.Input() <- kafkaMsg
	return nil
}

func (p *KafkaProducer) getTypeOfMessage(msg interface{}) string {
	if t := reflect.TypeOf(msg); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func (p *KafkaProducer) Close() {
	var wg sync.WaitGroup
	p.producer.AsyncClose()

	wg.Add(2)
	go func() {
		for range p.producer.Successes() {
			fmt.Println("Unexpected message on Successes()")
		}
		wg.Done()
	}()
	go func() {
		for msg := range p.producer.Errors() {
			fmt.Println(msg.Err)
		}
		wg.Done()
	}()
	wg.Wait()
}
