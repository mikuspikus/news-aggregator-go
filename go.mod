module github.com/mikuspikus/news-aggregator-go

replace github.com/mikuspikus/news-aggregator v0.0.0-20200525151209-9efb99e8e9c0 => ../news-aggregator-go

go 1.14

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-redis/redis/v8 v8.0.0-beta.6
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/jackc/pgx/v4 v4.7.1
	github.com/opentracing-contrib/go-stdlib v1.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/rs/cors v1.7.0
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
)
