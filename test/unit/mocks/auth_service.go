package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	authapi "example.com/gen/openapi/auth/go"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignUp(ctx context.Context, req authapi.SignupRequest) (*authapi.SignupResponse, error) {
	args := m.Called(ctx, req)
	if response := args.Get(0); response != nil {
		return response.(*authapi.SignupResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, req authapi.LoginRequest) (*authapi.LoginResponse, error) {
	args := m.Called(ctx, req)
	if response := args.Get(0); response != nil {
		return response.(*authapi.LoginResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) LookupUser(ctx context.Context, email string) (*authapi.UserLookupResponse, error) {
	args := m.Called(ctx, email)
	if response := args.Get(0); response != nil {
		return response.(*authapi.UserLookupResponse), args.Error(1)
	}
	return nil, args.Error(1)
}