package delete_question

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	entQ "test-question/internal/entity/question"
	"test-question/internal/pkg/rpc/rpc_auth"
	"test-question/internal/rpc/question/delete_question/mocks"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func router(h http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("DELETE /questions/{id}", h)
	return mux
}

func reqWithUser(method, url string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	ctx := rpc_auth.InjectUserID(req.Context(), "user-1")
	return req.WithContext(ctx)
}

func TestDeleteQuestion_Success(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteQuestion",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			10,
			"user-1",
		).
		Return(nil)

	h := NewHandler(mUC)
	srv := router(h)

	req := reqWithUser("DELETE", "/questions/10")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteQuestion_InvalidID(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)

	srv := router(h)

	req := reqWithUser("DELETE", "/questions/abc")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteQuestion_Unauthorized(t *testing.T) {
	mUC := mocks.NewUseCase(t)
	h := NewHandler(mUC)
	srv := router(h)

	req := httptest.NewRequest("DELETE", "/questions/10", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteQuestion_NotFound(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteQuestion",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			99,
			"user-1",
		).
		Return(entQ.ErrQuestionNotFound)

	h := NewHandler(mUC)
	srv := router(h)

	req := reqWithUser("DELETE", "/questions/99")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteQuestion_AccessDenied(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteQuestion",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			88,
			"user-1",
		).
		Return(entQ.ErrAccessDenied)

	h := NewHandler(mUC)
	srv := router(h)

	req := reqWithUser("DELETE", "/questions/88")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteQuestion_InternalError(t *testing.T) {
	mUC := mocks.NewUseCase(t)

	mUC.
		On("DeleteQuestion",
			mock.MatchedBy(func(ctx context.Context) bool { return true }),
			7,
			"user-1",
		).
		Return(errors.New("stub"))

	h := NewHandler(mUC)
	srv := router(h)

	req := reqWithUser("DELETE", "/questions/7")
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
