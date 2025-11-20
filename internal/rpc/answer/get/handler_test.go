package get

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	entA "test-question/internal/entity/answer"
	"test-question/internal/rpc/answer/get/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Get_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	now := time.Now()

	mUC.
		On("GetAnswer",
			mock.Anything,
			10,
		).
		Return(&entA.Answer{
			ID:        10,
			Text:      "hi",
			UserID:    "u1",
			CreatedAt: now,
		}, nil)

	h := NewHandler(mUC)

	req := httptest.NewRequest("GET", "/answers/10", nil)
	req.SetPathValue("id", "10")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, 10, resp.ID)
	require.Equal(t, "hi", resp.Text)
	require.Equal(t, "u1", resp.UserID)
	require.Equal(t, now.Format(time.RFC3339), resp.CreatedAt)
}

func TestHandler_Get_InvalidID(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	req := httptest.NewRequest("GET", "/answers/abc", nil)
	req.SetPathValue("id", "abc")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "invalid answer id", resp["message"])
}

func TestHandler_Get_NotFound(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("GetAnswer",
			mock.Anything,
			77,
		).
		Return(nil, entA.ErrAnswerNotFound)

	h := NewHandler(mUC)

	req := httptest.NewRequest("GET", "/answers/77", nil)
	req.SetPathValue("id", "77")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "answer_not_found", resp["message"])
}

func TestHandler_Get_UnexpectedError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("GetAnswer",
			mock.Anything,
			88,
		).
		Return(nil, fmt.Errorf("boom"))

	h := NewHandler(mUC)

	req := httptest.NewRequest("GET", "/answers/88", nil)
	req.SetPathValue("id", "88")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "internal error", resp["message"])
}
