package comments

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/comments/proto"
	"math"
	"testing"
	"time"
)

func _safeSlice(comments []*Comment, pageNumber, pageSize int32) []*Comment {
	if pageNumber < 0 || pageSize < 0 {
		return make([]*Comment, 0)
	}

	startIndex := pageNumber * pageSize
	if int(startIndex) > len(comments) {
		return make([]*Comment, 0)
	}
	endIndex := int(startIndex + pageSize)

	if endIndex > len(comments) {
		endIndex = len(comments)
	}

	return comments[startIndex:endIndex]
}

// mockDB is database mock
type mockDB struct {
	Comments []*Comment
	LastID   int32
}

func (db *mockDB) push(useruid uuid.UUID, newsuid uuid.UUID, body string) *Comment {
	comment := new(Comment)

	now := time.Now()

	comment.ID = db.LastID
	comment.User = useruid
	comment.News = newsuid
	comment.Body = body
	comment.Created = now
	comment.Edited = now

	update := append(db.Comments, comment)
	db.Comments = update
	db.LastID++

	return comment
}

func (db *mockDB) FilterByNewsUUID(newsUUID uuid.UUID) []*Comment {
	comments := make([]*Comment, 0)

	for _, comment := range db.Comments {
		if comment.News == newsUUID {
			comments = append(comments, comment)
		}
	}

	return comments
}

// List returns several comments with pageNumber and pageSize params
func (db *mockDB) List(pageNumber, pageSize int32, newsUUID uuid.UUID) ([]*Comment, int32, error) {
	comments := make([]*Comment, 0)

	if newsUUID == uuid.Nil {
		comments = db.Comments
	} else {
		comments = db.FilterByNewsUUID(newsUUID)
	}

	comments = _safeSlice(comments, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(db.Comments)) / float64(pageSize))
	return comments, int32(pageCount), nil
}

// Get returns comment by ID
func (db *mockDB) Get(id int32) (*Comment, error) {
	_, comment, err := db.FindComment(id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// Delete deletes comment by ID
func (db *mockDB) Delete(id int32) error {
	idx, _, err := db.FindComment(id)

	if err != nil {
		return err
	}

	comments := db.Comments
	comments[len(comments)-1], comments[idx] = comments[idx], comments[len(comments)-1]
	db.Comments = comments[:len(comments)-1]

	return nil
}

// FindComment returns index and comment by id
func (db *mockDB) FindComment(id int32) (int, *Comment, error) {
	for idx, comment := range db.Comments {
		if comment.ID == id {
			return idx, comment, nil
		}
	}

	return -1, nil, errNotFound
}

// Update returns changed comment
func (db *mockDB) Update(id int32, body string) (*Comment, error) {
	_, comment, err := db.FindComment(id)

	if err != nil {
		return nil, err
	}

	now := time.Now()

	comment.Body = body
	comment.Edited = now

	return comment, nil
}

// Create adds new comment
func (db *mockDB) Create(user uuid.UUID, news uuid.UUID, body string) (*Comment, error) {
	comment := db.push(user, news, body)
	return comment, nil
}

func (db *mockDB) Close() {}

func setupDB() *mockDB {
	db := &mockDB{}
	for i := 0; i < 20; i++ {
		db.push(uuid.New(), uuid.New(), fmt.Sprintf("comment-body-%d", i))
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

func TestService_ListComments(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	target := db.Comments[0]
	request := &pb.ListCommentsRequest{
		NewsUUID:   target.News.String(),
		PageSize:   10,
		PageNumber: 1,
	}

	_, err := server.ListComments(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_GetComment(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	target := db.Comments[0]
	request := &pb.GetCommentRequest{Id: target.ID}
	_, err := server.GetComment(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_DeleteComment(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	target := db.Comments[0]
	request := &pb.DeleteCommentRequest{
		Token: "",
		Id:    target.ID,
	}
	_, err := server.DeleteComment(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_EditComment(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	target := db.Comments[0]
	request := &pb.EditCommentRequest{
		Token: "",
		Id:    target.ID,
		Body:  "edited-comment-body",
	}
	_, err := server.EditComment(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AddComment(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.AddCommentRequest{
		Token:    "",
		NewsUUID: uuid.New().String(),
		UserUUID: uuid.New().String(),
		Body:     "add-comment-body",
	}
	_, err := server.AddComment(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}
