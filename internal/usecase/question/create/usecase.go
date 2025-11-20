package create

import (
	"context"
	"fmt"
	"time"

	entQ "test-question/internal/entity/question"
)

//go:generate mockery --name=questionRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=timer --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported

type (
	questionRepository interface {
		Create(ctx context.Context, q *entQ.Question) (*entQ.Question, error)
	}

	timer interface {
		Now() time.Time
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	repo   questionRepository
	timer  timer
	logger logger
}

func NewUseCase(
	questions questionRepository,
	timer timer,
	logger logger,
) *UseCase {
	return &UseCase{
		repo:   questions,
		timer:  timer,
		logger: logger,
	}
}

func (uc *UseCase) CreateQuestion(
	ctx context.Context,
	userID string,
	text string,
) (*entQ.Question, error) {
	q := &entQ.Question{
		Text:      text,
		UserID:    userID,
		CreatedAt: uc.timer.Now(),
	}

	out, err := uc.repo.Create(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("create question: %w", err)
	}

	uc.logger.DebugContext(ctx, "question created",
		"question_id", out.ID,
	)

	return out, nil
}
