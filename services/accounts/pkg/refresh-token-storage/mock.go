package refresh_token_storage

type MockStorage struct {
}

func (s *MockStorage) Add(token string) (string, error) {
	return "mock-refresh-token", nil
}

func (s *MockStorage) Get(refreshToken string) (string, error) {
	return "mock-token", nil
}

func (s *MockStorage) Check(token, refreshToken string) (bool, error) {
	return true, nil
}

func (s *MockStorage) Delete(refreshToken string) error {
	return nil
}
