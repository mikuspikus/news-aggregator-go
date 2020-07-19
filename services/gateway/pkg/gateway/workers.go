package gateway

import (
	"context"
	"encoding/json"
	"github.com/mikuspikus/news-aggregator-go/pkg/kafka-queue"
	stats "github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Message struct {
	UserUID string `json:"user_uid"`
	Action  string `json:"action"`
	Input   []byte `json:"input"`
	Output  []byte `json:"output"`
}

func convertMessageIntoRequest(msg *Message) (*stats.AddStatsRequest, error) {
	statsRequest := new(stats.AddStatsRequest)
	statsRequest.UserUID = msg.UserUID
	statsRequest.Action = msg.Action
	statsRequest.Input = msg.Input
	statsRequest.Output = msg.Output
	return statsRequest, nil
}

func (sc *StatsClient) PushAccountsStat(ctx context.Context, msg *Message) error {
	request, err := convertMessageIntoRequest(msg)
	if err != nil {
		return err
	}

	_, err = sc.client.AddAccountsStats(ctx, request)

	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = sc.UpdateToken(ctx)
			if err != nil {
				return err
			}
			_, err = sc.client.AddAccountsStats(ctx, request)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (sc *StatsClient) PushNewsStats(ctx context.Context, msg *Message) error {
	request, err := convertMessageIntoRequest(msg)
	if err != nil {
		return err
	}
	request.Token = sc.token

	_, err = sc.client.AddNewsStats(ctx, request)

	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = sc.UpdateToken(ctx)
			if err != nil {
				return err
			}
			request.Token = sc.token
			_, err = sc.client.AddNewsStats(ctx, request)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (sc *StatsClient) PushCommentsStats(ctx context.Context, msg *Message) error {
	request, err := convertMessageIntoRequest(msg)
	if err != nil {
		return err
	}

	_, err = sc.client.AddCommentsStats(ctx, request)

	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = sc.UpdateToken(ctx)
			if err != nil {
				return err
			}
			_, err = sc.client.AddCommentsStats(ctx, request)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (s *Server) AccountsWorker(brokerUrls []string, clientId, topic string) {
	reader := kafka_queue.NewReader(brokerUrls, clientId, topic)

	for {
		ctx := context.Background()
		value, _, err := reader.Read(ctx)
		if err != nil {
			log.Println(err)
			continue
		}
		msg := new(Message)
		err = json.Unmarshal(value, msg)
		if err != nil {
			log.Println(err)
			continue
		}
		err = s.Stats.PushAccountsStat(ctx, msg)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (s *Server) NewsWorker(brokerUrls []string, clientId, topic string) {
	log.Println("Worker started")
	reader := kafka_queue.NewReader(brokerUrls, clientId, topic)
	log.Println("Reader is up")

	for {
		ctx := context.Background()
		log.Println("Waiting for task")
		_, value, err := reader.Read(ctx)
		log.Println("Received new message from Kafka")
		if err != nil {
			log.Printf("Reader error: %v", err)
			continue
		}
		var msg Message
		err = json.Unmarshal(value, &msg)
		if err != nil {
			log.Printf("Formatting err: %v", err)
			continue
		}
		err = s.Stats.PushNewsStats(ctx, &msg)
		if err != nil {
			log.Printf("Pushing err: %v", err)
			continue
		}
	}
}

func (s *Server) CommentsWorker(brokerUrls []string, clientId, topic string) {
	reader := kafka_queue.NewReader(brokerUrls, clientId, topic)

	for {
		ctx := context.Background()
		value, _, err := reader.Read(ctx)
		if err != nil {
			log.Println(err)
			continue
		}
		msg := new(Message)
		err = json.Unmarshal(value, msg)
		if err != nil {
			log.Println(err)
			continue
		}
		err = s.Stats.PushCommentsStats(ctx, msg)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
