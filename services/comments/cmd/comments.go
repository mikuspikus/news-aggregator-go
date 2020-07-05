package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	"github.com/mikuspikus/news-aggregator-go/services/comments/pkg/comments"
)

// Config contains env vars
type Config struct {
	AppID         string `env:"COMMENTS_APP_ID" envDefault:"CommentsAppID"`
	AppSecret     string `env:"COMMENTS_APP_SECRET" envDefault:"CommentsAppSecret"`
	Port          int    `env:"COMMENTS_PORT" envDefault:"3009"`
	JaegerAddress string `env:"JAEGER_ADDRESS"`
	DbHost        string `env:"DB_HOST" envDefault:"localhost"`
	DbPort        int    `env:"DB_PORT" envDefault:"5432"`
	DbUser        string `env:"POSTGRES_USER" envDefault:"user"`
	DbPassword    string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	DbDatabase    string `env:"POSTGRES_DB"`
	RedisHost     string `env:"REDIS_HOST"`
	RedisPwd      string `env:"REDIS_PWD"`
	RedisDb       int    `env:"REDIS_DB"`
}

func main() {
	cfg := Config{}
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
	service, err := comments.NewService(connString, cfg.RedisHost, cfg.RedisPwd, cfg.RedisDb, apps)
	if err != nil {
		log.Fatal(err)
	}

	service.Start(cfg.Port, tracer)
}
