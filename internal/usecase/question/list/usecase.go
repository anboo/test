package list

import (
	"context"
	"fmt"

	entQ "test-question/internal/entity/question"
)

//go:generate mockery --name=questionRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported

type (
	questionRepository interface {
		List(ctx context.Context) ([]*entQ.Question, error)
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	repo   questionRepository
	logger logger
}

func NewUseCase(repo questionRepository, logger logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (uc *UseCase) ListQuestions(ctx context.Context) ([]*entQ.Question, error) {
	out, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list questions: %w", err)
	}

	uc.logger.DebugContext(ctx, "questions listed", "count", len(out))
	return out, nil
}
