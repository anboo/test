package list

import (
	"context"
	"errors"
	"testing"
	"time"

	entQ "test-question/internal/entity/question"
	mocks2 "test-question/internal/usecase/question/list/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_ListQuestions_Success(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks2.NewQuestionRepository(t)
	mLogger := mocks2.NewLogger(t)

	now := time.Now()

	mRepo.
		On("List", mock.MatchedBy(func(c context.Context) bool { return true })).
		Return([]*entQ.Question{
			{ID: 1, Text: "q1", UserID: "u1", CreatedAt: now},
			{ID: 2, Text: "q2", UserID: "u2", CreatedAt: now},
		}, nil)

	mLogger.
		On("DebugContext",
			mock.Anything,
			"questions listed",
			"count", 2,
		).
		Return()

	uc := NewUseCase(mRepo, mLogger)

	out, err := uc.ListQuestions(ctx)
	require.NoError(t, err)
	require.Len(t, out, 2)

	require.Equal(t, 1, out[0].ID)
	require.Equal(t, "q1", out[0].Text)
	require.Equal(t, "u1", out[0].UserID)
}

func TestUseCase_ListQuestions_Error(t *testing.T) {
	ctx := context.Background()

	mRepo := mocks2.NewQuestionRepository(t)
	mLogger := mocks2.NewLogger(t)

	mRepo.
		On("List", mock.MatchedBy(func(c context.Context) bool { return true })).
		Return(nil, errors.New("db_fail"))

	uc := NewUseCase(mRepo, mLogger)

	out, err := uc.ListQuestions(ctx)
	require.Nil(t, out)
	require.ErrorContains(t, err, "list questions: db_fail")
}
