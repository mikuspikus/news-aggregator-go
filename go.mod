module github.com/mikuspikus/news-aggregator-go

replace github.com/mikuspikus/news-aggregator v0.0.0-20200525151209-9efb99e8e9c0 => ../news-aggregator-go

go 1.14

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/go-redis/redis/v8 v8.0.0-beta.7
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/jackc/pgconn v1.6.4
	github.com/jackc/pgproto3/v2 v2.0.4 // indirect
	github.com/jackc/pgx/v4 v4.8.1
	github.com/klauspost/compress v1.10.11 // indirect
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/rs/cors v1.7.0
	github.com/segmentio/kafka-go v0.4.2
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.opentelemetry.io/otel v0.10.0 // indirect
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc // indirect
	golang.org/x/sys v0.0.0-20200814200057-3d37ad5750ed // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
)
