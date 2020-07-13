package gateway

import (
	"context"
	"fmt"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
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
	CommentsAppID     string `env:"COMMENTS_APP_ID" envDefault:"CommentsAppID"`
	CommentsAppSecret string `env:"COMMENTS_APP_SECRET" envDefault:"CommentsAppSecret"`
	CommentsAddr      string `env:"COMMENTS_ADDR"`

	AccountsAppID     string `env:"ACCOUNTS_APP_ID" envDefault:"AccountsAppID"`
	AccountsAppSecret string `env:"ACCOUNTS_APP_SECRET" envDefault:"AccountsAppSecret"`
	AccountsAddr      string `env:"ACCOUNTS_ADDR"`

	Port          int    `env:"GATEWAY_PORT" envDefault:"3009"`
	JaegerAddress string `env:"JAEGER_ADDRESS"`

	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:8080, "`
	AllowedMethods   []string `env:"ALLOWED_METHODS" envDefault:"GET, POST, PATCH, DELETE, OPTIONS"`
	AllowedHeaders   []string `env:"ALLOWED_HEADERS" envDefault:"Origin, X-Requested-With, Content-Type, Accept, Access-Control-Allow-Origin, Authorization"`
	AllowCredentials bool     `env:"ALLOWED_CREDENTIALS" envDefault:"true"`
}

type CommentsClient struct {
	client    comments.CommentsClient
	token     string
	appID     string
	appSECRET string
}

// UpdateToken requests new service-to-service token from [comments]
func (cc *CommentsClient) UpdateToken() error {
	token, err := cc.client.GetServiceToken(context.Background(), &comments.GetServiceTokenRequest{
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
	token     string
	appID     string
	appSECRET string
}

// UpdateToken requests new service-to-service token from [comments]
func (ac *AccountsClient) UpdateToken() error {
	token, err := ac.client.GetServiceToken(context.Background(), &accounts.GetServiceTokenRequest{
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
	Router   *tracer.TracedRouter
	Comments *CommentsClient
	Accounts *AccountsClient
}

func New(cc comments.CommentsClient, ac accounts.AccountsClient, tr opentracing.Tracer) *Server {
	return &Server{
		Comments: &CommentsClient{
			client:    cc,
			token:     "",
			appID:     "CommentsAppID",
			appSECRET: "CommentsAppSecret",
		},
		Accounts: &AccountsClient{
			client:    ac,
			token:     "",
			appID:     "AccountsAppID",
			appSECRET: "AccountsAppSecret",
		},
		Router: tracer.NewRouter(tr),
	}
}

func (s *Server) Start(port int, AllowedOrigins, AllowedMethods, AllowedHeaders []string, AllowCredentials bool) {
	cors := cors.New(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedMethods:   AllowedMethods,
		AllowedHeaders:   AllowedHeaders,
		AllowCredentials: AllowCredentials,
	})

	s.Router.Mux.Use(setContentType)
	s.routes()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
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
