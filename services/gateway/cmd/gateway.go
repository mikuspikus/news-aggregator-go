package main

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	news "github.com/mikuspikus/news-aggregator-go/services/news/proto"
	"google.golang.org/grpc"
	"log"

	"github.com/caarlos0/env"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
	"github.com/mikuspikus/news-aggregator-go/services/gateway/pkg/gateway"
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

	service := gateway.New(cc, ac, nc, tracer, cfg)
	if err != nil {
		log.Fatal(err)
	}

	service.Start(cfg)
}
