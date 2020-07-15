package kafka_queue

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	kafkaDialerDefaultTimeOut = time.Second * 10
	kafkaWriteDefaultTimeout  = time.Second * 10
	kafkaReadDefaultTimeout   = time.Second * 10
)

type KafkaWriter struct {
	*kafka.Writer
}

func NewWriter(kafkaBrokerUrls []string, clientId string, topic string) *KafkaWriter {
	dialer := &kafka.Dialer{
		Timeout:  kafkaDialerDefaultTimeOut,
		ClientID: clientId,
	}

	cfg := kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		WriteTimeout: kafkaWriteDefaultTimeout,
		ReadTimeout:  kafkaReadDefaultTimeout,
	}

	kw := kafka.NewWriter(cfg)
	return &KafkaWriter{kw}
}

func (kw *KafkaWriter) Push(ctx context.Context, key, value []byte) error {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return kw.WriteMessages(ctx, message)
}
