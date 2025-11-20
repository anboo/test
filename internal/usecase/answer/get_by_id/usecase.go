package get_by_id

import (
	"context"
	"fmt"

	entA "test-question/internal/entity/answer"

	"github.com/pkg/errors"
)

//go:generate mockery --name=answerRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported

type (
	answerRepository interface {
		GetByID(ctx context.Context, id int) (*entA.Answer, error)
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	repo   answerRepository
	logger logger
}

func NewUseCase(
	repo answerRepository,
	logger logger,
) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

func (uc *UseCase) GetAnswer(
	ctx context.Context,
	answerID int,
) (*entA.Answer, error) {
	a, err := uc.repo.GetByID(ctx, answerID)
	if err != nil {
		if errors.Is(err, entA.ErrAnswerNotFound) {
			return nil, entA.ErrAnswerNotFound
		}
		return nil, fmt.Errorf("get answer: %w", err)
	}

	uc.logger.DebugContext(ctx, "answer loaded",
		"answer_id", answerID,
		"user_id", a.UserID,
	)

	return a, nil
}
