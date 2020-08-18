package refresh_token_storage

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

const DefaultExpirationTime = time.Hour * 24 * 7 * 2

type RefreshTokenStorage interface {
	Add(token string) (string, error)
	Get(refreshToken string) (string, error)
	Check(token, refreshToken string) (bool, error)
	Delete(refreshToken string) error
}

type RTStorage struct {
	redis *redis.Client
}

func generateToken() string {
	return uuid.New().String()
}

func New(addr, password string, db int) (*RTStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RTStorage{redis: client}, nil
}

func (storage *RTStorage) Add(token string) (string, error) {
	refreshToken := generateToken()

	_, err := storage.redis.Set(context.Background(), refreshToken, token, DefaultExpirationTime).Result()
	if err != nil {
		return "", nil
	}
	return refreshToken, nil
}

func (storage *RTStorage) Get(refreshToken string) (string, error) {
	token, err := storage.redis.Get(context.Background(), refreshToken).Result()
	switch err {
	case nil:
		return token, nil
	case redis.Nil:
		return "", ErrTokenNotFound
	default:
		return "", err
	}
}

func (storage *RTStorage) Check(token, refreshToken string) (bool, error) {
	refreshTokenFromRedis, err := storage.Get(refreshToken)
	switch err {
	case nil:
		return refreshTokenFromRedis == token, nil
	case redis.Nil:
		return false, ErrTokenNotFound
	default:
		return false, err
	}
}

func (storage *RTStorage) Delete(refreshToken string) error {
	_, err := storage.redis.Del(context.Background(), refreshToken).Result()
	return err
}
