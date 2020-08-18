package accounts

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	ststorage "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	rtstorage "github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/refresh-token-storage"
	utstorage "github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/user-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Config contains env vars
type Config struct {
	AppID         string `env:"ACCOUNTS_APP_ID" envDefault:"AccountsAppID"`
	AppSecret     string `env:"ACCOUNTS_APP_SECRET" envDefault:"AccountsAppSecret"`
	Port          int    `env:"ACCOUNTS_PORT" envDefault:"3009"`
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

type Service struct {
	Store         DataStoreHandler
	ServiceTokens ststorage.APITokenStorage
	RefreshTokens rtstorage.RefreshTokenStorage
	UserTokens    utstorage.UserTokenStorage
}

func New(connString, addr, password string, db int, apps map[string]string) (*Service, error) {
	dbpool, err := NewDB(connString)
	if err != nil {
		return nil, err
	}

	ststorage, err := ststorage.New(addr, password, db, apps)
	if err != nil {
		return nil, err
	}

	rtstorage, err := rtstorage.New(addr, password, db+1)
	if err != nil {
		return nil, err
	}

	utstorage, err := utstorage.New(addr, password, db+2)
	if err != nil {
		return nil, err
	}

	return &Service{Store: dbpool, ServiceTokens: ststorage, RefreshTokens: rtstorage, UserTokens: utstorage}, nil
}

// Start starts Account Service server
func (s *Service) Start(port int, tracer opentracing.Tracer) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterAccountsServer(server, s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return server.Serve(listener)
}

func (s *Service) Close() {
	s.Store.Close()
}
