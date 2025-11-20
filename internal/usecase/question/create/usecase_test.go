package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	entQ "test-question/internal/entity/question"
	uc "test-question/internal/usecase/question/create"
	mocks2 "test-question/internal/usecase/question/create/mocks"

	"github.com/stretchr/testify/require"
)

func TestCreateQuestion_Success(t *testing.T) {
	ctx := context.Background()

	now := time.Date(2024, 11, 21, 10, 0, 0, 0, time.UTC)

	mRepo := mocks2.NewQuestionRepository(t)
	mTimer := mocks2.NewTimer(t)
	mLogger := mocks2.NewLogger(t)

	mTimer.
		On("Now").
		Return(now)

	expectedInput := &entQ.Question{
		Text:      "hello world",
		UserID:    "1",
		CreatedAt: now,
	}

	mRepo.
		On("Create", ctx, expectedInput).
		Return(&entQ.Question{
			ID:        101,
			Text:      "hello world",
			CreatedAt: now,
		}, nil)

	mLogger.
		On("DebugContext",
			ctx,
			"question created",
			"question_id", 101,
		).
		Return()

	ucase := uc.NewUseCase(mRepo, mTimer, mLogger)

	out, err := ucase.CreateQuestion(ctx, "1", "hello world")
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, 101, out.ID)
	require.Equal(t, "hello world", out.Text)
	require.Equal(t, now, out.CreatedAt)
}

func TestCreateQuestion_RepoError(t *testing.T) {
	ctx := context.Background()

	now := time.Date(2024, 11, 21, 10, 0, 0, 0, time.UTC)

	mRepo := mocks2.NewQuestionRepository(t)
	mTimer := mocks2.NewTimer(t)
	mLogger := mocks2.NewLogger(t)

	mTimer.
		On("Now").
		Return(now)

	expectedInput := &entQ.Question{
		Text:      "qqq",
		UserID:    "1",
		CreatedAt: now,
	}

	mRepo.
		On("Create", ctx, expectedInput).
		Return(nil, errors.New("db fail"))

	ucase := uc.NewUseCase(mRepo, mTimer, mLogger)

	out, err := ucase.CreateQuestion(ctx, "1", "qqq")

	require.Nil(t, out)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create question")
}
