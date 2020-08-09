package gateway

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type Stats struct {
	ID        int                    `json:"id"`
	User      string                 `json:"user"`
	Action    string                 `json:"action"`
	Timestamp time.Time              `json:"timestamp"`
	Input     map[string]interface{} `json:"input"`
	Output    map[string]interface{} `json:"output"`
}

func convertSingleStat(singleStat *stats.SingleStat) (*Stats, error) {
	stats := new(Stats)
	var err error

	stats.ID = int(singleStat.Id)
	stats.User = singleStat.UserUID
	stats.Action = singleStat.Action
	stats.Timestamp, err = ptypes.Timestamp(singleStat.Timestamp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(singleStat.Input, &stats.Input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(singleStat.Output, &stats.Output)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (sc *StatsClient) ListAccountsStats(ctx context.Context, pageNumber, pageSize int) ([]Stats, int, error) {
	response, err := sc.client.ListAccountsStats(ctx, &stats.ListStatsRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
		Token:      sc.token,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = sc.UpdateToken(ctx)
			if err != nil {
				return nil, 0, err
			}
			response, err = sc.client.ListAccountsStats(ctx, &stats.ListStatsRequest{
				PageSize:   int32(pageSize),
				PageNumber: int32(pageNumber),
				Token:      sc.token,
			})
			if err != nil {
				return nil, 0, err
			}
		} else {
			return nil, 0, err
		}
	}
	stat_s := make([]Stats, len(response.Stats))
	for idx, singleStat := range response.Stats {
		stat, err := convertSingleStat(singleStat)
		if err != nil {
			return nil, 0, err
		}
		stat_s[idx] = *stat
	}
	return stat_s, int(response.PageCount), nil
}

func (sc *StatsClient) ListNewsStats(ctx context.Context, pageNumber, pageSize int) ([]Stats, int, error) {
	response, err := sc.client.ListNewsStats(ctx, &stats.ListStatsRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
		Token:      sc.token,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = sc.UpdateToken(ctx)
			if err != nil {
				return nil, 0, err
			}
			response, err = sc.client.ListNewsStats(ctx, &stats.ListStatsRequest{
				PageSize:   int32(pageSize),
				PageNumber: int32(pageNumber),
				Token:      sc.token,
			})
			if err != nil {
				return nil, 0, err
			}
		} else {
			return nil, 0, err
		}
	}
	stat_s := make([]Stats, len(response.Stats))
	for idx, singleStat := range response.Stats {
		stat, err := convertSingleStat(singleStat)
		if err != nil {
			return nil, 0, err
		}
		stat_s[idx] = *stat
	}
	return stat_s, int(response.PageCount), nil
}

func (sc *StatsClient) ListCommentsStats(ctx context.Context, pageNumber, pageSize int) ([]Stats, int, error) {
	response, err := sc.client.ListCommentsStats(ctx, &stats.ListStatsRequest{
		PageSize:   int32(pageSize),
		PageNumber: int32(pageNumber),
		Token:      sc.token,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok && status.Code() == codes.Unauthenticated {
			err = sc.UpdateToken(ctx)
			if err != nil {
				return nil, 0, err
			}
			response, err = sc.client.ListCommentsStats(ctx, &stats.ListStatsRequest{
				PageSize:   int32(pageSize),
				PageNumber: int32(pageNumber),
				Token:      sc.token,
			})
			if err != nil {
				return nil, 0, err
			}
		} else {
			return nil, 0, err
		}
	}
	stat_s := make([]Stats, len(response.Stats))
	for idx, singleStat := range response.Stats {
		stat, err := convertSingleStat(singleStat)
		if err != nil {
			return nil, 0, err
		}
		stat_s[idx] = *stat
	}
	return stat_s, int(response.PageCount), nil
}

func (s *Server) listAccountsStats() http.HandlerFunc {
	type Response struct {
		Stats      []Stats `json:"stats"`
		PageSize   int     `json:"page_size"`
		PageNumber int     `json:"page_number"`
		PageCount  int     `json:"page_count"`
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
		if !user.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
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

		stat_s, pageCount, err := s.Stats.ListAccountsStats(ctx, page, size)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		response := Response{
			Stats:      stat_s,
			PageSize:   size,
			PageNumber: page,
			PageCount:  pageCount,
		}
		json, err := json.Marshal(response)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) listNewsStats() http.HandlerFunc {
	type Response struct {
		Stats      []Stats `json:"stats"`
		PageSize   int     `json:"page_size"`
		PageNumber int     `json:"page_number"`
		PageCount  int     `json:"page_count"`
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
		if !user.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
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

		stat_s, pageCount, err := s.Stats.ListNewsStats(ctx, page, size)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		response := Response{
			Stats:      stat_s,
			PageSize:   size,
			PageNumber: page,
			PageCount:  pageCount,
		}
		json, err := json.Marshal(response)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *Server) listCommentsStats() http.HandlerFunc {
	type Response struct {
		Stats      []Stats `json:"stats"`
		PageSize   int     `json:"page_size"`
		PageNumber int     `json:"page_number"`
		PageCount  int     `json:"page_count"`
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
		if !user.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
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

		stat_s, pageCount, err := s.Stats.ListCommentsStats(ctx, page, size)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		response := Response{
			Stats:      stat_s,
			PageSize:   size,
			PageNumber: page,
			PageCount:  pageCount,
		}
		json, err := json.Marshal(response)
		if err != nil {
			handleRPCErrors(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}
