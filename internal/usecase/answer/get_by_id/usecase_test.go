package get_by_id

import (
	"context"
	"errors"
	"testing"
	"time"

	entA "test-question/internal/entity/answer"
	"test-question/internal/usecase/answer/get_by_id/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAnswer_Success(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	now := time.Now()

	mRepo.
		On("GetByID",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			10,
		).
		Return(&entA.Answer{
			ID:        10,
			Text:      "hello",
			UserID:    "u1",
			CreatedAt: now,
		}, nil)

	mLogger.
		On("DebugContext",
			mock.Anything,
			"answer loaded",
			"answer_id", 10,
			"user_id", "u1",
		)

	uc := NewUseCase(mRepo, mLogger)

	out, err := uc.GetAnswer(ctx, 10)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, 10, out.ID)
	require.Equal(t, "hello", out.Text)
	require.Equal(t, "u1", out.UserID)
}

func TestGetAnswer_NotFound(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			50,
		).
		Return(nil, entA.ErrAnswerNotFound)

	uc := NewUseCase(mRepo, mLogger)

	out, err := uc.GetAnswer(ctx, 50)
	require.Nil(t, out)
	require.ErrorIs(t, err, entA.ErrAnswerNotFound)
}

func TestGetAnswer_GetError(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			77,
		).
		Return(nil, errors.New("db down"))

	uc := NewUseCase(mRepo, mLogger)

	out, err := uc.GetAnswer(ctx, 77)
	require.Nil(t, out)
	require.Error(t, err)
	require.Contains(t, err.Error(), "get answer")
}

func TestGetAnswer_LoggerCalled(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			1,
		).
		Return(&entA.Answer{
			ID:     1,
			Text:   "t",
			UserID: "u1",
		}, nil)

	mLogger.
		On("DebugContext",
			mock.Anything,
			"answer loaded",
			"answer_id", 1,
			"user_id", "u1")

	uc := NewUseCase(mRepo, mLogger)

	_, err := uc.GetAnswer(ctx, 1)
	require.NoError(t, err)
}
