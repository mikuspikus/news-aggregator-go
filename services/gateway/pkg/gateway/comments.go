package gateway

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	comments "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
)

type Comment struct {
	ID       int
	UserUUID string
	NewsUUID string
	Body     string
	Created  time.Time
	Edited   time.Time
}

//
func extractIntFromString(stringValue string, defaultValue int) (int, error) {
	if stringValue != "" {
		value, err := strconv.Atoi(stringValue)
		if err != nil {
			return 0, err
		}
		return value, nil
	}

	return defaultValue, nil
}

// convertSingleComment converts comments.SingleComment into Comment
func convertSingleComment(singleComment *comments.SingleComment) (*Comment, error) {
	comment := new(Comment)
	var err error
	comment.ID = int(singleComment.Id)
	comment.NewsUUID = singleComment.NewsUUID
	comment.UserUUID = singleComment.UserUUID
	comment.Body = singleComment.Body
	comment.Created, err = ptypes.Timestamp(singleComment.Created)
	if err != nil {
		return nil, err
	}
	comment.Edited, err = ptypes.Timestamp(singleComment.Edited)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (cc *CommentsClient) GetComment(ctx context.Context, id int32) (*Comment, error) {
	commentsResponse, err := cc.client.GetComment(ctx, &comments.GetCommentRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return convertSingleComment(commentsResponse.Comment)
}

func (cc *CommentsClient) AddComment(ctx context.Context, body, userUUID, newsUUID string) (*Comment, error) {
	commentsResponse, err := cc.client.AddComment(ctx, &comments.AddCommentRequest{
		Body:     body,
		UserUUID: userUUID,
		NewsUUID: newsUUID,
		Token:    cc.token,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = cc.UpdateToken()
			if err != nil {
				return nil, err
			}
			commentsResponse, err = cc.client.AddComment(ctx, &comments.AddCommentRequest{
				Body:     body,
				UserUUID: userUUID,
				NewsUUID: newsUUID,
				Token:    cc.token,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return convertSingleComment(commentsResponse.Comment)
}

func (cc *CommentsClient) UpdateComment(ctx context.Context, body string) (*Comment, error) {
	commentsResponse, err := cc.client.EditComment(ctx, &comments.EditCommentRequest{
		Body:  body,
		Token: cc.token,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = cc.UpdateToken()
			if err != nil {
				return nil, err
			}
			commentsResponse, err = cc.client.EditComment(ctx, &comments.EditCommentRequest{
				Body:  body,
				Token: cc.token,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return convertSingleComment(commentsResponse.Comment)
}

func (cc *CommentsClient) DeleteComment(ctx context.Context, id int32) error {
	_, err := cc.client.DeleteComment(ctx, &comments.DeleteCommentRequest{
		Id:    id,
		Token: cc.token,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = cc.UpdateToken()
			if err != nil {
				return err
			}
			_, err := cc.client.DeleteComment(ctx, &comments.DeleteCommentRequest{
				Id:    id,
				Token: cc.token,
			})
			if err != nil {
				return err
			}
		} else {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// getNewsComments returns paged comment for news passed through Mux router
func (s *Server) getNewsComments() http.HandlerFunc {
	type Response struct {
		Comments   []Comment
		PageSize   int
		PageNumber int
		PagesCount int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		strpage, strsize := r.URL.Query().Get("page"), r.URL.Query().Get("number")
		page, err := extractIntFromString(strpage, 0)
		if err != nil {
			http.Error(w, "can't parse URL parameter 'page'", http.StatusBadRequest)
			return
		}
		size, err := extractIntFromString(strsize, 10)
		if err != nil {
			http.Error(w, "can't parse URL parameter `size`", http.StatusBadRequest)
			return
		}
		vars := mux.Vars(r)
		newsUUID := vars["news"]
		ctx := r.Context()
		grpcResponse, err := s.Comments.client.ListComments(ctx, &comments.ListCommentsRequest{
			NewsUUID:   newsUUID,
			PageNumber: int32(page),
			PageSize:   int32(size),
		})
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		comments := make([]Comment, len(grpcResponse.Comments))
		for idx, singleComment := range grpcResponse.Comments {
			comment, err := convertSingleComment(singleComment)
			if err != nil {
				handleRPCErrors(w, err)
				return
			}

			comments[idx] = *comment
		}
		httpResponse := Response{
			Comments:   comments,
			PageSize:   size,
			PageNumber: page,
			PagesCount: int(grpcResponse.PageCount),
		}
		json, err := json.Marshal(httpResponse)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

// getSingleComment returns comment data with id passed through Mux router
func (s *Server) getSingleComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := extractIntFromString(vars["id"], -1)
		if err != nil {
			http.Error(w, "can't parse 'id' route parameter", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		grpcResponse, err := s.Comments.client.GetComment(ctx, &comments.GetCommentRequest{Id: int32(id)})
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		comment, err := convertSingleComment(grpcResponse.Comment)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(*comment)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

// deleteComment deletes comment by id through Mux router from [comments]
func (s *Server) deleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := getAuthorizationToken(r)
		if userToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		user, err := s.Accounts.GetUserByToken(ctx, userToken)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		vars := mux.Vars(r)
		id, err := extractIntFromString(vars["id"], -1)
		if err != nil {
			http.Error(w, "can't parse 'id' route parameter", http.StatusBadRequest)
			return
		}

		comment, err := s.Comments.GetComment(ctx, int32(id))
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		if user.UID != comment.UserUUID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err = s.Comments.DeleteComment(ctx, int32(id))
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// updateComment changes the comment by id through Mux router
func (s *Server) updateComment() http.HandlerFunc {
	type request struct {
		Body string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := getAuthorizationToken(r)
		if userToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		user, err := s.Accounts.GetUserByToken(ctx, userToken)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		var req request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		err = json.Unmarshal(b, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		vars := mux.Vars(r)
		id, err := extractIntFromString(vars["id"], -1)
		if err != nil {
			http.Error(w, "can't parse 'id' route parameter", http.StatusBadRequest)
			return
		}

		comment, err := s.Comments.GetComment(ctx, int32(id))
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		if user.UID != comment.UserUUID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		comment, err = s.Comments.UpdateComment(ctx, req.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(*comment)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}

// createComment add new comment to [comments] with newsUUID from newsuuid from Mux router
func (s *Server) createComment() http.HandlerFunc {
	type request struct {
		Body string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := getAuthorizationToken(r)
		if userToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		user, err := s.Accounts.GetUserByToken(ctx, userToken)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		var req request
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		err = json.Unmarshal(b, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		vars := mux.Vars(r)
		newsUUID := vars["newsuuid"]

		comment, err := s.Comments.AddComment(ctx, req.Body, user.UID, newsUUID)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(*comment)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}
