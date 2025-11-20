package delete_answer

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
		Delete(ctx context.Context, id int) error
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	answerRepo answerRepository
	logger     logger
}

func NewUseCase(
	answerRepo answerRepository,
	logger logger,
) *UseCase {
	return &UseCase{
		answerRepo: answerRepo,
		logger:     logger,
	}
}

func (uc *UseCase) DeleteAnswer(
	ctx context.Context,
	answerID int,
	userID string,
) error {
	a, err := uc.answerRepo.GetByID(ctx, answerID)
	if err != nil {
		if errors.Is(err, entA.ErrAnswerNotFound) {
			return err
		}
		return fmt.Errorf("get answer: %w", err)
	}

	if a.UserID != userID {
		return entA.ErrAccessDenied
	}

	if err = uc.answerRepo.Delete(ctx, answerID); err != nil {
		return fmt.Errorf("delete answer: %w", err)
	}

	uc.logger.DebugContext(ctx, "answer deleted",
		"answer_id", answerID,
		"user_id", userID,
	)

	return nil
}
