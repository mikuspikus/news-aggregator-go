package stats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusInvalidUUID         = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidToken        = status.Error(codes.Unauthenticated, "invalid token")
	statusAppIDNotFound       = status.Error(codes.NotFound, "app ID not found")
	statusInvalidSecret       = status.Error(codes.InvalidArgument, "invalid secret")
	statusInvalidServiceToken = status.Error(codes.Unauthenticated, "invalid service token")
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
	tokenStorage *simple_token_storage.APITokenStorage
}

func (s *Service) genericAddStats(ctx context.Context, r *pb.AddStatsRequest, target target) (*pb.StatsResponse, error) {
	apiToken := r.Token
	valid, err := s.tokenStorage.CheckToken(apiToken)
	if err == simple_token_storage.ErrTokenNotFound {
		return nil, statusInvalidServiceToken
	} else if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}

	var userUID uuid.UUID
	if r.UserUID == "" {
		userUID = uuid.Nil
	} else {
		userUID, err = uuid.Parse(r.UserUID)
		if err != nil {
			return nil, statusInvalidUUID
		}
	}

	var input, output map[string]interface{}
	err = unmarshal(r.Input, &input)
	if err != nil {
		return nil, internalServerError(err)
	}
	err = unmarshal(r.Output, &output)
	if err != nil {
		return nil, internalServerError(err)
	}

	stats := new(Stats)
	switch target {
	case ACCOUNTS:
		stats, err = s.db.AddAccountsStats(userUID, r.Action, input, output)
	case NEWS:
		stats, err = s.db.AddNewsStats(userUID, r.Action, input, output)
	case COMMENTS:
		stats, err = s.db.AddCommentsStats(userUID, r.Action, input, output)
	default:
		return nil, internalServerError(errors.New("unknown type"))
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
	case simple_token_storage.ErrIDNotFound:
		return nil, statusAppIDNotFound
	case simple_token_storage.ErrWrongSecret:
		return nil, statusInvalidSecret
	default:
		return nil, internalServerError(err)
	}
}

func (s *Service) ListAccountsStats(ctx context.Context, req *pb.ListStatsRequest) (*pb.ListStatsResponse, error) {
	valid, err := s.tokenStorage.CheckToken(req.Token)
	if err != nil {
		return nil, internalServerError(err)
	}
	if !valid {
		return nil, statusInvalidToken
	}

	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	stats, pageCount, err := s.db.ListAccountsStats(req.PageNumber, pageSize)
	if err != nil {
		return nil, err
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

func (s *Service) ListNewsStats(ctx context.Context, req *pb.ListStatsRequest) (*pb.ListStatsResponse, error) {
	valid, err := s.tokenStorage.CheckToken(req.Token)
	if err != nil {
		return nil, internalServerError(err)
	}
	if !valid {
		return nil, statusInvalidToken
	}

	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	stats, pageCount, err := s.db.ListNews(req.PageNumber, pageSize)
	if err != nil {
		return nil, err
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

func (s *Service) ListCommentsStats(ctx context.Context, req *pb.ListStatsRequest) (*pb.ListStatsResponse, error) {
	valid, err := s.tokenStorage.CheckToken(req.Token)
	if err != nil {
		return nil, internalServerError(err)
	}
	if !valid {
		return nil, statusInvalidToken
	}

	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	stats, pageCount, err := s.db.ListComments(req.PageNumber, pageSize)
	if err != nil {
		return nil, err
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

func (s *Service) AddAccountsStats(ctx context.Context, req *pb.AddStatsRequest) (*pb.StatsResponse, error) {
	return s.genericAddStats(ctx, req, ACCOUNTS)
}

func (s *Service) AddNewsStats(ctx context.Context, req *pb.AddStatsRequest) (*pb.StatsResponse, error) {
	return s.genericAddStats(ctx, req, NEWS)
}

func (s *Service) AddCommentsStats(ctx context.Context, req *pb.AddStatsRequest) (*pb.StatsResponse, error) {
	return s.genericAddStats(ctx, req, COMMENTS)
}
