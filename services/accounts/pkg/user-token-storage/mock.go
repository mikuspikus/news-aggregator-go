package user_token_storage

import "github.com/google/uuid"

type MockStorage struct {
}

func (s *MockStorage) Add(uid uuid.UUID) (string, error) {
	return "mock-token", nil
}

func (s *MockStorage) Get(token string) (string, error) {
	return "mock-user-uuid", nil
}

func (s *MockStorage) Delete(token string) error {
	return nil
}
