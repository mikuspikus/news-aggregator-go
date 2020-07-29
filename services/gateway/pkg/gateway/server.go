package gateway

import (
	"context"
	"fmt"
	kqueue "github.com/mikuspikus/news-aggregator-go/pkg/kafka-queue"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
	news "github.com/mikuspikus/news-aggregator-go/services/news/proto"
	stats "github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/cors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// Config contains env vars
type Config struct {
	KafkaBrokerURLs    []string `env:"KAFKA_BROKER_URLS" envDefault:"kafka:9092"`
	KafkaClientID      string   `env:"KAFKA_CLIENT_ID" envDefault:"1"`
	KafkaAccountsTopic string   `env:"KAFKA_ACCOUNTS_TOPIC" envDefault:"Accounts"`
	KafkaNewsTopic     string   `env:"KAFKA_NEWS_TOPIC" envDefault:"News"`
	KafkaCommentsTopic string   `env:"KAFKA_COMMENTS_TOPIC" envDefault:"Comments"`

	CommentsAppID     string `env:"COMMENTS_APP_ID" envDefault:"CommentsAppID"`
	CommentsAppSecret string `env:"COMMENTS_APP_SECRET" envDefault:"CommentsAppSecret"`
	CommentsAddr      string `env:"COMMENTS_ADDR"`

	AccountsAppID     string `env:"ACCOUNTS_APP_ID" envDefault:"AccountsAppID"`
	AccountsAppSecret string `env:"ACCOUNTS_APP_SECRET" envDefault:"AccountsAppSecret"`
	AccountsAddr      string `env:"ACCOUNTS_ADDR"`

	NewsAppID     string `env:"NEWS_APP_ID" envDefault:"CommentsAppID"`
	NewsAppSecret string `env:"NEWS_APP_SECRET" envDefault:"CommentsAppSecret"`
	NewsAddr      string `env:"NEWS_ADDR"`

	StatsAppID     string `env:"STATS_APP_ID" envDefault:"StatsAppID"`
	StatsAppSecret string `env:"STATS_APP_SECRET" envDefault:"StatsAppSecret"`
	StatsAddr      string `env:"STATS_ADDR"`

	Port          int    `env:"GATEWAY_PORT" envDefault:"8080"`
	JaegerAddress string `env:"JAEGER_ADDRESS"`

	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:8000"`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envSeparator:"," envDefault:"GET,POST,PATCH,DELETE,OPTIONS"`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envSeparator:"," envDefault:"Origin,X-Requested-With,Content-Type,Accept,Access-Control-Allow-Origin,Authorization"`
	AllowCredentials bool     `env:"ALLOWED_CREDENTIALS" envDefault:"true"`
}

type StatsClient struct {
	client    stats.StatsClient
	token     string
	appID     string
	appSECRET string
}

func (sc *StatsClient) UpdateToken(ctx context.Context) error {
	token, err := sc.client.GetServiceToken(ctx, &stats.GetServiceTokenRequest{
		AppID:     sc.appID,
		AppSECRET: sc.appSECRET,
	})
	if err != nil {
		return err
	}
	sc.token = token.Token
	return nil
}

type NewsClient struct {
	client    news.NewsClient
	Writer    *kqueue.KafkaWriter
	token     string
	appID     string
	appSECRET string
}

func (nc *NewsClient) UpdateToken(ctx context.Context) error {
	token, err := nc.client.GetServiceToken(ctx, &news.GetServiceTokenRequest{
		AppSECRET: nc.appSECRET,
		AppID:     nc.appID,
	})
	if err != nil {
		return err
	}
	nc.token = token.Token
	return nil
}

type CommentsClient struct {
	client    comments.CommentsClient
	Writer    *kqueue.KafkaWriter
	token     string
	appID     string
	appSECRET string
}

// UpdateToken requests new service-to-service token from [comments]
func (cc *CommentsClient) UpdateToken(ctx context.Context) error {
	token, err := cc.client.GetServiceToken(ctx, &comments.GetServiceTokenRequest{
		AppSECRET: cc.appSECRET,
		AppID:     cc.appID,
	})
	if err != nil {
		return err
	}
	cc.token = token.Token
	return nil
}

type AccountsClient struct {
	client    accounts.AccountsClient
	Writer    *kqueue.KafkaWriter
	token     string
	appID     string
	appSECRET string
}

// UpdateToken requests new service-to-service token from [comments]
func (ac *AccountsClient) UpdateToken(ctx context.Context) error {
	token, err := ac.client.GetServiceToken(ctx, &accounts.GetServiceTokenRequest{
		AppSECRET: ac.appSECRET,
		AppID:     ac.appID,
	})
	if err != nil {
		return err
	}
	ac.token = token.Token
	return nil
}

type Server struct {
	Router *tracer.TracedRouter

	Comments *CommentsClient
	Accounts *AccountsClient
	News     *NewsClient
	Stats    *StatsClient
}

func New(cc comments.CommentsClient, ac accounts.AccountsClient, nc news.NewsClient, sc stats.StatsClient, tr opentracing.Tracer, cfg Config) *Server {
	return &Server{
		Comments: &CommentsClient{
			client:    cc,
			token:     "",
			Writer:    kqueue.NewWriter(cfg.KafkaBrokerURLs, cfg.KafkaClientID, cfg.KafkaCommentsTopic),
			appID:     cfg.CommentsAppID,
			appSECRET: cfg.CommentsAppSecret,
		},
		Accounts: &AccountsClient{
			client:    ac,
			token:     "",
			Writer:    kqueue.NewWriter(cfg.KafkaBrokerURLs, cfg.KafkaClientID, cfg.KafkaAccountsTopic),
			appID:     cfg.AccountsAppID,
			appSECRET: cfg.AccountsAppSecret,
		},
		News: &NewsClient{
			client:    nc,
			token:     "",
			Writer:    kqueue.NewWriter(cfg.KafkaBrokerURLs, cfg.KafkaClientID, cfg.KafkaNewsTopic),
			appID:     cfg.NewsAppID,
			appSECRET: cfg.NewsAppSecret,
		},
		Stats: &StatsClient{
			client:    sc,
			token:     "",
			appID:     cfg.StatsAppID,
			appSECRET: cfg.StatsAppSecret,
		},
		Router: tracer.NewRouter(tr),
	}
}

func (s *Server) Start(cfg Config) {
	cors := cors.New(cors.Options{
		//Debug: true,
		//AllowedOrigins:   []string{"http://localhost:8000"},
		//AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Access-Control-Allow-Origin", "Authorization"},
		//AllowCredentials: true,
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: cfg.AllowCredentials,
	})

	s.Router.Mux.Use(setContentType)
	s.routes()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors.Handler(s.Router),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("HTTP server is shutting down")
	os.Exit(0)
}

func (s *Server) Close() {
	s.Accounts.Writer.Close()
	s.News.Writer.Close()
	s.Comments.Writer.Close()
}

func getAuthorizationToken(req *http.Request) string {
	header := req.Header.Get("Authorization")
	splittedHeader := strings.Split(header, " ")
	if len(splittedHeader) == 2 {
		return splittedHeader[1]
	}
	return ""
}

func setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// handleRPCErrors converts RPC error into http ones
func handleRPCErrors(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)

	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		http.Error(w, st.Message(), http.StatusNotFound)
		return

	case codes.InvalidArgument:
		http.Error(w, st.Message(), http.StatusBadRequest)
		return

	case codes.Unauthenticated:
		w.WriteHeader(http.StatusForbidden)
		return

	case codes.Unavailable:
		w.WriteHeader(http.StatusServiceUnavailable)
		return

	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
