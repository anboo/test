package auth_test

import (
	"context"
	"errors"
	"testing"

	ent "test-question/internal/entity/user"
	repo "test-question/internal/repository/user"
	"test-question/internal/usecase/auth"
	"test-question/internal/usecase/auth/mocks"

	"github.com/stretchr/testify/require"
)

func TestAuthorizeUser_Success(t *testing.T) {
	ctx := context.Background()

	r := mocks.NewRepository(t)
	l := mocks.NewLogger(t)

	r.On(
		"GetUserByUsernamePassword",
		ctx,
		"john",
		"pass123",
	).Return(&ent.User{Username: "john"}, nil)

	uc := auth.NewUseCase(r, l)

	u, err := uc.AuthorizeUser(ctx, "john", "pass123")
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, "john", u.Username)
}

func TestAuthorizeUser_IncorrectCredentials(t *testing.T) {
	ctx := context.Background()

	r := mocks.NewRepository(t)
	l := mocks.NewLogger(t)

	r.On(
		"GetUserByUsernamePassword",
		ctx,
		"john",
		"wrong",
	).Return(nil, repo.ErrUserNotFound)

	l.On(
		"InfoContext",
		ctx,
		"fail attempt login with username",
		"john",
	).Return()

	uc := auth.NewUseCase(r, l)

	u, err := uc.AuthorizeUser(ctx, "john", "wrong")
	require.Error(t, err)
	require.ErrorIs(t, err, ent.ErrUsernameOrPasswordIncorrect)
	require.Nil(t, u)
}

func TestAuthorizeUser_InternalError(t *testing.T) {
	ctx := context.Background()

	r := mocks.NewRepository(t)
	l := mocks.NewLogger(t)

	r.On(
		"GetUserByUsernamePassword",
		ctx,
		"john",
		"pass123",
	).Return(nil, errors.New("db down"))

	uc := auth.NewUseCase(r, l)

	u, err := uc.AuthorizeUser(ctx, "john", "pass123")
	require.Nil(t, u)
	require.Error(t, err)
	require.Contains(t, err.Error(), "get user by username and password")
}
