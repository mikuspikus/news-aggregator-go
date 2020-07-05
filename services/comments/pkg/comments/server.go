package comments

import (
	"fmt"
	"log"
	"net"

	// Import the generated protobuf code
	pb "github.com/mikuspikus/news-aggregator-go/services/comments/proto"

	"github.com/mikuspikus/news-aggregator-go/pkg/token-storage"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
)

// NewService returns new Service instance
func NewService(connString, addr, password string, db int, apps map[string]string) (*Service, error) {
	datastore, err := NewDataStore(connString)
	if err != nil {
		return nil, err
	}

	storage, err := token_storage.New(addr, password, db, apps)
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
