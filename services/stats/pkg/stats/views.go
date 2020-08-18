package stats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	stst "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errUnknownType             = errors.New("unknown type")
	statusInvalidUUID          = status.Error(codes.InvalidArgument, "invalid UUID")
	statusAppIDNotFound        = status.Error(codes.NotFound, "app ID not found")
	statusInvalidSecret        = status.Error(codes.InvalidArgument, "invalid secret")
	statusInvalidServiceToken  = status.Error(codes.Unauthenticated, "invalid service token")
	statusServiceTokenNotFound = status.Error(codes.Unauthenticated, "service token not found")
)

type target = int32

const (
	ACCOUNTS target = 0
	NEWS     target = 1
	COMMENTS target = 2
)

func internalServerError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func unmarshal(input []byte, target interface{}) error {
	if len(input) == 0 {
		return nil
	}
	return json.Unmarshal(input, target)
}

func (stats *Stats) SingleStat() (*pb.SingleStat, error) {
	var err error
	singleStat := new(pb.SingleStat)
	singleStat.Id = stats.ID
	singleStat.UserUID = stats.User.String()
	singleStat.Action = stats.Action
	singleStat.Timestamp, err = ptypes.TimestampProto(stats.Timestamp)
	if err != nil {
		return nil, err
	}
	singleStat.Input, err = json.Marshal(stats.Input)
	if err != nil {
		return nil, err
	}
	singleStat.Output, err = json.Marshal(stats.Output)
	if err != nil {
		return nil, err
	}
	return singleStat, nil
}

type Service struct {
	db           DataStoreHandler
	tokenStorage stst.APITokenStorage
}

func (s *Service) validateServiceToken(token string) error {
	valid, err := s.tokenStorage.CheckToken(token)
	if err == stst.ErrTokenNotFound {
		return statusServiceTokenNotFound
	}
	if err != nil {
		return internalServerError(err)
	}
	if !valid {
		return statusInvalidServiceToken
	}
	return nil
}

func (s *Service) genericListStats(ctx context.Context, req *pb.ListStatsRequest, target target) (*pb.ListStatsResponse, error) {
	statusError := s.validateServiceToken(req.Token)
	if statusError != nil {
		return nil, statusError
	}

	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	var stats []*Stats
	var pageCount int32
	var err error
	switch target {
	case ACCOUNTS:
		stats, pageCount, err = s.db.ListAccounts(req.PageNumber, pageSize)
	case NEWS:
		stats, pageCount, err = s.db.ListNews(req.PageNumber, pageSize)
	case COMMENTS:
		stats, pageCount, err = s.db.ListComments(req.PageNumber, pageSize)
	default:
		err = errUnknownType
	}

	if err != nil {
		return nil, internalServerError(err)
	}
	response := new(pb.ListStatsResponse)
	for _, stat := range stats {
		singleStat, err := stat.SingleStat()
		if err != nil {
			return nil, internalServerError(err)
		}
		response.Stats = append(response.Stats, singleStat)
	}

	response.PageNumber = req.PageNumber
	response.PageSize = pageSize
	response.PageCount = pageCount

	return response, nil
}

func (s *Service) genericAddStats(ctx context.Context, req *pb.AddStatsRequest, target target) (*pb.StatsResponse, error) {
	statusError := s.validateServiceToken(req.Token)
	if statusError != nil {
		return nil, statusError
	}

	var err error
	var userUID uuid.UUID
	if req.UserUID == "" {
		userUID = uuid.Nil
	} else {
		userUID, err = uuid.Parse(req.UserUID)
		if err != nil {
			return nil, statusInvalidUUID
		}
	}

	var input, output map[string]interface{}
	err = unmarshal(req.Input, &input)
	if err != nil {
		return nil, internalServerError(err)
	}
	err = unmarshal(req.Output, &output)
	if err != nil {
		return nil, internalServerError(err)
	}

	stats := new(Stats)
	switch target {
	case ACCOUNTS:
		stats, err = s.db.AddAccounts(userUID, req.Action, input, output)

	case NEWS:
		stats, err = s.db.AddNews(userUID, req.Action, input, output)

	case COMMENTS:
		stats, err = s.db.AddComments(userUID, req.Action, input, output)

	default:
		err = errUnknownType
	}
	if err != nil {
		return nil, internalServerError(err)
	}

	singleStat, err := stats.SingleStat()
	if err != nil {
		return nil, internalServerError(err)
	}
	response := new(pb.StatsResponse)
	response.Stats = singleStat
	return response, nil
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

func (s *Service) ListAccountsStats(ctx context.Context, req *pb.ListStatsRequest) (*pb.ListStatsResponse, error) {
	return s.genericListStats(ctx, req, ACCOUNTS)
}

func (s *Service) ListNewsStats(ctx context.Context, req *pb.ListStatsRequest) (*pb.ListStatsResponse, error) {
	return s.genericListStats(ctx, req, NEWS)
}

func (s *Service) ListCommentsStats(ctx context.Context, req *pb.ListStatsRequest) (*pb.ListStatsResponse, error) {
	return s.genericListStats(ctx, req, COMMENTS)
}

func (s *Service) AddAccountsStats(ctx context.Context, req *pb.AddStatsRequest) (*pb.StatsResponse, error) {
	return s.genericAddStats(ctx, req, ACCOUNTS)
}

func (s *Service) AddNewsStats(ctx context.Context, req *pb.AddStatsRequest) (*pb.StatsResponse, error) {
	return s.genericAddStats(ctx, req, NEWS)
}

func (s *Service) AddCommentsStats(ctx context.Context, req *pb.AddStatsRequest) (*pb.StatsResponse, error) {
	return s.genericAddStats(ctx, req, COMMENTS)
}
