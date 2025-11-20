package get_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	entA "test-question/internal/entity/answer"
	entQ "test-question/internal/entity/question"
	"test-question/internal/rpc/question/get"
	"test-question/internal/rpc/question/get/mocks"
	qwa "test-question/internal/usecase/question/get_with_answers"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Get_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	now := time.Now()

	mUC.
		On("GetQuestionWithAnswers",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			10,
		).
		Return(&qwa.QuestionWithAnswers{
			Question: &entQ.Question{
				ID:        10,
				Text:      "hello",
				UserID:    "user-1",
				CreatedAt: now,
			},
			Answers: []*entA.Answer{
				{
					ID:         1,
					QuestionID: 10,
					UserID:     "a1",
					Text:       "first",
					CreatedAt:  now.Add(time.Minute),
				},
			},
		}, nil)

	h := get.NewHandler(mUC)
	mux := http.NewServeMux()
	mux.Handle("GET /questions/{id}", h)

	req := httptest.NewRequest("GET", "/questions/10", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp get.Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Equal(t, 10, resp.ID)
	require.Equal(t, "hello", resp.Text)
	require.Equal(t, "user-1", resp.UserID)
	require.Equal(t, now.Format(time.RFC3339), resp.CreatedAt)

	require.Len(t, resp.Answers, 1)
	require.Equal(t, 1, resp.Answers[0].ID)
	require.Equal(t, "first", resp.Answers[0].Text)
	require.Equal(t, "a1", resp.Answers[0].UserID)
	require.Equal(t, now.Add(time.Minute).Format(time.RFC3339), resp.Answers[0].CreatedAt)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := get.NewHandler(mUC)

	mux := http.NewServeMux()
	mux.Handle("GET /questions/{id}", h)

	req := httptest.NewRequest("GET", "/questions/xxx", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "invalid question id", resp["message"])
}

func TestHandler_Get_NotFound(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("GetQuestionWithAnswers",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			55,
		).
		Return(nil, entQ.ErrQuestionNotFound)

	h := get.NewHandler(mUC)
	mux := http.NewServeMux()
	mux.Handle("GET /questions/{id}", h)

	req := httptest.NewRequest("GET", "/questions/55", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "question_not_found", resp["message"])
}

func TestHandler_Get_UnexpectedError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("GetQuestionWithAnswers",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			99,
		).
		Return(nil, fmt.Errorf("boom"))

	h := get.NewHandler(mUC)
	mux := http.NewServeMux()
	mux.Handle("GET /questions/{id}", h)

	req := httptest.NewRequest("GET", "/questions/99", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "internal error", resp["message"])
}
