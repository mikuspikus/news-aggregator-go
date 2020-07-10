package gateway

import (
	"context"
	"github.com/mikuspikus/news-aggregator-go/pkg/tracer"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

type gRPCServiceClient interface {
}

type CommentsClient struct {
	client    comments.CommentsClient
	token     string
	appID     string
	appSECRET string
}

// UpdateToken requests new service-to-service token from [comments]
func (cc *CommentsClient) UpdateToken() error {
	token, err := cc.client.GetToken(context.Background(), &comments.GetTokenRequest{
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
			appID:     "",
			appSECRET: "",
		},
		Accounts: &AccountsClient{
			client:    ac,
			token:     "",
			appID:     "",
			appSECRET: "",
		},
		Router: tracer.NewRouter(tr),
	}
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
