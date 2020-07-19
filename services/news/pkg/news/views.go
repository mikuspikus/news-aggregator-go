package news

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"

	stst "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/news/proto"
)

var (
	statusServiceTokenNotFound = status.Error(codes.Unauthenticated, "service token not found")
	statusInvalidServiceToken  = status.Error(codes.Unauthenticated, "invalid service token")
	statusInvalidUUID          = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidURI           = status.Error(codes.InvalidArgument, "invalid URI")
	statusNotFound             = status.Error(codes.NotFound, "news not found")
	statusInvalidToken         = status.Error(codes.Unauthenticated, "invalid token")
	statusAppIDNotFound        = status.Error(codes.NotFound, "app ID not found")
	statusInvalidSecret        = status.Error(codes.InvalidArgument, "invalid secret")
)

func internalServerError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func (news *News) SingleNews() (*pb.SingleNews, error) {
	created, err := ptypes.TimestampProto(news.Created)
	if err != nil {
		return nil, internalServerError(err)
	}
	edited, err := ptypes.TimestampProto(news.Edited)
	if err != nil {
		return nil, internalServerError(err)
	}

	singleNews := new(pb.SingleNews)
	singleNews.Uid = news.UID.String()
	singleNews.UserUUID = news.User.String()
	singleNews.Title = news.Title
	singleNews.Uri = news.URI.String()
	singleNews.Created = created
	singleNews.Edited = edited

	return singleNews, nil
}

type Service struct {
	db           DataStoreHandler
	tokenStorage *stst.APITokenStorage
}

func (s *Service) validateServiceToken(token string) error {
	valid, err := s.tokenStorage.CheckToken(token)
	if err == stst.ErrTokenNotFound {
		return statusServiceTokenNotFound
	}
	if err != nil {
		return err
	}
	if !valid {
		return statusInvalidServiceToken
	}
	return nil
}

func (s *Service) GetServiceToken(ctx context.Context, req *pb.GetServiceTokenRequest) (*pb.GetServiceTokenResponse, error) {
	appID, appSECRET := req.AppID, req.AppSECRET
	token, err := s.tokenStorage.AddToken(appID, appSECRET)
	switch err {
	case nil:
		response := new(pb.GetServiceTokenResponse)
		response.Token = token
		return response, nil
	case stst.ErrIDNotFound:
		return nil, statusAppIDNotFound
	case stst.ErrWrongSecret:
		return nil, statusInvalidSecret
	default:
		return nil, internalServerError(err)
	}
}

func (s *Service) ListNews(ctx context.Context, req *pb.ListNewsRequest) (*pb.ListNewsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	news, pageCount, err := s.db.List(req.PageNumber, pageSize)
	if err != nil {
		return nil, internalServerError(err)
	}
	response := new(pb.ListNewsResponse)
	for _, snews := range news {
		singleNews, err := snews.SingleNews()
		if err != nil {
			return nil, internalServerError(err)
		}
		response.News = append(response.News, singleNews)
	}

	response.PageNumber = req.PageNumber
	response.PageSize = pageSize
	response.PageCount = pageCount

	return response, nil
}

func (s *Service) GetNews(ctx context.Context, req *pb.GetNewsRequest) (*pb.NewsResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	news, err := s.db.Get(uid)
	if err != nil {
		return nil, internalServerError(err)
	}
	singleNews, err := news.SingleNews()
	if err != nil {
		return nil, internalServerError(err)
	}
	response := new(pb.NewsResponse)
	response.News = singleNews
	return response, nil

}

func (s *Service) AddNews(ctx context.Context, req *pb.AddNewsRequest) (*pb.NewsResponse, error) {
	statusError := s.validateServiceToken(req.Token)
	if statusError != nil {
		return nil, statusError
	}

	user, err := uuid.Parse(req.UserUUID)
	if err != nil {
		return nil, statusInvalidUUID
	}
	uri, err := url.ParseRequestURI(req.Uri)
	if err != nil {
		return nil, statusInvalidURI
	}
	news, err := s.db.Create(user, req.Title, *uri)
	if err != nil {
		return nil, internalServerError(err)
	}
	singleNews, err := news.SingleNews()
	if err != nil {
		return nil, internalServerError(err)
	}
	response := new(pb.NewsResponse)
	response.News = singleNews
	return response, nil
}

func (s *Service) EditNews(ctx context.Context, req *pb.EditNewsRequest) (*pb.NewsResponse, error) {
	statusError := s.validateServiceToken(req.Token)
	if statusError != nil {
		return nil, statusError
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}
	uri, err := url.ParseRequestURI(req.Uri)
	if err != nil {
		return nil, statusInvalidURI
	}
	news, err := s.db.Update(uid, req.Title, *uri)
	if err != nil {
		return nil, internalServerError(err)
	}
	singleNews, err := news.SingleNews()
	if err != nil {
		return nil, internalServerError(err)
	}
	response := new(pb.NewsResponse)
	response.News = singleNews
	return response, nil
}

func (s *Service) DeleteNews(ctx context.Context, req *pb.DeleteNewsRequest) (*pb.DeleteNewsResponse, error) {
	statusError := s.validateServiceToken(req.Token)
	if statusError != nil {
		return nil, statusError
	}

	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.Delete(uid)
	switch err {
	case nil:
		return new(pb.DeleteNewsResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalServerError(err)
	}
}
