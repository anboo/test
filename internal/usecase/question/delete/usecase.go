package delete //nolint:predeclared

import (
	"context"
	"fmt"

	entQ "test-question/internal/entity/question"

	"github.com/pkg/errors"
)

//go:generate mockery --name=questionRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=answerRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=unitOfWork --output=mocks --outpkg=mocks --exported

type (
	questionRepository interface {
		GetByID(ctx context.Context, id int) (*entQ.Question, error)
		Delete(ctx context.Context, id int) error
	}

	answerRepository interface {
		DeleteByQuestionID(ctx context.Context, questionID int) error
	}

	unitOfWork interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	questionRepo questionRepository
	answerRepo   answerRepository
	uow          unitOfWork
	logger       logger
}

func NewUseCase(
	questionRepo questionRepository,
	answerRepo answerRepository,
	uow unitOfWork,
	logger logger,
) *UseCase {
	return &UseCase{
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		uow:          uow,
		logger:       logger,
	}
}

func (uc *UseCase) DeleteQuestion(
	ctx context.Context,
	questionID int,
	userID string,
) error {
	q, err := uc.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		if errors.Is(err, entQ.ErrQuestionNotFound) {
			return err
		}
		return fmt.Errorf("get question: %w", err)
	}

	if q.UserID != userID {
		return entQ.ErrAccessDenied
	}

	return uc.uow.Do(ctx, func(ctx context.Context) error {
		if err = uc.questionRepo.Delete(ctx, questionID); err != nil {
			return fmt.Errorf("delete question: %w", err)
		}

		if err = uc.answerRepo.DeleteByQuestionID(ctx, questionID); err != nil {
			return fmt.Errorf("delete answers: %w", err)
		}

		uc.logger.DebugContext(ctx, "question deleted with all answers",
			"question_id", questionID,
			"user_id", userID,
		)

		return nil
	})
}
