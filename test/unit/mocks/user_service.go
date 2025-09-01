package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	authapi "example.com/gen/openapi/auth/go"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) LookupUser(ctx context.Context, email string) (*authapi.UserLookupResponse, error) {
	args := m.Called(ctx, email)
	if response := args.Get(0); response != nil {
		return response.(*authapi.UserLookupResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
