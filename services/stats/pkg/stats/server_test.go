package stats

import (
	"context"
	"github.com/google/uuid"
	simple_token_storage "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/stats/proto"
	"math"
	"testing"
	"time"
)

type mockTable struct {
	Stats  []*Stats
	LastID int
}

func (t *mockTable) push(useruid uuid.UUID, action string, input, output Attrs) *Stats {
	stats := new(Stats)

	stats.ID = int32(t.LastID)
	stats.User = useruid
	stats.Action = action
	stats.Input = input
	stats.Output = output
	stats.Timestamp = time.Now()

	update := append(t.Stats, stats)
	t.Stats = update
	t.LastID++

	return stats
}

func _safeSlice(stats []*Stats, pageNumber, pageSize int32) []*Stats {
	if pageNumber < 0 || pageSize < 0 {
		return make([]*Stats, 0)
	}

	startIndex := pageNumber * pageSize
	if int(startIndex) > len(stats) {
		return make([]*Stats, 0)
	}
	endIndex := int(startIndex + pageSize)

	if endIndex > len(stats) {
		endIndex = len(stats)
	}

	return stats[startIndex:endIndex]
}

type mockDB struct {
	Accounts *mockTable
	News     *mockTable
	Comments *mockTable
}

func (db *mockDB) ListAccounts(pageNumber, pageSize int32) ([]*Stats, int32, error) {
	stats := _safeSlice(db.Accounts.Stats, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(db.Accounts.Stats)) / float64(pageSize))
	return stats, int32(pageCount), nil
}

func (db *mockDB) AddAccounts(user uuid.UUID, action string, input, output Attrs) (*Stats, error) {
	stats := db.Accounts.push(user, action, input, output)
	return stats, nil
}

func (db *mockDB) ListNews(pageNumber, pageSize int32) ([]*Stats, int32, error) {
	stats := _safeSlice(db.News.Stats, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(db.News.Stats)) / float64(pageSize))
	return stats, int32(pageCount), nil
}

func (db *mockDB) AddNews(user uuid.UUID, action string, input, output Attrs) (*Stats, error) {
	stats := db.News.push(user, action, input, output)
	return stats, nil
}

func (db *mockDB) ListComments(pageNumber, pageSize int32) ([]*Stats, int32, error) {
	stats := _safeSlice(db.Comments.Stats, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(db.Comments.Stats)) / float64(pageSize))
	return stats, int32(pageCount), nil
}

func (db *mockDB) AddComments(user uuid.UUID, action string, input, output Attrs) (*Stats, error) {
	stats := db.Comments.push(user, action, input, output)
	return stats, nil
}

func (db *mockDB) Close() {}

func setupTable() *mockTable {
	table := &mockTable{}
	for i := 0; i < 35; i++ {
		table.push(uuid.New(), "some-test-action", nil, nil)
	}
	return table
}

func setupDB() *mockDB {
	db := &mockDB{
		Accounts: setupTable(),
		News:     setupTable(),
		Comments: setupTable(),
	}

	return db
}

func setupServer(db *mockDB) *Service {
	server := &Service{
		db:           db,
		tokenStorage: &simple_token_storage.MockStorage{},
	}
	return server
}

func TestService_ListAccountsStats(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.ListStatsRequest{
		PageSize:   10,
		PageNumber: 1,
		Token:      "",
	}
	_, err := server.ListAccountsStats(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AddAccountsStats(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.AddStatsRequest{
		UserUID: uuid.New().String(),
		Action:  "some-test-action",
		Input:   nil,
		Output:  nil,
		Token:   "",
	}
	_, err := server.AddAccountsStats(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_ListNewsStats(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.ListStatsRequest{
		PageSize:   10,
		PageNumber: 1,
		Token:      "",
	}
	_, err := server.ListNewsStats(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AddNewsStats(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.AddStatsRequest{
		UserUID: uuid.New().String(),
		Action:  "some-test-action",
		Input:   nil,
		Output:  nil,
		Token:   "",
	}
	_, err := server.AddNewsStats(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_ListCommentsStats(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.ListStatsRequest{
		PageSize:   10,
		PageNumber: 1,
		Token:      "",
	}
	_, err := server.ListCommentsStats(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AddCommentsStats(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.AddStatsRequest{
		UserUID: uuid.New().String(),
		Action:  "some-test-action",
		Input:   nil,
		Output:  nil,
		Token:   "",
	}
	_, err := server.AddCommentsStats(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}
