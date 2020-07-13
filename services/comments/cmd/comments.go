package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	"github.com/mikuspikus/news-aggregator-go/services/comments/pkg/comments"
)

func main() {
	cfg := comments.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("could not parse env vars for config: %v", err)
	}

	tracer, closer, err := tracer.NewTracer("comments", cfg.JaegerAddress)
	defer closer.Close()

	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbDatabase)
	apps := map[string]string{cfg.AppID: cfg.AppSecret}
	service, err := comments.New(connString, cfg.RedisHost, cfg.RedisPwd, cfg.RedisDb, apps)
	if err != nil {
		log.Fatal(err)
	}

	service.Start(cfg.Port, tracer)
}
