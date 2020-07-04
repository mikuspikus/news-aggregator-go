package comments

import (
	"math"
	"time"

	// "database/sql"
	"github.com/google/uuid"
)

// Datastore is database mock
type Datastore struct {
	Comments []*Comment
	LastID   int32
}

func (store *Datastore) FilterByNewsUUID(newsUUID uuid.UUID) []*Comment {
	comments := make([]*Comment, 0)

	for _, comment := range store.Comments {
		if comment.News == newsUUID {
			comments = append(comments, comment)
		}
	}

	return comments
}

func SafeSlice(comments []*Comment, pageNumber, pageSize int32) []*Comment {
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

// List returns several comments with pageNumber and pageSize params
func (store *Datastore) List(pageNumber, pageSize int32, newsUUID uuid.UUID) ([]*Comment, int32, error) {
	comments := make([]*Comment, 0)

	if newsUUID == uuid.Nil {
		comments = store.Comments
	} else {
		comments = store.FilterByNewsUUID(newsUUID)
	}

	comments = SafeSlice(comments, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(store.Comments)) / float64(pageSize))
	return comments, int32(pageCount), nil
}

// Get returns comment by ID
func (store *Datastore) Get(id int32) (*Comment, error) {
	_, comment, err := store.FindComment(id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// Delete deletes comment by ID
func (store *Datastore) Delete(id int32) error {
	idx, _, err := store.FindComment(id)

	if err != nil {
		return err
	}

	comments := store.Comments
	comments[len(comments)-1], comments[idx] = comments[idx], comments[len(comments)-1]
	store.Comments = comments[:len(comments)-1]

	return nil
}

// FindComment returns index and comment by id
func (store *Datastore) FindComment(id int32) (int, *Comment, error) {
	for idx, comment := range store.Comments {
		if comment.ID == id {
			return idx, comment, nil
		}
	}

	return -1, nil, errNotFound
}

// Update returns changed comment
func (store *Datastore) Update(id int32, body string) (*Comment, error) {
	_, comment, err := store.FindComment(id)

	if err != nil {
		return nil, err
	}

	now := time.Now()

	comment.Body = body
	comment.Edited = now

	return comment, nil
}

// Create adds new comment
func (store *Datastore) Create(user uuid.UUID, news uuid.UUID, body string) (*Comment, error) {
	comment := new(Comment)

	now := time.Now()

	comment.ID = store.LastID
	comment.User = user
	comment.News = news
	comment.Body = body
	comment.Created = now
	comment.Edited = now

	update := append(store.Comments, comment)
	store.Comments = update
	store.LastID++

	return comment, nil
}
