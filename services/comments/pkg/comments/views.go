package comments

import (

	// Import the generated protobuf code
	"context"

	pb "github.com/mikuspikus/news-aggregator-go/services/comments/proto"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusInvalidUUID  = status.Error(codes.InvalidArgument, "Invalid UUID")
	statusNotFound     = status.Error(codes.NotFound, "Comment not found")
	statusInvalidToken = status.Error(codes.Unauthenticated, "Invalid authentication token")
)

func internalServerError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

// SingleComment conerts Comment structure into SingleComment from pb
func (comment *Comment) SingleComment() (*pb.SingleComment, error) {
	created, err := ptypes.TimestampProto(comment.Created)
	if err != nil {
		return nil, internalServerError(err)
	}

	edited, err := ptypes.TimestampProto(comment.Edited)
	if err != nil {
		return nil, internalServerError(err)
	}

	singleComment := new(pb.SingleComment)
	singleComment.Id = comment.ID
	singleComment.Body = comment.Body
	singleComment.NewsUUID = comment.News.String()
	singleComment.UserUUID = comment.User.String()
	singleComment.Created = created
	singleComment.Edited = edited

	return singleComment, nil
}

// Service structure implements gRPC interface for Comments Service
type Service struct {
	db DataStoreHandler
}

// GetToken returns new authorization token for appID and appSECRET
func (s *Service) GetToken(ctx context.Context, req *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	// appID, appSECRET := req.AppID, req.AppSECRET

	// token, err := s.auth.Add(appID, appSECRET)
	// switch err {
	// case nil:
	// 	response := new(pb.GetTokenResponse)
	// 	response.Token = response.Token
	// 	return response, nil

	// case auth.ErrorNotFound:
	// 	return nil, statusNotFound

	// case auth.ErrorWrongCredentials:
	// 	return nil, statusWrongCredentials

	// default:
	// 	return nil, internalServerError(err)
	// }

	return new(pb.GetTokenResponse), nil
}

// DeleteComment deletes comment by ID
func (s *Service) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	// valid, err := s.CheckToken(req.Token)
	// if err != nil {
	// 	return nil, err
	// }
	// if !valid {
	// 	return nil, statusInvalidToken
	// }
	err := s.db.Delete(req.Id)

	switch err {
	case nil:
		return new(pb.DeleteCommentResponse), nil

	case errNotFound:
		return nil, statusNotFound

	default:
		return nil, internalServerError(err)
	}
}

// EditComment changes comment Body by ID
func (s *Service) EditComment(ctx context.Context, req *pb.EditCommentRequest) (*pb.EditCommentResponse, error) {
	// valid, err := s.CheckToken(req.Token)
	// if err != nil {
	// 	return nil, err
	// }
	// if !valid {
	// 	return nil, statusInvalidToken
	// }
	comment, err := s.db.Update(req.Id, req.Body)

	if err != nil {
		if err == errNotFound {
			return nil, statusNotFound
		}

		return nil, internalServerError(err)

	}

	singleComment, err := comment.SingleComment()
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.EditCommentResponse)
	response.Comment = singleComment

	return response, nil
}

// AddComment adds new comment
func (s *Service) AddComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.AddCommentResponse, error) {
	// valid, err := s.CheckToken(req.Token)
	// if err != nil {
	// 	return nil, err
	// }
	// if !valid {
	// 	return nil, statusInvalidToken
	// }

	newsUUID, err := uuid.Parse(req.NewsUUID)
	if err != nil {
		return nil, statusInvalidUUID
	}

	userUUID, err := uuid.Parse(req.UserUUID)
	if err != nil {
		return nil, statusInvalidUUID
	}

	comment, err := s.db.Create(userUUID, newsUUID, req.Body)
	if err != nil {
		return nil, internalServerError(err)
	}

	singleComment, err := comment.SingleComment()
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.AddCommentResponse)
	response.Comment = singleComment

	return response, nil
}

// GetComment returns comment info by ID
func (s *Service) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	comment, err := s.db.Get(req.Id)

	if err != nil {
		return nil, statusInvalidUUID
	}

	singleComment, err := comment.SingleComment()

	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.GetCommentResponse)
	response.Comment = singleComment

	return response, nil
}

// ListComments returns several comments with pageSize and PageNumber
func (s *Service) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	var pageSize int32

	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	var newsUUID uuid.UUID
	if req.NewsUUID == "" {
		newsUUID = uuid.Nil
	}

	comments, pageCount, err := s.db.List(req.PageNumber, pageSize, newsUUID)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := new(pb.ListCommentsResponse)

	for _, comment := range comments {
		singleComment, err := comment.SingleComment()
		if err != nil {
			return nil, internalServerError(err)
		}

		response.Comments = append(response.Comments, singleComment)
	}

	response.PageNumber = req.PageNumber
	response.PageSize = pageSize
	response.PageCount = pageCount

	return response, nil
}
