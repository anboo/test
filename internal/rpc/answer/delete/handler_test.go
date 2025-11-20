package delete //nolint:predeclared

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	entA "test-question/internal/entity/answer"
	"test-question/internal/pkg/rpc/rpc_auth"
	"test-question/internal/rpc/answer/delete/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Delete_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteAnswer",
			mock.AnythingOfType("*context.valueCtx"),
			10,
			"user-1",
		).
		Return(nil)

	h := NewHandler(mUC)

	req := httptest.NewRequest("DELETE", "/answers/10", nil)
	req.SetPathValue("id", "10")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestHandler_Delete_InvalidID(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	req := httptest.NewRequest("DELETE", "/answers/abc", nil)
	req.SetPathValue("id", "abc")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-1"))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "invalid answer id", resp["message"])
}

func TestHandler_Delete_Unauthorized(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	req := httptest.NewRequest("DELETE", "/answers/10", nil)
	req.SetPathValue("id", "10")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	require.Equal(t, "unauthorized", resp["message"])
}

func TestHandler_Delete_NotFound(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteAnswer",
			mock.AnythingOfType("*context.valueCtx"),
			55,
			"user-x",
		).
		Return(entA.ErrAnswerNotFound)

	h := NewHandler(mUC)

	req := httptest.NewRequest("DELETE", "/answers/55", nil)
	req.SetPathValue("id", "55")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-x"))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "answer_not_found", resp["message"])
}

func TestHandler_Delete_AccessDenied(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteAnswer",
			mock.AnythingOfType("*context.valueCtx"),
			77,
			"user-2",
		).
		Return(entA.ErrAccessDenied)

	h := NewHandler(mUC)

	req := httptest.NewRequest("DELETE", "/answers/77", nil)
	req.SetPathValue("id", "77")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-2"))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "access_denied", resp["message"])
}

func TestHandler_Delete_UnexpectedError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteAnswer",
			mock.AnythingOfType("*context.valueCtx"),
			99,
			"user-e",
		).
		Return(fmt.Errorf("boom"))

	h := NewHandler(mUC)

	req := httptest.NewRequest("DELETE", "/answers/99", nil)
	req.SetPathValue("id", "99")
	req = req.WithContext(rpc_auth.InjectUserID(req.Context(), "user-e"))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "internal error", resp["message"])
}
