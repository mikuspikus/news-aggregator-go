package gateway

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	accounts "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	UID      string    `json:"uid"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
	Edited   time.Time `json:"edited"`
	IsAdmin  bool      `json:"is_admin"`
}

// convertUserInfo converts accounts.UserInfo into User
func convertUserInfo(userInfo *accounts.UserInfo) (*User, error) {
	user := new(User)
	var err error
	user.UID = userInfo.Uid
	user.Username = userInfo.Username
	user.IsAdmin = userInfo.IsAdmin
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

func (ac *AccountsClient) GetUserByToken(ctx context.Context, userToken string) (*User, error) {
	accountsResponse, err := ac.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
		ApiToken:  ac.token,
		UserToken: userToken,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return nil, err
			}
			accountsResponse, err = ac.client.GetUserByToken(ctx, &accounts.GetUserByTokenRequest{
				ApiToken:  ac.token,
				UserToken: userToken,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return convertUserInfo(accountsResponse)
}

func (ac *AccountsClient) GetUser(ctx context.Context, userUUID string) (*User, error) {
	response, err := ac.client.GetUser(ctx, &accounts.GetUserRequest{Uid: userUUID})
	if err != nil {
		return nil, err
	}
	return convertUserInfo(response)
}

func (ac *AccountsClient) ListUser(ctx context.Context, pageNumber, pageSize int) ([]*User, int, error) {
	request := &accounts.ListUsersRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
	}
	response, err := ac.client.ListUsers(ctx, request)
	if err != nil {
		return nil, 0, err
	}
	users := make([]*User, len(response.Users))
	for idx, userinfo := range response.Users {
		user, err := convertUserInfo(userinfo)
		if err != nil {
			return nil, 0, err
		}
		users[idx] = user
	}
	return users, int(response.PageCount), nil
}

func (ac *AccountsClient) AddUser(ctx context.Context, username, password string) (*User, error) {
	response, err := ac.client.AddUser(ctx, &accounts.AddUserRequest{
		ApiToken: ac.token,
		Username: username,
		Password: password,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return nil, err
			}
			response, err = ac.client.AddUser(ctx, &accounts.AddUserRequest{
				ApiToken: ac.token,
				Username: username,
				Password: password,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return convertUserInfo(response)
}

func (ac *AccountsClient) EditUser(ctx context.Context, username, password string) (*User, error) {
	response, err := ac.client.EditUser(ctx, &accounts.EditUserRequest{
		ApiToken: ac.token,
		Username: username,
		Password: password,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return nil, err
			}
			response, err = ac.client.EditUser(ctx, &accounts.EditUserRequest{
				ApiToken: ac.token,
				Username: username,
				Password: password,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return convertUserInfo(response)
}

func (ac *AccountsClient) AdminEditUser(ctx context.Context, uid, username string, isAdmin bool) (*User, error) {
	request := &accounts.AdminEditUserRequest{
		ApiToken: ac.token,
		Uid:      uid,
		IsAdmin:  isAdmin,
		Username: username,
	}

	response, err := ac.client.AdminEditUser(ctx, request)
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return nil, err
			}

			request.ApiToken = ac.token

			response, err = ac.client.AdminEditUser(ctx, request)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return convertUserInfo(response)
}

func (ac *AccountsClient) DeleteUser(ctx context.Context, userUUID string) error {
	_, err := ac.client.DeleteUser(ctx, &accounts.DeleteUserRequest{
		ApiToken: ac.token,
		Uid:      userUUID,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return err
			}
			_, err = ac.client.DeleteUser(ctx, &accounts.DeleteUserRequest{
				ApiToken: ac.token,
				Uid:      userUUID,
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// User authentication
func (ac *AccountsClient) GetUserToken(ctx context.Context, username, password string) (string, string, error) {
	response, err := ac.client.GetUserToken(ctx, &accounts.GetUserTokenRequest{
		ApiToken: ac.token,
		Username: username,
		Password: password,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return "", "", err
			}
			response, err = ac.client.GetUserToken(ctx, &accounts.GetUserTokenRequest{
				ApiToken: ac.token,
				Username: username,
				Password: password,
			})
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}

	return response.Token, response.RefreshToken, nil
}

func (ac *AccountsClient) RefreshUserToken(ctx context.Context, token, refreshToken string) (string, string, error) {
	response, err := ac.client.RefreshUserToken(ctx, &accounts.RefreshUserTokenRequest{
		ApiToken:     ac.token,
		Token:        token,
		RefreshToken: refreshToken,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = ac.UpdateToken(ctx)
			if err != nil {
				return "", "", err
			}
			response, err = ac.client.RefreshUserToken(ctx, &accounts.RefreshUserTokenRequest{
				ApiToken:     ac.token,
				Token:        token,
				RefreshToken: refreshToken,
			})
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", err
		}
	}
	return response.Token, response.RefreshToken, nil
}

// User administration

func (s *Server) listUsers() http.HandlerFunc {
	type Response struct {
		Users      []*User `json:"users"`
		PagesCount int     `json:"pages_count"`
		PageNumber int     `json:"page_number"`
		PageSize   int     `json:"page_size"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := getAuthorizationToken(r)
		if userToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

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

		ctx := r.Context()

		users, pageCount, err := s.Accounts.ListUser(ctx, page, size)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		httpResponse := Response{
			Users:      users,
			PageSize:   size,
			PageNumber: page,
			PagesCount: pageCount,
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

func (s *Server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		ctx := r.Context()

		user, err := s.Accounts.GetUser(ctx, uid)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		json, err := json.Marshal(user)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) addUser() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		ctx := r.Context()

		user, err := s.Accounts.AddUser(ctx, req.Username, req.Password)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(user)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}

func (s *Server) updateUser() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := getAuthorizationToken(r)
		if userToken != "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req request
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		err = json.Unmarshal(bytes, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		ctx := r.Context()

		user, err := s.Accounts.GetUserByToken(ctx, userToken)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		vars := mux.Vars(r)
		userUUID := vars["useruuid"]

		if userUUID != user.UID {
			w.WriteHeader(http.StatusForbidden)
		}

		user, err = s.Accounts.EditUser(ctx, req.Username, req.Password)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		json, err := json.Marshal(*user)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}

func (s *Server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := getAuthorizationToken(r)
		if userToken != "" {
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
		userUUID := vars["useruuid"]

		if userUUID != user.UID {
			w.WriteHeader(http.StatusForbidden)
		}
		err = s.Accounts.DeleteUser(ctx, userUUID)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) getUserToken() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type Response struct {
		Uid          string `json:"uid"`
		Username     string `json:"username"`
		IsAdmin      bool   `json:"is_admin"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		ctx := r.Context()

		token, refreshToken, err := s.Accounts.GetUserToken(ctx, req.Username, req.Password)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		user, err := s.Accounts.GetUserByToken(ctx, token)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		json, err := json.Marshal(&Response{
			Uid:          user.UID,
			Username:     user.Username,
			IsAdmin:      user.IsAdmin,
			Token:        token,
			RefreshToken: refreshToken,
		})
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) refreshUserToken() http.HandlerFunc {
	type request struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	type response struct {
		Uid          string `json:"uid"`
		Username     string `json:"username"`
		IsAdmin      bool   `json:"is_admin"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		err = json.Unmarshal(bytes, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		ctx := r.Context()

		token, refreshToken, err := s.Accounts.RefreshUserToken(ctx, req.Token, req.RefreshToken)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		user, err := s.Accounts.GetUserByToken(ctx, token)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		json, err := json.Marshal(&response{
			Uid:          user.UID,
			Username:     user.Username,
			IsAdmin:      user.IsAdmin,
			Token:        token,
			RefreshToken: refreshToken,
		})
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) adminEditUser() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		IsAdmin  bool   `json:"is_admin"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//userToken := getAuthorizationToken(r)
		//if userToken == "" {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	return
		//}

		ctx := r.Context()
		//
		//user, err := s.Accounts.GetUserByToken(ctx, userToken)
		//if err != nil {
		//	handleRPCErrors(w, err)
		//	return
		//}
		//if !user.IsAdmin {
		//	w.WriteHeader(http.StatusForbidden)
		//	return
		//}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		var request Request
		err = json.Unmarshal(bytes, &request)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		vars := mux.Vars(r)
		uid := vars["useruid"]

		user, err := s.Accounts.AdminEditUser(ctx, uid, request.Username, request.IsAdmin)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		json, err := json.Marshal(user)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}
