package mocks

import (
	"github.com/stretchr/testify/mock"

	"example.com/internal/infrastructure/logger"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(append([]interface{}{msg}, args...)...)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(append([]interface{}{msg}, args...)...)
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(append([]interface{}{msg}, args...)...)
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(append([]interface{}{msg}, args...)...)
}

func (m *MockLogger) With(args ...interface{}) logger.Logger {
	return m
}