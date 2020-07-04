FROM golang:1.14.4-alpine as builder

RUN apk update && apk upgrade && apk add --no-cache git

RUN mkdir /news-aggregator-go
WORKDIR /news-aggregator-go
COPY . .

ENV GO111MODULE=on
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./services/comments/cmd/...

# Run container
FROM scratch

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /news-aggregator-go/comments .

CMD ["./comments"]