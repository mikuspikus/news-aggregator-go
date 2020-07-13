package news

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	ststorage "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/news/proto"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
	"net"
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

func New(connString, addr, password string, db int, apps map[string]string) (*Service, error) {
	datastore, err := NewDataStore(connString)
	if err != nil {
		return nil, err
	}
	storage, err := ststorage.New(addr, password, db, apps)
	if err != nil {
		return nil, err
	}

	return &Service{db: datastore, tokenStorage: storage}, nil
}

func (s *Service) Start(port int, tracer opentracing.Tracer) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterNewsServer(server, s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return server.Serve(listener)
}

func NewDataStore(connString string) (DataStoreHandler, error) {
	return newDB(connString)
}
