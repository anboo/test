package get_with_answers

import (
	"context"
	"fmt"

	entA "test-question/internal/entity/answer"
	entQ "test-question/internal/entity/question"

	"github.com/pkg/errors"
)

type QuestionWithAnswers struct {
	Question *entQ.Question `json:"question"`
	Answers  []*entA.Answer `json:"answers"`
}

//go:generate mockery --name=questionRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=answerRepository --output=mocks --outpkg=mocks --exported
//go:generate mockery --name=logger --output=mocks --outpkg=mocks --exported

type (
	questionRepository interface {
		GetByID(ctx context.Context, id int) (*entQ.Question, error)
	}

	answerRepository interface {
		ListByQuestionID(ctx context.Context, questionID int) ([]*entA.Answer, error)
	}

	logger interface {
		DebugContext(ctx context.Context, msg string, args ...any)
	}
)

type UseCase struct {
	questions questionRepository
	answers   answerRepository
	logger    logger
}

func NewUseCase(
	qRepo questionRepository,
	aRepo answerRepository,
	logger logger,
) *UseCase {
	return &UseCase{questions: qRepo, answers: aRepo, logger: logger}
}

func (uc *UseCase) GetQuestionWithAnswers(
	ctx context.Context,
	questionID int,
) (*QuestionWithAnswers, error) {

	q, err := uc.questions.GetByID(ctx, questionID)
	if err != nil {
		if errors.Is(err, entQ.ErrQuestionNotFound) {
			return nil, entQ.ErrQuestionNotFound
		}
		return nil, fmt.Errorf("get question: %w", err)
	}

	ans, err := uc.answers.ListByQuestionID(ctx, questionID)
	if err != nil {
		return nil, fmt.Errorf("list answers: %w", err)
	}

	uc.logger.
		DebugContext(ctx, "loaded question with answers",
			"question_id", questionID,
			"answers", len(ans),
		)

	return &QuestionWithAnswers{
		Question: q,
		Answers:  ans,
	}, nil
}
