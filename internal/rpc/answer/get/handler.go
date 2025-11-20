package get

import (
	"context"
	"net/http"
	"strconv"
	"time"

	entA "test-question/internal/entity/answer"
	"test-question/internal/pkg/rpc"

	"github.com/pkg/errors"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		GetAnswer(ctx context.Context, answerID int) (*entA.Answer, error)
	}
)

type Response struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	answerID, err := strconv.Atoi(idStr)
	if err != nil {
		rpc.WriteBadRequest(w, "invalid answer id")
		return
	}

	a, err := h.uc.GetAnswer(r.Context(), answerID)
	if err != nil {
		switch {
		case errors.Is(err, entA.ErrAnswerNotFound):
			rpc.WriteNotFound(w, "answer_not_found")
			return
		default:
			rpc.WriteUnexpectedError(w, err)
			return
		}
	}

	rpc.WriteJSON(w, http.StatusOK, Response{
		ID:        a.ID,
		Text:      a.Text,
		UserID:    a.UserID,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	})
}
