package kafka_queue

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	minBytes         = 10e3
	maxBytes         = 10e6
	kafkaDefaultWait = time.Second * 1
)

type KafkaReader struct {
	*kafka.Reader
}

func NewReader(brokerUrls []string, clientId, topic string) *KafkaReader {
	dialer := &kafka.Dialer{
		ClientID: clientId,
		Timeout:  kafkaDialerDefaultTimeOut,
	}
	cfg := kafka.ReaderConfig{
		Brokers:  brokerUrls,
		GroupID:  clientId,
		Topic:    topic,
		MinBytes: minBytes,
		MaxBytes: maxBytes,
		MaxWait:  kafkaDefaultWait,
		Dialer:   dialer,
	}

	kr := kafka.NewReader(cfg)
	return &KafkaReader{kr}
}

func (kr *KafkaReader) Read(ctx context.Context) (key []byte, value []byte, err error) {
	message, err := kr.ReadMessage(ctx)
	if err != nil {
		return nil, nil, err
	}

	key, value = message.Key, message.Value
	return key, value, nil
}
