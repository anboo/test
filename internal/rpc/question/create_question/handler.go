package create_question

import (
	"context"
	"net/http"

	entQ "test-question/internal/entity/question"
	"test-question/internal/pkg/rpc"
	"test-question/internal/pkg/rpc/rpc_auth"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		CreateQuestion(ctx context.Context, userID, text string) (*entQ.Question, error)
	}
)

type CreateQuestionRequest struct {
	Text string `json:"text" validate:"required,min=1"`
}

type CreateQuestionResponse struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) http.Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req CreateQuestionRequest
	if !rpc.ShouldBindJSON(r, w, &req) {
		return
	}

	userID := rpc_auth.GetUserID(r.Context())
	if userID == "" {
		rpc.WriteUnauthorized(w)
		return
	}

	q, err := h.uc.CreateQuestion(r.Context(), userID, req.Text)
	if err != nil {
		rpc.WriteUnexpectedError(w, err)
		return
	}

	rpc.WriteJSON(w, http.StatusCreated, CreateQuestionResponse{
		ID:   q.ID,
		Text: q.Text,
	})
}
