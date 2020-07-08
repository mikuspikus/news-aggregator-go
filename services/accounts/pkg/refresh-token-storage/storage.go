package refresh_token_storage

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"
)

var (
	ErrNotFound   = errors.New("refresh token not found")
	ErrWrongToken = errors.New("wrong token")
)

const DefaultExpirationTime = time.Hour * 24 * 7 * 2

type RefreshTokenStorage struct {
	redis *redis.Client
}

func generateToken() string {
	return uuid.New().String()
}

func New(addr, password string, db int) (*RefreshTokenStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RefreshTokenStorage{redis: client}, nil
}

func (storage *RefreshTokenStorage) Add(token string) (string, error) {
	refreshToken := generateToken()

	_, err := storage.redis.Set(context.Background(), refreshToken, token, DefaultExpirationTime).Result()
	if err != nil {
		return "", nil
	}
	return refreshToken, nil
}

func (storage *RefreshTokenStorage) Get(refreshToken string) (string, error) {
	token, err := storage.redis.Get(context.Background(), refreshToken).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return token, nil
}

func (storage *RefreshTokenStorage) Check(token, refreshToken string) (bool, error) {
	redisToken, err := storage.Get(refreshToken)
	if err == redis.Nil {
		return false, ErrNotFound
	} else if err != nil {
		return false, err
	}
	return redisToken == token, nil
}

func (storage *RefreshTokenStorage) Del(refreshToken string) error {
	_, err := storage.redis.Del(context.Background(), refreshToken).Result()
	return err
}
