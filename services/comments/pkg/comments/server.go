package comments

import (
	"fmt"
	"log"
	"net"

	// Import the generated protobuf code
	pb "github.com/mikuspikus/news-aggregator-go/services/comments/proto"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc"
)

// NewService returns new Service instance
func NewService(connString string) (*Service, error) {
	datastore, err := NewDataStore(connString)
	if err != nil {
		return nil, err
	}

	return &Service{db: datastore}, nil

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
