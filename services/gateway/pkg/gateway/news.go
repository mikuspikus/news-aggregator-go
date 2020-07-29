package gateway

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	news "github.com/mikuspikus/news-aggregator-go/services/news/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"time"
)

type News struct {
	UID     string    `json:"uid"`
	User    string    `json:"user"`
	Title   string    `json:"title"`
	URI     string    `json:"uri"`
	Created time.Time `json:"created"`
	Edited  time.Time `json:"edited"`
}

func convertSingleNews(singleNews *news.SingleNews) (*News, error) {
	news := new(News)
	var err error

	news.UID = singleNews.Uid
	news.User = singleNews.UserUUID
	news.Title = singleNews.Title
	news.URI = singleNews.Uri
	news.Created, err = ptypes.Timestamp(singleNews.Created)
	if err != nil {
		return nil, err
	}
	news.Edited, err = ptypes.Timestamp(singleNews.Edited)
	if err != nil {
		return nil, err
	}
	return news, err
}

func (nc *NewsClient) ListNews(ctx context.Context, pageNumber, pageSize int) ([]News, int, error) {
	response, err := nc.client.ListNews(ctx, &news.ListNewsRequest{
		PageNumber: int32(pageNumber),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		return nil, 0, err
	}
	newses := make([]News, len(response.News))
	for idx, singleNews := range response.News {
		news, err := convertSingleNews(singleNews)
		if err != nil {
			return nil, 0, err
		}
		newses[idx] = *news
	}
	return newses, int(response.PageCount), nil
}

func (nc *NewsClient) GetNews(ctx context.Context, uid string) (*News, error) {
	response, err := nc.client.GetNews(ctx, &news.GetNewsRequest{
		Uid: uid,
	})
	if err != nil {
		return nil, err
	}
	return convertSingleNews(response.News)
}

func (nc *NewsClient) AddNews(ctx context.Context, user, title, uri string) (*News, error) {
	response, err := nc.client.AddNews(ctx, &news.AddNewsRequest{
		Token:    nc.token,
		UserUUID: user,
		Title:    title,
		Uri:      uri,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = nc.UpdateToken(ctx)
			if err != nil {
				return nil, err
			}
			response, err = nc.client.AddNews(ctx, &news.AddNewsRequest{
				Token:    nc.token,
				UserUUID: user,
				Title:    title,
				Uri:      uri,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return convertSingleNews(response.News)
}

func (nc *NewsClient) UpdateNews(ctx context.Context, uid, title, uri string) (*News, error) {
	response, err := nc.client.EditNews(ctx, &news.EditNewsRequest{
		Token: nc.token,
		Uid:   uid,
		Title: title,
		Uri:   uri,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = nc.UpdateToken(ctx)
			if err != nil {
				return nil, err
			}
			response, err = nc.client.EditNews(ctx, &news.EditNewsRequest{
				Token: nc.token,
				Uid:   uid,
				Title: title,
				Uri:   uri,
			})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return convertSingleNews(response.News)
}

func (nc *NewsClient) DeleteNews(ctx context.Context, uid string) error {
	_, err := nc.client.DeleteNews(ctx, &news.DeleteNewsRequest{
		Token: nc.token,
		Uid:   uid,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = nc.UpdateToken(ctx)
			if err != nil {
				return err
			}
			_, err := nc.client.DeleteNews(ctx, &news.DeleteNewsRequest{
				Token: nc.token,
				Uid:   uid,
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

func (s *Server) getNews() http.HandlerFunc {
	type Response struct {
		News       []News `json:"news"`
		PageSize   int    `json:"page_size"`
		PageNumber int    `json:"page_number"`
		PageCount  int    `json:"page_count"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userToken := getAuthorizationToken(r)
		msg := new(Message)
		msg.UserUID = ""
		msg.Action = "getNews"
		if user, err := s.Accounts.GetUserByToken(ctx, userToken); err == nil {
			msg.UserUID = user.UID
		}

		strpage, strsize := r.URL.Query().Get("page"), r.URL.Query().Get("size")
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
		news, pageCount, err := s.News.ListNews(ctx, page, size)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		id_s := make([]string, len(news))
		for idx, new := range news {
			id_s[idx] = new.UID
		}

		input, err := json.Marshal(map[string][]string{"id_s": id_s})
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		msg.Input = input
		msg_, err := json.Marshal(msg)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		err = s.News.Writer.Push(ctx, nil, msg_)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		response := Response{
			News:       news,
			PageSize:   size,
			PageNumber: page,
			PageCount:  pageCount,
		}
		json, err := json.Marshal(&response)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) getSingleNews() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["newsuid"]
		ctx := r.Context()

		news, err := s.News.GetNews(ctx, uid)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(news)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) deleteNews() http.HandlerFunc {
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
		uid := vars["newsuid"]

		news, err := s.News.GetNews(ctx, uid)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}

		if user.UID != news.User {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err = s.News.DeleteNews(ctx, uid)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) updateNews() http.HandlerFunc {
	type Request struct {
		Title string `json:"title"`
		URI   string `json:"uri"`
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
		vars := mux.Vars(r)
		uid := vars["newsuid"]
		news, err := s.News.GetNews(ctx, uid)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		if user.UID != news.User {
			w.WriteHeader(http.StatusForbidden)
			return
		}

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

		news, err = s.News.UpdateNews(ctx, uid, req.Title, req.URI)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(news)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}

func (s *Server) addNews() http.HandlerFunc {
	type Request struct {
		Title string `json:"title"`
		URI   string `json:"uri"`
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

		news, err := s.News.AddNews(ctx, user.UID, req.Title, req.URI)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		json, err := json.Marshal(news)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(json)
	}
}
