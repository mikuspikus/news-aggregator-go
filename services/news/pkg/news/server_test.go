package news

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	simple_token_storage "github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/news/proto"
	"math"
	"net/url"
	"testing"
	"time"
)

type mockDb struct {
	News   []*News
	LastID int
}

func _safeSlice(news []*News, pageNumber, pageSize int32) []*News {
	if pageNumber < 0 || pageSize < 0 {
		return make([]*News, 0)
	}

	startIndex := pageNumber * pageSize
	if int(startIndex) > len(news) {
		return make([]*News, 0)
	}
	endIndex := int(startIndex + pageSize)

	if endIndex > len(news) {
		endIndex = len(news)
	}

	return news[startIndex:endIndex]
}

func (db *mockDb) push(useruid uuid.UUID, title string, uri url.URL) *News {
	news := new(News)

	news.UID = uuid.New()
	news.User = useruid
	news.Title = title
	news.URI = &uri
	news.Created = time.Now()
	news.Edited = time.Now()

	update := append(db.News, news)
	db.News = update
	db.LastID++

	return news
}

func (db *mockDb) find(uid uuid.UUID) (int, *News, error) {
	for idx, news := range db.News {
		if news.UID == uid {
			return idx, news, nil
		}
	}

	return -1, nil, errNotFound
}

func (db *mockDb) List(pageNumber, pageSize int32) ([]*News, int32, error) {
	news := _safeSlice(db.News, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(db.News)) / float64(pageSize))
	return news, int32(pageCount), nil
}

func (db *mockDb) Get(uid uuid.UUID) (*News, error) {
	_, news, err := db.find(uid)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func (db *mockDb) Create(user uuid.UUID, title string, uri url.URL) (*News, error) {
	news := db.push(user, title, uri)
	return news, nil
}

func (db *mockDb) Update(uid uuid.UUID, title string, uri url.URL) (*News, error) {
	_, news, err := db.find(uid)
	if err != nil {
		return nil, err
	}

	news.Title = title
	news.URI = &uri
	news.Edited = time.Now()

	return news, nil
}

func (db *mockDb) Delete(uid uuid.UUID) error {
	idx, _, err := db.find(uid)

	if err != nil {
		return err
	}

	news := db.News
	news[len(news)-1], news[idx] = news[idx], news[len(news)-1]
	return nil
}

func (db *mockDb) Close() {}

func setupDB() (*mockDb, error) {
	db := &mockDb{}

	url, err := url.Parse("http://www.google.com")
	if err != nil {
		return nil, err
	}
	for i := 0; i < 20; i++ {
		db.push(uuid.New(), fmt.Sprintf("test-title-%d", i), *url)
	}

	return db, nil
}

func setupServer(db *mockDb) *Service {
	server := &Service{
		db:           db,
		tokenStorage: &simple_token_storage.MockStorage{},
	}
	return server
}

func TestService_ListNews(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
	server := setupServer(db)
	request := &pb.ListNewsRequest{
		PageSize:   10,
		PageNumber: 1,
	}
	_, err = server.ListNews(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_GetNews(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
	server := setupServer(db)
	target := db.News[0]
	request := &pb.GetNewsRequest{Uid: target.UID.String()}
	_, err = server.GetNews(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AddNews(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
	server := setupServer(db)
	request := &pb.AddNewsRequest{
		Token:    "",
		UserUUID: uuid.New().String(),
		Title:    "add-news-title",
		Uri:      "https://google.com",
	}
	_, err = server.AddNews(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_EditNews(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
	server := setupServer(db)
	target := db.News[0]
	request := &pb.EditNewsRequest{
		Token: "",
		Uid:   target.UID.String(),
		Title: "edit-news-title",
		Uri:   "https://localhost.com",
	}
	_, err = server.EditNews(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_DeleteNews(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
	server := setupServer(db)
	target := db.News[0]
	request := &pb.DeleteNewsRequest{
		Token: "",
		Uid:   target.UID.String(),
	}
	_, err = server.DeleteNews(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}
