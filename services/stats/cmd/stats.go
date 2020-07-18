package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	"github.com/mikuspikus/news-aggregator-go/services/stats/pkg/stats"
	"log"
)

func main() {
	cfg := stats.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("could not parse env vars for config: %v", err)
	}
	tracer, closer, err := tracer.NewTracer("stats", cfg.JaegerAddress)
	defer closer.Close()
	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbDatabase)
	apps := map[string]string{cfg.AppID: cfg.AppSecret}
	service, err := stats.New(connString, cfg.RedisHost, cfg.RedisPwd, cfg.RedisDb, apps)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()
	service.Start(cfg.Port, tracer)
}
