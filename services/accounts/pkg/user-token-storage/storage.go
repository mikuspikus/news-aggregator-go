package user_token_storage

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

var (
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrWrongToken    = errors.New("wrong token")
)

const DefaultExpirationTime = time.Minute * 30

type UserTokenStorage interface {
	Add(uid uuid.UUID) (string, error)
	Get(token string) (string, error)
	Delete(token string) error
}

type UTStorage struct {
	redis *redis.Client
}

func generateToken() string {
	return uuid.New().String()
}

func New(addr, password string, db int) (*UTStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &UTStorage{redis: client}, nil
}

func (storage *UTStorage) Add(uid uuid.UUID) (string, error) {
	token := generateToken()

	_, err := storage.redis.Set(context.Background(), token, uid.String(), DefaultExpirationTime).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (storage *UTStorage) Get(token string) (string, error) {
	uid, err := storage.redis.Get(context.Background(), token).Result()
	switch err {
	case nil:
		return uid, nil
	case redis.Nil:
		return "", ErrTokenNotFound
	default:
		return "", err
	}
}

func (storage *UTStorage) Delete(token string) error {
	_, err := storage.redis.Del(context.Background(), token).Result()
	return err
}
