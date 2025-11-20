package auth

import (
	"context"
	"fmt"

	ent "test-question/internal/entity/user"
	repo "test-question/internal/repository/user"

	"github.com/pkg/errors"
)

//go:generate mockery --name=repository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported

type (
	repository interface {
		GetUserByUsernamePassword(ctx context.Context, username, password string) (*ent.User, error)
	}

	logger interface {
		InfoContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	rep    repository
	logger logger
}

func NewUseCase(rep repository, logger logger) *UseCase {
	return &UseCase{rep: rep, logger: logger}
}

func (uc *UseCase) AuthorizeUser(ctx context.Context, username, password string) (*ent.User, error) {
	user, err := uc.rep.GetUserByUsernamePassword(ctx, username, password)
	switch {
	case errors.Is(err, repo.ErrUserNotFound):
		uc.logger.InfoContext(ctx, "fail attempt login with username", username)
		return nil, ent.ErrUsernameOrPasswordIncorrect
	case err != nil:
		return nil, fmt.Errorf("get user by username and password: %w", err)
	}
	return user, nil
}
