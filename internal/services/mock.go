package services

import (
	"os"
)

type MockService struct{}

func NewMockService() *MockService {
	return &MockService{}
}

func (m *MockService) OpenMockFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}
