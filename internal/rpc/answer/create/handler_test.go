package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	entA "test-question/internal/entity/answer"
	"test-question/internal/pkg/rpc/rpc_auth"
	"test-question/internal/rpc/answer/create/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	now := time.Now()

	mUC.
		On("CreateAnswer",
			mock.AnythingOfType("*context.valueCtx"),
			10,
			"user-1",
			"hello answer",
		).
		Return(&entA.Answer{
			ID:         100,
			Text:       "hello answer",
			UserID:     "user-1",
			QuestionID: 10,
			CreatedAt:  now,
		}, nil)

	h := NewHandler(mUC)

	body := `{"text":"hello answer"}`
	req := httptest.NewRequest("POST", "/questions/10/answers", bytes.NewBufferString(body))

	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))
	req.SetPathValue("id", "10")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp CreateAnswerResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Equal(t, 100, resp.ID)
	require.Equal(t, "hello answer", resp.Text)
	require.Equal(t, "user-1", resp.UserID)
	require.Equal(t, 10, resp.QuestionID)
}

func TestHandler_Create_InvalidQuestionID(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	body := `{"text":"xxx"}`

	req := httptest.NewRequest("POST", "/questions/abc/answers", bytes.NewBufferString(body))
	req.SetPathValue("id", "abc")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "invalid question id", resp["message"])
}

func TestHandler_Create_Unauthorized(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	body := `{"text":"xxx"}`

	req := httptest.NewRequest("POST", "/questions/10/answers", bytes.NewBufferString(body))
	req.SetPathValue("id", "10")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "unauthorized", resp["message"])
}

func TestHandler_Create_QuestionNotFound(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("CreateAnswer",
			mock.Anything,
			55,
			"user-1",
			"hi",
		).
		Return(nil, entA.ErrRequestedQuestionNotFound)

	h := NewHandler(mUC)

	body := `{"text":"hi"}`

	req := httptest.NewRequest("POST", "/questions/55/answers", bytes.NewBufferString(body))
	req.SetPathValue("id", "55")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "question_not_found", resp["message"])
}

func TestHandler_Create_ValidationError(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	body := `{}`

	req := httptest.NewRequest("POST", "/questions/10/answers", bytes.NewBufferString(body))
	req.SetPathValue("id", "10")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "validation_failed", resp["message"])
}

func TestHandler_Create_UnexpectedError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("CreateAnswer",
			mock.Anything,
			99,
			"user-1",
			"xx",
		).
		Return(nil, fmt.Errorf("boom"))

	h := NewHandler(mUC)

	body := `{"text":"xx"}`

	req := httptest.NewRequest("POST", "/questions/99/answers", bytes.NewBufferString(body))
	req.SetPathValue("id", "99")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "internal error", resp["message"])
}
