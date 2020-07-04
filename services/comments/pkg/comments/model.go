package comments

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	errNotCreated = errors.New("comment not created")
	errNotFound   = errors.New("comment not found")
)

// Comment is inner representation of comment
type Comment struct {
	ID      int32
	User    uuid.UUID
	News    uuid.UUID
	Body    string
	Created time.Time
	Edited  time.Time
}

// DataStoreHandler is interface for data storage manipulation
type DataStoreHandler interface {
	List(int32, int32, uuid.UUID) ([]*Comment, int32, error)
	Get(int32) (*Comment, error)
	Create(uuid.UUID, uuid.UUID, string) (*Comment, error)
	Update(int32, string) (*Comment, error)
	Delete(int32) error
}

type db struct {
	*pgxpool.Pool
}

func newDB(connString string) (*db, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connString)
	return &db{dbpool}, err
}

func (db *db) pageCount(pageSize int32, newsUUID uuid.UUID) (int32, error) {
	query :=
		"select count(*) " +
			"from comments_view " +
			"where news_uid=$1"

	row := db.QueryRow(context.Background(), query, newsUUID)
	var rowsCount int32
	err := row.Scan(&rowsCount)
	if err != nil {
		return 0, err
	}
	pageCount := math.Ceil(float64(rowsCount) / float64(pageSize))

	return int32(pageCount), nil
}

func (db *db) List(pageNumber, pageSize int32, newsUUID uuid.UUID) ([]*Comment, int32, error) {
	query := "select * from comments_view " +
		"where news_uid=$1 " +
		"limit $2 offset $3"
	lastRecord := pageNumber * pageSize

	rows, err := db.Query(context.Background(), query, newsUUID.String(), pageSize, lastRecord)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	comments := make([]*Comment, 0)
	for rows.Next() {
		comment := new(Comment)
		var id int32
		var userUUID, newsUUID string

		err := rows.Scan(&id, &userUUID, &newsUUID, &comment.Body, &comment.Created, &comment.Edited)
		if err != nil {
			return nil, 0, err
		}
		comment.ID = id

		comment.User, err = uuid.Parse(userUUID)
		if err != nil {
			return nil, 0, err
		}

		comment.News, err = uuid.Parse(newsUUID)
		if err != nil {
			return nil, 0, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	pageCount, err := db.pageCount(pageSize, newsUUID)
	if err != nil {
		return nil, 0, err
	}

	return comments, pageCount, nil
}

func (db *db) Get(id int32) (*Comment, error) {
	query := "select * from comments_view where id=$1"

	row := db.QueryRow(context.Background(), query, id)

	comment := new(Comment)
	var userUUID, newsUUID string

	err := row.Scan(&id, &userUUID, &newsUUID, &comment.Body, &comment.Created, &comment.Edited)
	if err != nil {
		return nil, err
	}
	comment.ID = id
	comment.User, err = uuid.Parse(userUUID)
	if err != nil {
		return nil, err
	}
	comment.News, err = uuid.Parse(newsUUID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (db *db) Create(userUUID, newsUUID uuid.UUID, body string) (*Comment, error) {
	query :=
		"insert into comments (user_uid, news_uid, body, created_at, edited_at) " +
			"values ($1, $2, $3, $4, $5) " +
			"returning id"

	now := time.Now()
	var id int32

	err := db.QueryRow(context.Background(), query, userUUID.String(), newsUUID.String(), body, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	comment := new(Comment)
	comment.ID = id
	comment.User = userUUID
	comment.News = newsUUID
	comment.Body = body
	comment.Created = now
	comment.Edited = now

	return comment, err
}

func (db *db) Update(id int32, body string) (*Comment, error) {
	query := "update comments set body=$1, edited_at=$2 " +
		"where id=$3 " +
		"returning user_uid, news_uid, created_at"
	now := time.Now()
	var userUUID, newsUUID string

	comment := new(Comment)
	comment.ID = id
	comment.Edited = now
	comment.Body = body

	err := db.QueryRow(context.Background(), query, body, now, id).Scan(&userUUID, &newsUUID, &comment.Created)
	if err != nil {
		return nil, errNotFound
	}

	comment.User, err = uuid.Parse(userUUID)
	if err != nil {
		return nil, err
	}

	comment.News, err = uuid.Parse(newsUUID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (db *db) Delete(id int32) error {
	query := "delete from comments where id=$1"
	cmd, err := db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsCount := cmd.RowsAffected()
	if rowsCount == 0 {
		return errNotFound
	}

	return nil
}
