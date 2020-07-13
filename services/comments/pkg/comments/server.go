package comments

import (
	"fmt"
	"log"
	"net"

	// Import the generated protobuf code
	pb "github.com/mikuspikus/news-aggregator-go/services/comments/proto"

	"github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
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

// New returns new Service instance
func New(connString, addr, password string, db int, apps map[string]string) (*Service, error) {
	datastore, err := NewDataStore(connString)
	if err != nil {
		return nil, err
	}

	storage, err := simple_token_storage.New(addr, password, db, apps)
	if err != nil {
		return nil, err
	}

	return &Service{db: datastore, tokenStorage: storage}, nil

}

// NewDataStore returns DataStoreHandler instance
func NewDataStore(connString string) (DataStoreHandler, error) {
	return newDB(connString)
}

// Start starts Comments Service server
func (s *Service) Start(port int, tracer opentracing.Tracer) error {
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	pb.RegisterCommentsServer(server, s)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return server.Serve(listener)
}
