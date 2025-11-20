package delete_test

import (
	"context"
	"errors"
	"testing"

	entQ "test-question/internal/entity/question"
	uc "test-question/internal/usecase/question/delete"
	mocks2 "test-question/internal/usecase/question/delete/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newMocks(t *testing.T) (*mocks2.QuestionRepository, *mocks2.AnswerRepository, *mocks2.UnitOfWork, *mocks2.Logger) { //nolint:thelper
	return mocks2.NewQuestionRepository(t),
		mocks2.NewAnswerRepository(t),
		mocks2.NewUnitOfWork(t),
		mocks2.NewLogger(t)
}

func TestDeleteQuestion_Success(t *testing.T) {
	ctx := context.Background()

	qRepo, aRepo, uow, log := newMocks(t)

	qRepo.
		On("GetByID", mock.Anything, 10).
		Return(&entQ.Question{
			ID:     10,
			UserID: "owner-1",
		}, nil)

	qRepo.
		On("Delete", mock.Anything, 10).
		Return(nil)

	aRepo.
		On("DeleteByQuestionID", mock.Anything, 10).
		Return(nil)

	uow.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error) //nolint:forcetypeassert
			require.NoError(t, fn(ctx))
		}).
		Return(nil)

	log.
		On("DebugContext",
			mock.Anything,
			"question deleted with all answers",
			"question_id", 10,
			"user_id", "owner-1",
		).Return()

	ucase := uc.NewUseCase(qRepo, aRepo, uow, log)

	err := ucase.DeleteQuestion(ctx, 10, "owner-1")
	require.NoError(t, err)
}

func TestDeleteQuestion_NotFound(t *testing.T) {
	ctx := context.Background()

	qRepo, aRepo, uow, log := newMocks(t)

	qRepo.
		On("GetByID", mock.Anything, 99).
		Return(nil, entQ.ErrQuestionNotFound)

	uow.AssertNotCalled(t, "Do")

	ucase := uc.NewUseCase(qRepo, aRepo, uow, log)

	err := ucase.DeleteQuestion(ctx, 99, "user-x")
	require.ErrorIs(t, err, entQ.ErrQuestionNotFound)
}

func TestDeleteQuestion_AccessDenied(t *testing.T) {
	ctx := context.Background()

	qRepo, aRepo, uow, log := newMocks(t)

	qRepo.
		On("GetByID", mock.Anything, 7).
		Return(&entQ.Question{
			ID:     7,
			UserID: "owner-7",
		}, nil)

	uow.AssertNotCalled(t, "Do")

	ucase := uc.NewUseCase(qRepo, aRepo, uow, log)

	err := ucase.DeleteQuestion(ctx, 7, "other-user")
	require.ErrorIs(t, err, entQ.ErrAccessDenied)
}

func TestDeleteQuestion_GetByIDError(t *testing.T) {
	ctx := context.Background()

	qRepo, aRepo, uow, log := newMocks(t)

	qRepo.
		On("GetByID", mock.Anything, 5).
		Return(nil, errors.New("db error"))

	uow.AssertNotCalled(t, "Do")

	ucase := uc.NewUseCase(qRepo, aRepo, uow, log)

	err := ucase.DeleteQuestion(ctx, 5, "u1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get question")
}

func TestDeleteQuestion_DeleteError(t *testing.T) {
	ctx := context.Background()

	qRepo, aRepo, uow, log := newMocks(t)

	qRepo.
		On("GetByID", mock.Anything, 12).
		Return(&entQ.Question{
			ID:     12,
			UserID: "u12",
		}, nil)

	qRepo.
		On("Delete", mock.Anything, 12).
		Return(errors.New("delete fail"))

	uow.
		On("Do", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context) error) //nolint:forcetypeassert
			_ = fn(ctx)
		}).
		Return(errors.New("delete fail"))

	ucase := uc.NewUseCase(qRepo, aRepo, uow, log)

	err := ucase.DeleteQuestion(ctx, 12, "u12")
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete fail")
}
