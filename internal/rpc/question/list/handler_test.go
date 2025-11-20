package list

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	entQ "test-question/internal/entity/question"
	"test-question/internal/rpc/question/list/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_List_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	now := time.Now()

	mUC.
		On("ListQuestions",
			mock.Anything, // ← ВАЖНО!
		).
		Return([]*entQ.Question{
			{ID: 1, Text: "hello", UserID: "u1", CreatedAt: now},
			{ID: 2, Text: "world", UserID: "u2", CreatedAt: now.Add(-time.Hour)},
		}, nil)

	h := NewHandler(mUC)

	req := httptest.NewRequest("GET", "/questions", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []ResponseItem
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp, 2)

	require.Equal(t, 1, resp[0].ID)
	require.Equal(t, "hello", resp[0].Text)
	require.Equal(t, "u1", resp[0].UserID)
	require.Equal(t, now.Format(time.RFC3339), resp[0].CreatedAt)
}

func TestHandler_List_Error(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("ListQuestions",
			mock.Anything, // ← ТАКЖЕ ВАЖНО
		).
		Return(nil, assertErr())

	h := NewHandler(mUC)

	req := httptest.NewRequest("GET", "/questions", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	require.Equal(t, "internal error", resp["message"])
}

func assertErr() error { return fmt.Errorf("boom") }
