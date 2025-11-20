package delete_answer

import (
	"context"
	"errors"
	"testing"

	entA "test-question/internal/entity/answer"
	"test-question/internal/usecase/answer/delete/mocks"

	"github.com/stretchr/testify/require"
)

func TestDeleteAnswer_Success(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID", ctx, 10).
		Return(&entA.Answer{
			ID:     10,
			UserID: "owner-1",
		}, nil)

	mRepo.
		On("Delete", ctx, 10).
		Return(nil)

	mLogger.
		On("DebugContext",
			ctx,
			"answer deleted",
			"answer_id", 10,
			"user_id", "owner-1",
		).Return()

	ucase := NewUseCase(mRepo, mLogger)

	err := ucase.DeleteAnswer(ctx, 10, "owner-1")
	require.NoError(t, err)
}

func TestDeleteAnswer_NotFound(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID", ctx, 99).
		Return(nil, entA.ErrAnswerNotFound)

	ucase := NewUseCase(mRepo, mLogger)

	err := ucase.DeleteAnswer(ctx, 99, "user-x")
	require.ErrorIs(t, err, entA.ErrAnswerNotFound)
}

func TestDeleteAnswer_AccessDenied(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID", ctx, 7).
		Return(&entA.Answer{
			ID:     7,
			UserID: "owner-7",
		}, nil)

	ucase := NewUseCase(mRepo, mLogger)

	err := ucase.DeleteAnswer(ctx, 7, "another-user")
	require.ErrorIs(t, err, entA.ErrAccessDenied)
}

func TestDeleteAnswer_GetByIDError(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID", ctx, 5).
		Return(nil, errors.New("db down"))

	ucase := NewUseCase(mRepo, mLogger)

	err := ucase.DeleteAnswer(ctx, 5, "u1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get answer")
}

func TestDeleteAnswer_DeleteError(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks.NewAnswerRepository(t)
	mLogger := mocks.NewLogger(t)

	mRepo.
		On("GetByID", ctx, 12).
		Return(&entA.Answer{
			ID:     12,
			UserID: "user12",
		}, nil)

	mRepo.
		On("Delete", ctx, 12).
		Return(errors.New("delete fail"))

	ucase := NewUseCase(mRepo, mLogger)

	err := ucase.DeleteAnswer(ctx, 12, "user12")
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete answer")
}
