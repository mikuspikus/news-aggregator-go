module github.com/mikuspikus/news-aggregator-go

replace github.com/mikuspikus/news-aggregator v0.0.0-20200525151209-9efb99e8e9c0 => ../news-aggregator-go

go 1.14

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/jackc/pgx/v4 v4.7.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
)
