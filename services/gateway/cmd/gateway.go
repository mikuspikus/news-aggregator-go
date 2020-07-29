package main

import (
	"github.com/caarlos0/env"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
	"github.com/mikuspikus/news-aggregator-go/services/gateway/pkg/gateway"
	news "github.com/mikuspikus/news-aggregator-go/services/news/proto"
	stats "github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cfg := gateway.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("could not parse env vars for config: %v", err)
	}

	tracer, closer, err := tracer.NewTracer("gateway", cfg.JaegerAddress)
	defer closer.Close()
	if err != nil {
		log.Fatal(err)
	}

	commentsConnection, err := grpc.Dial(cfg.CommentsAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}
	defer commentsConnection.Close()
	cc := comments.NewCommentsClient(commentsConnection)

	accountsConnection, err := grpc.Dial(cfg.AccountsAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}
	defer accountsConnection.Close()

	ac := accounts.NewAccountsClient(accountsConnection)
	newsConnection, err := grpc.Dial(cfg.NewsAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}
	defer newsConnection.Close()
	nc := news.NewNewsClient(newsConnection)

	statsConnection, err := grpc.Dial(cfg.StatsAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
	}
	defer statsConnection.Close()
	sc := stats.NewStatsClient(statsConnection)

	service := gateway.New(cc, ac, nc, sc, tracer, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()
	go service.NewsWorker(cfg.KafkaBrokerURLs, "my-client-id", cfg.KafkaNewsTopic)

	service.Start(cfg)
}
