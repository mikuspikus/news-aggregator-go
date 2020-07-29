package news

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"net/url"
	"time"
)

var (
	errNotCreated = errors.New("news not created")
	errNotFound   = errors.New("news not found")
)

type News struct {
	UID     uuid.UUID
	User    uuid.UUID
	Title   string
	URI     *url.URL
	Created time.Time
	Edited  time.Time
}

type DataStoreHandler interface {
	List(pageNumber, pageSize int32) ([]*News, int32, error)
	Get(uid uuid.UUID) (*News, error)
	Create(user uuid.UUID, title string, uri url.URL) (*News, error)
	Update(uid uuid.UUID, title string, uri url.URL) (*News, error)
	Delete(uid uuid.UUID) error
	Close()
}

type db struct {
	*pgxpool.Pool
}

func newDB(connstring string) (*db, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	return &db{dbpool}, err
}

func (db *db) pageCount(pageSize int32) (int32, error) {
	query := "select count(*) from news_view"
	var rowsCount int32
	err := db.QueryRow(context.Background(), query).Scan(&rowsCount)
	if err != nil {
		return 0, err
	}
	pageCount := math.Ceil(float64(rowsCount) / float64(pageSize))
	return int32(pageCount), nil
}

func (db *db) List(pageNumber, pageSize int32) ([]*News, int32, error) {
	query := "select * from news_view limit $1 offset $2"
	lastRecord := pageNumber * pageSize

	rows, err := db.Query(context.Background(), query, pageSize, lastRecord)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	newsS := make([]*News, 0)
	for rows.Next() {
		news := new(News)
		var uid, user, uri string

		err := rows.Scan(&uid, &user, &news.Title, &uri, &news.Created, &news.Edited)
		if err != nil {
			return nil, 0, err
		}

		news.UID, err = uuid.Parse(uid)
		if err != nil {
			return nil, 0, err
		}

		news.User, err = uuid.Parse(user)
		if err != nil {
			return nil, 0, err
		}

		news.URI, err = url.ParseRequestURI(uri)
		if err != nil {
			return nil, 0, err
		}

		newsS = append(newsS, news)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	pageCount, err := db.pageCount(pageSize)
	if err != nil {
		return nil, 0, err
	}
	return newsS, pageCount, nil
}

func (db *db) Get(uid uuid.UUID) (*News, error) {
	query := "select * from news_view where uid=$1"
	var uidMock, user, uri string
	news := new(News)

	err := db.QueryRow(context.Background(), query, uid.String()).Scan(&uidMock, &user, &news.Title, &uri, &news.Created, &news.Edited)
	if err != nil {
		return nil, err
	}
	news.UID = uid
	news.User, err = uuid.Parse(user)
	if err != nil {
		return nil, err
	}
	news.URI, err = url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (db *db) Create(user uuid.UUID, title string, uri url.URL) (*News, error) {
	query := "insert into news (uid, user_uid, title, uri, created_at, edited_at) values($1, $2, $3, $4, $5, $6)"
	now := time.Now()
	uid := uuid.New()

	cmd, err := db.Exec(context.Background(), query, uid, user.String(), title, uri.String(), now, now)
	if err != nil {
		return nil, err
	}
	if rowsCount := cmd.RowsAffected(); rowsCount == 0 {
		return nil, errNotCreated
	}

	news := new(News)
	news.UID = uid
	news.User = user
	news.Title = title
	news.URI = &uri
	news.Created = now
	news.Edited = now

	return news, nil
}

func (db *db) Update(uid uuid.UUID, title string, uri url.URL) (*News, error) {
	query := "update news set title=$1, uri=$2, edited_at=$3 where uid=$4 returning user_uid, created_at"

	now := time.Now()
	var user string

	news := new(News)
	news.UID = uid
	news.Title = title
	news.URI = &uri
	news.Edited = now

	err := db.QueryRow(context.Background(), query, title, uri.String(), now, uid.String()).Scan(&user, &news.Created)
	if err != nil {
		return nil, err
	}

	news.User, err = uuid.Parse(user)
	if err != nil {
		return nil, err
	}

	return news, nil
}

func (db *db) Delete(uid uuid.UUID) error {
	query := "delete from news where uid=$1"
	cmd, err := db.Exec(context.Background(), query, uid.String())
	if err != nil {
		return err
	}
	if rowsAffected := cmd.RowsAffected(); rowsAffected == 0 {
		return errNotFound
	}
	return nil
}
