package simple_token_storage

type MockStorage struct{}

func (s *MockStorage) AddToken(appID, appSECRET string) (string, error) {
	return "mock-token", nil
}

func (s *MockStorage) CheckToken(token string) (bool, error) {
	return true, nil
}
