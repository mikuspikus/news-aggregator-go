package accounts

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	ststorage "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	rtstorage "github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/refresh-token-storage"
	utstorage "github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/user-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusInvalidUUID            = status.Error(codes.InvalidArgument, "invalid UUID")
	statusNotFound               = status.Error(codes.NotFound, "not found by UUID")
	statusAppIDNotFound          = status.Error(codes.NotFound, "app ID not found")
	statusInvalidSecret          = status.Error(codes.InvalidArgument, "invalid secret")
	statusInvalidServiceToken    = status.Error(codes.Unauthenticated, "invalid service token")
	statusInvalidUserCredentials = status.Error(codes.Unauthenticated, "invalid username or password")
	statusInvalidToken           = status.Error(codes.InvalidArgument, "invalid token")
	statusInvalidRefreshToken    = status.Error(codes.InvalidArgument, "invalid refresh token")
)

func internalServerError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func (user *User) UserInfo() (*pb.UserInfo, error) {
	created, err := ptypes.TimestampProto(user.Created)
	if err != nil {
		return nil, err
	}
	edited, err := ptypes.TimestampProto(user.Edited)
	if err != nil {
		return nil, err
	}

	userinfo := new(pb.UserInfo)
	userinfo.Uid = user.Uid.String()
	userinfo.Username = user.Username
	userinfo.Created = created
	userinfo.Edited = edited

	return userinfo, nil
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserInfo, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	user, err := s.Store.Get(uid)
	if err != nil {
		return nil, internalServerError(err)
	}

	response, err := user.UserInfo()
	if err != nil {
		return nil, internalServerError(err)
	}
	return response, nil
}

func (s *Service) CreateTokens(uid uuid.UUID) (string, string, error) {
	token, err := s.UserTokens.Add(uid)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.RefreshTokens.Add(token)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (s *Service) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.UserInfo, error) {
	apiToken := req.ApiToken
	valid, err := s.ServiceTokens.CheckToken(apiToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}
	username := req.Username
	password := req.Password

	user, err := s.Store.Create(username, password)
	if err != nil {
		return nil, internalServerError(err)
	}
	response, err := user.UserInfo()
	if err != nil {
		return nil, internalServerError(err)
	}
	return response, nil
}

func (s *Service) EditUser(ctx context.Context, req *pb.EditUserRequest) (*pb.UserInfo, error) {
	apiToken := req.ApiToken
	valid, err := s.ServiceTokens.CheckToken(apiToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}
	user, err := s.Store.Update(uid, req.Username, req.Password)
	switch err {
	case nil:
		userinfo, err := user.UserInfo()
		if err != nil {
			return nil, internalServerError(err)
		}
		return userinfo, nil
	case errNotFound:
		return nil, statusInvalidUUID
	default:
		return nil, internalServerError(err)

	}
}

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	apiToken := req.ApiToken
	valid, err := s.ServiceTokens.CheckToken(apiToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.Store.Delete(uid)

	switch err {
	case nil:
		return new(pb.DeleteUserResponse), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalServerError(err)
	}
}

func (s *Service) GetServiceToken(ctx context.Context, req *pb.GetServiceTokenRequest) (*pb.GetServiceTokenResponse, error) {
	appID, appSECRET := req.AppID, req.AppSECRET

	token, err := s.ServiceTokens.AddToken(appID, appSECRET)
	switch err {
	case nil:
		response := new(pb.GetServiceTokenResponse)
		response.Token = token
		return response, nil

	case ststorage.ErrIDNotFound:
		return nil, statusAppIDNotFound

	case ststorage.ErrWrongSecret:
		return nil, statusInvalidSecret

	default:
		return nil, internalServerError(err)
	}
}

func (s *Service) GetUserToken(ctx context.Context, req *pb.GetUserTokenRequest) (*pb.UserTokenResponse, error) {
	apiToken := req.ApiToken
	valid, err := s.ServiceTokens.CheckToken(apiToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}

	username, password := req.Username, req.Password
	user, err := s.Store.GetUserByUsername(username)
	if err == errNotFound {
		return nil, statusNotFound
	} else if err != nil {
		return nil, internalServerError(err)
	}
	valid, err = s.Store.CheckPassword(user.Uid, password)
	if err == errNotFound {
		return nil, statusNotFound
	} else if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidUserCredentials
	}

	token, refreshToken, err := s.CreateTokens(user.Uid)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.UserTokenResponse)
	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *Service) RefreshUserToken(ctx context.Context, req *pb.RefreshUserTokenRequest) (*pb.UserTokenResponse, error) {
	apiToken := req.ApiToken
	valid, err := s.ServiceTokens.CheckToken(apiToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}

	token, refreshToken := req.Token, req.RefreshToken
	valid, err = s.RefreshTokens.Check(token, refreshToken)
	if err == rtstorage.ErrNotFound {
		return nil, statusInvalidRefreshToken
	} else if err != nil {
		return nil, internalServerError(err)
	}
	if !valid {
		return nil, statusInvalidToken
	}

	struid, err := s.RefreshTokens.Get(refreshToken)
	if err != nil {
		return nil, internalServerError(err)
	}
	uid, err := uuid.Parse(struid)
	if err != nil {
		return nil, internalServerError(err)
	}

	err = s.RefreshTokens.Del(refreshToken)
	if err != nil {
		return nil, internalServerError(err)
	}

	token, refreshToken, err = s.CreateTokens(uid)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.UserTokenResponse)
	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *Service) GetUserByToken(ctx context.Context, req *pb.GetUserByTokenRequest) (*pb.UserInfo, error) {
	apiToken := req.ApiToken
	valid, err := s.ServiceTokens.CheckToken(apiToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidServiceToken
	}

	token := req.UserToken
	struid, err := s.UserTokens.Get(token)
	if err == utstorage.ErrNotFound {
		return nil, statusInvalidToken
	} else if err != nil {
		return nil, internalServerError(err)
	}
	uid, err := uuid.Parse(struid)
	if err != nil {
		return nil, internalServerError(err)
	}
	user, err := s.Store.Get(uid)
	if err == errNotFound {
		return nil, statusNotFound
	} else if err != nil {
		return nil, internalServerError(err)
	}

	response, err := user.UserInfo()
	if err != nil {
		return nil, internalServerError(err)
	}
	return response, nil
}
