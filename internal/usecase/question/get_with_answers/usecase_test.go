package get_with_answers

import (
	"context"
	"errors"
	"testing"
	"time"

	entA "test-question/internal/entity/answer"
	entQ "test-question/internal/entity/question"
	mocks2 "test-question/internal/usecase/question/get_with_answers/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetQuestionWithAnswers_Success(t *testing.T) {
	ctx := context.Background()

	mQ := mocks2.NewQuestionRepository(t)
	mA := mocks2.NewAnswerRepository(t)
	mL := mocks2.NewLogger(t)

	now := time.Now()

	mQ.
		On("GetByID", ctx, 10).
		Return(&entQ.Question{
			ID:        10,
			Text:      "hello",
			UserID:    "u1",
			CreatedAt: now,
		}, nil)

	mA.
		On("ListByQuestionID", ctx, 10).
		Return([]*entA.Answer{
			{ID: 1, QuestionID: 10, UserID: "u1", Text: "ok"},
			{ID: 2, QuestionID: 10, UserID: "u2", Text: "yo"},
		}, nil)

	mL.
		On("DebugContext",
			mock.MatchedBy(func(_ context.Context) bool { return true }),
			"loaded question with answers",
			"question_id", 10,
			"answers", 2,
		).
		Return()

	ucase := NewUseCase(mQ, mA, mL)

	out, err := ucase.GetQuestionWithAnswers(ctx, 10)
	require.NoError(t, err)

	require.Equal(t, 10, out.Question.ID)
	require.Equal(t, "hello", out.Question.Text)
	require.Len(t, out.Answers, 2)
	require.Equal(t, 1, out.Answers[0].ID)
}

func TestGetQuestionWithAnswers_QuestionNotFound(t *testing.T) {
	ctx := context.Background()

	mQ := mocks2.NewQuestionRepository(t)
	mA := mocks2.NewAnswerRepository(t)
	mL := mocks2.NewLogger(t)

	mQ.
		On("GetByID", ctx, 99).
		Return(nil, entQ.ErrQuestionNotFound)

	ucase := NewUseCase(mQ, mA, mL)

	out, err := ucase.GetQuestionWithAnswers(ctx, 99)
	require.Nil(t, out)
	require.ErrorIs(t, err, entQ.ErrQuestionNotFound)
}

func TestGetQuestionWithAnswers_GetQuestionError(t *testing.T) {
	ctx := context.Background()

	mQ := mocks2.NewQuestionRepository(t)
	mA := mocks2.NewAnswerRepository(t)
	mL := mocks2.NewLogger(t)

	mQ.
		On("GetByID", ctx, 10).
		Return(nil, errors.New("db down"))

	ucase := NewUseCase(mQ, mA, mL)

	out, err := ucase.GetQuestionWithAnswers(ctx, 10)
	require.Nil(t, out)
	require.Contains(t, err.Error(), "get question")
}

func TestGetQuestionWithAnswers_ListAnswersError(t *testing.T) {
	ctx := context.Background()

	mQ := mocks2.NewQuestionRepository(t)
	mA := mocks2.NewAnswerRepository(t)
	mL := mocks2.NewLogger(t)

	mQ.
		On("GetByID", ctx, 10).
		Return(&entQ.Question{ID: 10}, nil)

	mA.
		On("ListByQuestionID", ctx, 10).
		Return(nil, errors.New("answers fail"))

	ucase := NewUseCase(mQ, mA, mL)

	out, err := ucase.GetQuestionWithAnswers(ctx, 10)
	require.Nil(t, out)
	require.Contains(t, err.Error(), "list answers")
}
