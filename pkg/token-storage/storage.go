package token_storage

import (
	"context"
	"errors"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("app with ID not found")
	ErrWrongSecret = errors.New("invalid credentials")
)

const DefaultExpirationTime = time.Minute * 30

func generateToken() string {
	return uuid.New().String()
}

type APITokenStorage struct {
	redis *redis.Client
	apps  map[string]string
}

func New(addr, password string, db int, apps map[string]string) (*APITokenStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &APITokenStorage{redis: client, apps: apps}, nil
}

func (storage *APITokenStorage) AddToken(appID, appSECRET string) (string, error) {
	secret, exists := storage.apps[appID]

	if !exists {
		return "", ErrNotFound
	}
	if secret != appSECRET {
		return "", ErrWrongSecret
	}

	token := generateToken()
	_, err := storage.redis.Set(context.Background(), token, true, DefaultExpirationTime).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (storage *APITokenStorage) CheckToken(token string) (bool, error) {
	exists, err := storage.redis.Get(context.Background(), token).Result()
	return exists == "1", err
}
