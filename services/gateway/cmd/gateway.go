package main

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
	"log"

	"github.com/caarlos0/env"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
	"github.com/mikuspikus/news-aggregator-go/services/gateway/pkg/gateway"
)

// Config contains env vars
type Config struct {
	CommentsAppID     string `env:"COMMENTS_APP_ID" envDefault:"CommentsAppID"`
	CommentsAppSecret string `env:"COMMENTS_APP_SECRET" envDefault:"CommentsAppSecret"`
	CommentsAddr      string `env:"COMMENTS_ADDR"`

	AccountsAppID     string `env:"ACCOUNTS_APP_ID" envDefault:"AccountsAppID"`
	AccountsAppSecret string `env:"ACCOUNTS_APP_SECRET" envDefault:"AccountsAppSecret"`
	AccountsAddr      string `env:"ACCOUNTS_ADDR"`

	Port          int    `env:"GATEWAY_PORT" envDefault:"3009"`
	JaegerAddress string `env:"JAEGER_ADDRESS"`

	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:8080, "`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envDefault:"GET, POST, PATCH, DELETE, OPTIONS"`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envDefault:"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Authorization"`
	AllowCredentials bool     `env:"ALLOWED_CREDENTIALS" envDefault:"true"`
}

func main() {
	cfg := Config{}
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

	service := gateway.New(cc, ac, tracer)
	if err != nil {
		log.Fatal(err)
	}

	service.Start(cfg.Port, cfg.AllowedOrigins, cfg.AllowedMethods, cfg.AllowedHeaders, cfg.AllowCredentials)
}
