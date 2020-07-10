package gateway

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
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

type User struct {
	UID      string
	Username string
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

// convertUserInfo converts accounts.UserInfo into User
func convertUserInfo(userInfo *accounts.UserInfo) (*User, error) {
	user := new(User)
	var err error
	user.UID = userInfo.Uid
	user.Username = userInfo.Username
	user.Created, err = ptypes.Timestamp(userInfo.Created)
	if err != nil {
		return nil, err
	}
	user.Edited, err = ptypes.Timestamp(userInfo.Edited)
	if err != nil {
		return nil, err
	}
	return user, nil
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

		accountsResponse, err := s.Accounts.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
			ApiToken:  s.Accounts.token,
			UserToken: userToken,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Accounts.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				accountsResponse, err = s.Accounts.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
					ApiToken:  s.Accounts.token,
					UserToken: userToken,
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}
		user, err := convertUserInfo(accountsResponse)
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
		commentsResponse, err := s.Comments.client.GetComment(ctx, &comments.GetCommentRequest{
			Id: int32(id),
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Comments.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				commentsResponse, err = s.Comments.client.GetComment(ctx, &comments.GetCommentRequest{
					Id: int32(id),
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}
		comment, err := convertSingleComment(commentsResponse.Comment)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		if user.UID != comment.UserUUID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		_, err = s.Comments.client.DeleteComment(ctx, &comments.DeleteCommentRequest{
			Id:    int32(id),
			Token: s.Comments.token,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Comments.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				_, err = s.Comments.client.DeleteComment(ctx, &comments.DeleteCommentRequest{
					Id:    int32(id),
					Token: s.Comments.token,
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
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

		accountsResponse, err := s.Accounts.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
			ApiToken:  s.Accounts.token,
			UserToken: userToken,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Accounts.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				accountsResponse, err = s.Accounts.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
					ApiToken:  s.Accounts.token,
					UserToken: userToken,
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}
		user, err := convertUserInfo(accountsResponse)
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
		commentsResponse, err := s.Comments.client.GetComment(ctx, &comments.GetCommentRequest{
			Id: int32(id),
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Comments.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				commentsResponse, err = s.Comments.client.GetComment(ctx, &comments.GetCommentRequest{
					Id: int32(id),
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}
		comment, err := convertSingleComment(commentsResponse.Comment)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		if user.UID != comment.UserUUID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		updatedResponse, err := s.Comments.client.EditComment(ctx, &comments.EditCommentRequest{
			Token: s.Comments.token,
			Id:    int32(id),
			Body:  req.Body,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Comments.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				updatedResponse, err = s.Comments.client.EditComment(ctx, &comments.EditCommentRequest{
					Token: s.Comments.token,
					Id:    int32(id),
					Body:  req.Body,
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}

		comment, err = convertSingleComment(updatedResponse.Comment)
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

		accountsResponse, err := s.Accounts.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
			ApiToken:  s.Accounts.token,
			UserToken: userToken,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Accounts.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				accountsResponse, err = s.Accounts.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
					ApiToken:  s.Accounts.token,
					UserToken: userToken,
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}
		user, err := convertUserInfo(accountsResponse)
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

		commentsResponse, err := s.Comments.client.AddComment(ctx, &comments.AddCommentRequest{
			Body:     req.Body,
			UserUUID: user.UID,
			NewsUUID: newsUUID,
			Token:    s.Comments.token,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
				err = s.Accounts.UpdateToken()
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
				commentsResponse, err = s.Comments.client.AddComment(ctx, &comments.AddCommentRequest{
					Body:     req.Body,
					UserUUID: user.UID,
					NewsUUID: newsUUID,
					Token:    s.Comments.token,
				})
				if err != nil {
					handleRPCErrors(w, err)
					return
				}
			} else {
				handleRPCErrors(w, err)
				return
			}
		}

		comment, err := convertSingleComment(commentsResponse.Comment)
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
