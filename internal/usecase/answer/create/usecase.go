package create

import (
	"context"
	"fmt"
	"time"

	entA "test-question/internal/entity/answer"
	entQ "test-question/internal/entity/question"

	"github.com/pkg/errors"
)

//go:generate mockery --name=answerRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=questionRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=timer --output=mocks --outpkg=mocks --exported

type (
	answerRepository interface {
		Create(ctx context.Context, a *entA.Answer) (*entA.Answer, error)
	}

	questionRepository interface {
		GetByID(ctx context.Context, id int) (*entQ.Question, error)
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}

	timer interface {
		Now() time.Time
	}
)

type UseCase struct {
	repo      answerRepository
	questions questionRepository
	timer     timer
	logger    logger
}

func NewUseCase(
	answers answerRepository,
	questions questionRepository,
	timer timer,
	logger logger,
) *UseCase {
	return &UseCase{
		repo:      answers,
		questions: questions,
		timer:     timer,
		logger:    logger,
	}
}

func (uc *UseCase) CreateAnswer(
	ctx context.Context,
	questionID int,
	userID string,
	text string,
) (*entA.Answer, error) {
	_, err := uc.questions.GetByID(ctx, questionID)
	if err != nil {
		if errors.Is(err, entQ.ErrQuestionNotFound) {
			return nil, entA.ErrRequestedQuestionNotFound
		}
		return nil, fmt.Errorf("check question exists: %w", err)
	}

	a := &entA.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
		CreatedAt:  uc.timer.Now(),
	}

	out, err := uc.repo.Create(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("create answer: %w", err)
	}

	uc.logger.DebugContext(ctx, "answer created",
		"answer_id", out.ID,
		"question_id", questionID,
		"user_id", userID,
	)

	return out, nil
}
