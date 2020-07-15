FROM golang:1.14.4-alpine as builder

RUN mkdir /news-aggregator-go
WORKDIR /news-aggregator-go
COPY . .
RUN go build -mod=vendor ./services/stats/cmd/stats.go

# Run container
FROM golang:1.14.4-alpine
RUN apk --no-cache add ca-certificates
ENV APP_HOME=/home/app
RUN mkdir $APP_HOME
WORKDIR $APP_HOME

# create the app user
RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /news-aggregator-go/stats .
COPY ./services/stats/entrypoint.sh .

RUN chmod +x $APP_HOME/entrypoint.sh
# chown all the files to the app user
RUN chown -R app:app $APP_HOME

USER app

CMD ["/home/app/entrypoint.sh"]