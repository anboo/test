package create_question

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	entQ "test-question/internal/entity/question"
	"test-question/internal/pkg/rpc/rpc_auth"
	"test-question/internal/rpc/question/create_question/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create_Unauthorized(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	h := NewHandler(mUC)

	req := httptest.NewRequest("POST", "/questions", bytes.NewBufferString(`{"text":"hello"}`))
	req.Header.Set("Content-Type", "application/json")

	// В контексте НЕТ user_id
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_Create_ValidationError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	h := NewHandler(mUC)

	req := httptest.NewRequest("POST", "/questions", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	ctx := rpc_auth.InjectUserID(req.Context(), "test-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var body map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)

	require.Equal(t, "validation_failed", body["message"])

	fields := body["fields"].(map[string]any) //nolint:forcetypeassert
	require.Equal(t, "required", fields["Text"])
}

func TestHandler_Create_UseCaseError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.On(
		"CreateQuestion",
		mock.AnythingOfType("*context.valueCtx"),
		"test-user",
		"hello",
	).Return(nil, errors.New("fail"))

	h := NewHandler(mUC)

	req := httptest.NewRequest("POST", "/questions", bytes.NewBufferString(`{"text":"hello"}`))
	req.Header.Set("Content-Type", "application/json")

	ctx := rpc_auth.InjectUserID(req.Context(), "test-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_Create_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.On(
		"CreateQuestion",
		mock.AnythingOfType("*context.valueCtx"),
		"test-user",
		"hello",
	).Return(&entQ.Question{ID: 10, Text: "hello"}, nil)

	h := NewHandler(mUC)

	req := httptest.NewRequest("POST", "/questions", bytes.NewBufferString(`{"text":"hello"}`))
	req.Header.Set("Content-Type", "application/json")

	ctx := rpc_auth.InjectUserID(req.Context(), "test-user")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp CreateQuestionResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, 10, resp.ID)
	require.Equal(t, "hello", resp.Text)
}
