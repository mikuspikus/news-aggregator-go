package accounts

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	stst "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	rtst "github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/refresh-token-storage"
	utst "github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/user-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusInvalidUUID         = status.Error(codes.InvalidArgument, "invalid UUID")
	statusInvalidToken        = status.Error(codes.InvalidArgument, "invalid token")
	statusInvalidRefreshToken = status.Error(codes.InvalidArgument, "invalid refresh token")
	statusInvalidSecret       = status.Error(codes.InvalidArgument, "invalid secret")

	statusInvalidServiceToken    = status.Error(codes.Unauthenticated, "invalid service token")
	statusInvalidUserCredentials = status.Error(codes.Unauthenticated, "invalid username or password")
	statusServiceTokenNotFound   = status.Error(codes.Unauthenticated, "service token not found")

	statusNotFoundByUUID      = status.Error(codes.NotFound, "not found by UUID")
	statusUserNotFoundByToken = status.Error(codes.NotFound, "user not found by token")
	statusAppIDNotFound       = status.Error(codes.NotFound, "app ID not found")
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

func (s *Service) validateServiceToken(token string) error {
	valid, err := s.ServiceTokens.CheckToken(token)
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

func (s *Service) createUserTokens(uid uuid.UUID) (string, string, error) {
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
	statusError := s.validateServiceToken(req.ApiToken)
	if statusError != nil {
		return nil, statusError
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
	statusError := s.validateServiceToken(req.ApiToken)
	if statusError != nil {
		return nil, statusError
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
	statusError := s.validateServiceToken(req.ApiToken)
	if statusError != nil {
		return nil, statusError
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
		return nil, statusNotFoundByUUID
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

	case stst.ErrIDNotFound:
		return nil, statusAppIDNotFound

	case stst.ErrWrongSecret:
		return nil, statusInvalidSecret

	default:
		return nil, internalServerError(err)
	}
}

func (s *Service) GetUserToken(ctx context.Context, req *pb.GetUserTokenRequest) (*pb.UserTokenResponse, error) {
	statusError := s.validateServiceToken(req.ApiToken)
	if statusError != nil {
		return nil, statusError
	}

	username, password := req.Username, req.Password
	user, err := s.Store.GetUserByUsername(username)
	if err == errNotFound {
		return nil, statusNotFoundByUUID
	} else if err != nil {
		return nil, internalServerError(err)
	}
	valid, err := s.Store.CheckPassword(user.Uid, password)
	if err == errNotFound {
		return nil, statusNotFoundByUUID
	} else if err != nil {
		return nil, err
	}
	if !valid {
		return nil, statusInvalidUserCredentials
	}

	token, refreshToken, err := s.createUserTokens(user.Uid)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.UserTokenResponse)
	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *Service) RefreshUserToken(ctx context.Context, req *pb.RefreshUserTokenRequest) (*pb.UserTokenResponse, error) {
	statusError := s.validateServiceToken(req.ApiToken)
	if statusError != nil {
		return nil, statusError
	}

	token, refreshToken := req.Token, req.RefreshToken
	valid, err := s.RefreshTokens.Check(token, refreshToken)
	if err == rtst.ErrTokenNotFound {
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

	token, refreshToken, err = s.createUserTokens(uid)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.UserTokenResponse)
	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *Service) GetUserByToken(ctx context.Context, req *pb.GetUserByTokenRequest) (*pb.UserInfo, error) {
	statusError := s.validateServiceToken(req.ApiToken)
	if statusError != nil {
		return nil, statusError
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
