package tests

import (
	"github.com/stretchr/testify/mock"
	"time"
)

// Mock implementations
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheRepository) Set(key string, value []byte, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}