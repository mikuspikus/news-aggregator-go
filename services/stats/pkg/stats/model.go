package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"math"
	"time"
)

var (
	errNotCreated = errors.New("stats not created")
)

type Stats struct {
	ID        int32
	User      uuid.UUID
	Action    string
	Timestamp time.Time
	Input     Attrs
	Output    Attrs
}

type Attrs map[string]interface{}

func newDB(connstring string) (*db, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connstring)
	return &db{dbpool}, err
}

type DataStoreHandler interface {
	ListAccounts(pageNumber, pageSize int32) (stats []*Stats, pageCount int32, err error)
	AddAccounts(user uuid.UUID, action string, input, output Attrs) (stats *Stats, err error)

	ListNews(pageNumber, pageSize int32) (stats []*Stats, pageCount int32, err error)
	AddNews(user uuid.UUID, action string, input, output Attrs) (stats *Stats, err error)

	ListComments(pageNumber, pageSize int32) (stats []*Stats, pageCount int32, err error)
	AddComments(user uuid.UUID, action string, input, output Attrs) (stats *Stats, err error)

	Close()
}

type db struct {
	*pgxpool.Pool
}

func (db *db) pageCount(viewName string, pageSize int32) (pageCount int32, err error) {
	queryTemplate := "select count(*) from %s"
	query := fmt.Sprintf(queryTemplate, viewName)
	var rowsCount int32
	err = db.QueryRow(context.Background(), query).Scan(&rowsCount)
	if err != nil {
		return 0, err
	}
	pageCount = int32(math.Ceil(float64(rowsCount) / float64(pageSize)))
	return pageCount, nil
}

func (db *db) listStats(viewName string, pageNumber, pageSize int32) ([]*Stats, int32, error) {
	queryTemplate := "select  id, user_uid, action, timestamp, input, output from %s limit $1 offset $2"
	query := fmt.Sprintf(queryTemplate, viewName)
	lastRecord := pageNumber * pageSize

	rows, err := db.Query(context.Background(), query, pageSize, lastRecord)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	stats_s := make([]*Stats, 0)
	for rows.Next() {
		stats := new(Stats)

		var userUID string
		err := rows.Scan(&stats.ID, &userUID, &stats.Action, &stats.Timestamp, &stats.Input, &stats.Output)
		if err != nil {
			return nil, 0, err
		}

		stats.User, err = uuid.Parse(userUID)
		if err != nil {
			return nil, 0, err
		}

		stats_s = append(stats_s, stats)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	pageCount, err := db.pageCount(viewName, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return stats_s, pageCount, err
}

func (db *db) addStats(tableName string, user uuid.UUID, action string, input, output Attrs) (*Stats, error) {
	queryTemplate := "insert into %s (user_uid, action, timestamp, input, output) " +
		"values ($1, $2, $3, $4, $5) returning id"

	query := fmt.Sprintf(queryTemplate, tableName)
	now := time.Now()
	stats := new(Stats)
	var id int32
	binput, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	boutput, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow(context.Background(), query, user.String(), action, now, binput, boutput).Scan(&id)
	if err != nil {
		return nil, err
	}

	stats.ID = id
	stats.User = user
	stats.Action = action
	stats.Timestamp = now
	stats.Input = input
	stats.Output = output

	return stats, nil
}

func (db *db) ListAccounts(pageNumber, pageSize int32) (stats []*Stats, pageCount int32, err error) {
	return db.listStats("accounts_view", pageNumber, pageSize)
}

func (db *db) AddAccounts(user uuid.UUID, action string, input, output Attrs) (stats *Stats, err error) {
	return db.addStats("accounts", user, action, input, output)
}

func (db *db) ListNews(pageNumber, pageSize int32) (stats []*Stats, pageCount int32, err error) {
	return db.listStats("news_view", pageNumber, pageSize)
}

func (db *db) AddNews(user uuid.UUID, action string, input, output Attrs) (stats *Stats, err error) {
	return db.addStats("news", user, action, input, output)
}

func (db *db) ListComments(pageNumber, pageSize int32) (stats []*Stats, pageCount int32, err error) {
	return db.listStats("comments_view", pageNumber, pageSize)
}

func (db *db) AddComments(user uuid.UUID, action string, input, output Attrs) (stats *Stats, err error) {
	return db.addStats("comments", user, action, input, output)
}
