package list

import (
	"context"
	"net/http"
	"time"

	entQ "test-question/internal/entity/question"
	"test-question/internal/pkg/rpc"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		ListQuestions(ctx context.Context) ([]*entQ.Question, error)
	}
)

type ResponseItem struct {
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
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	qs, err := h.uc.ListQuestions(r.Context())
	if err != nil {
		rpc.WriteUnexpectedError(w, err)
		return
	}

	resp := make([]ResponseItem, len(qs))
	for i, q := range qs {
		resp[i] = ResponseItem{
			ID:        q.ID,
			Text:      q.Text,
			UserID:    q.UserID,
			CreatedAt: q.CreatedAt.Format(time.RFC3339),
		}
	}

	rpc.WriteJSON(w, http.StatusOK, resp)
}
