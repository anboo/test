package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	entA "test-question/internal/entity/answer"
	entQ "test-question/internal/entity/question"
	uc "test-question/internal/usecase/answer/create"
	"test-question/internal/usecase/answer/create/mocks"

	"github.com/stretchr/testify/require"
)

func TestCreateAnswer_Success(t *testing.T) {
	ctx := context.Background()

	now := time.Date(2024, 11, 20, 12, 0, 0, 0, time.UTC)

	mAnswers := mocks.NewAnswerRepository(t)
	mQuestions := mocks.NewQuestionRepository(t)
	mTimer := mocks.NewTimer(t)
	mLogger := mocks.NewLogger(t)

	mQuestions.
		On("GetByID", ctx, 10).
		Return(&entQ.Question{ID: 10}, nil)

	mTimer.
		On("Now").
		Return(now)

	expectedInput := &entA.Answer{
		QuestionID: 10,
		UserID:     "u1",
		Text:       "hello",
		CreatedAt:  now,
	}

	mAnswers.
		On("Create", ctx, expectedInput).
		Return(&entA.Answer{
			ID:         55,
			QuestionID: 10,
			UserID:     "u1",
			Text:       "hello",
			CreatedAt:  now,
		}, nil)

	mLogger.
		On("DebugContext",
			ctx,
			"answer created",
			"answer_id", 55,
			"question_id", 10,
			"user_id", "u1",
		).
		Return()

	ucase := uc.NewUseCase(mAnswers, mQuestions, mTimer, mLogger)

	out, err := ucase.CreateAnswer(ctx, 10, "u1", "hello")
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, 55, out.ID)
	require.Equal(t, 10, out.QuestionID)
	require.Equal(t, "u1", out.UserID)
	require.Equal(t, "hello", out.Text)
	require.Equal(t, now, out.CreatedAt)
}

func TestCreateAnswer_QuestionNotFound(t *testing.T) {
	ctx := context.Background()

	mAnswers := mocks.NewAnswerRepository(t)
	mQuestions := mocks.NewQuestionRepository(t)
	mTimer := mocks.NewTimer(t)
	mLogger := mocks.NewLogger(t)

	mQuestions.
		On("GetByID", ctx, 99).
		Return(nil, entQ.ErrQuestionNotFound)

	ucase := uc.NewUseCase(mAnswers, mQuestions, mTimer, mLogger)

	out, err := ucase.CreateAnswer(ctx, 99, "u1", "aaa")

	require.Nil(t, out)
	require.ErrorIs(t, err, entA.ErrRequestedQuestionNotFound)
}

func TestCreateAnswer_QuestionRepoError(t *testing.T) {
	ctx := context.Background()

	mAnswers := mocks.NewAnswerRepository(t)
	mQuestions := mocks.NewQuestionRepository(t)
	mTimer := mocks.NewTimer(t)
	mLogger := mocks.NewLogger(t)

	mQuestions.
		On("GetByID", ctx, 5).
		Return(nil, errors.New("db down"))

	ucase := uc.NewUseCase(mAnswers, mQuestions, mTimer, mLogger)

	out, err := ucase.CreateAnswer(ctx, 5, "u1", "aaa")

	require.Nil(t, out)
	require.Error(t, err)
	require.Contains(t, err.Error(), "check question exists")
}

func TestCreateAnswer_CreateError(t *testing.T) {
	ctx := context.Background()

	now := time.Date(2024, 11, 20, 12, 0, 0, 0, time.UTC)

	mAnswers := mocks.NewAnswerRepository(t)
	mQuestions := mocks.NewQuestionRepository(t)
	mTimer := mocks.NewTimer(t)
	mLogger := mocks.NewLogger(t)

	mQuestions.
		On("GetByID", ctx, 7).
		Return(&entQ.Question{ID: 7}, nil)

	mTimer.
		On("Now").
		Return(now)

	expectedInput := &entA.Answer{
		QuestionID: 7,
		UserID:     "u1",
		Text:       "xxx",
		CreatedAt:  now,
	}

	mAnswers.
		On("Create", ctx, expectedInput).
		Return(nil, errors.New("insert failed"))

	ucase := uc.NewUseCase(mAnswers, mQuestions, mTimer, mLogger)

	out, err := ucase.CreateAnswer(ctx, 7, "u1", "xxx")

	require.Nil(t, out)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create answer")
}
