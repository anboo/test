package create

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"test-question/internal/entity/answer"
	"test-question/internal/pkg/rpc"
	"test-question/internal/pkg/rpc/rpc_auth"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		CreateAnswer(
			ctx context.Context,
			questionID int,
			userID string,
			text string,
		) (*answer.Answer, error)
	}
)

type CreateAnswerRequest struct {
	Text string `json:"text" validate:"required,min=1"`
}

type CreateAnswerResponse struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	UserID     string `json:"user_id"`
	QuestionID int    `json:"question_id"`
}

type Handler struct {
	uc useCase
}

func NewHandler(uc useCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req CreateAnswerRequest
	if !rpc.ShouldBindJSON(r, w, &req) {
		return
	}

	qIDStr := r.PathValue("id")
	qID, err := strconv.Atoi(qIDStr)
	if err != nil {
		rpc.WriteBadRequest(w, "invalid question id")
		return
	}

	userID := rpc_auth.GetUserID(r.Context())
	if userID == "" {
		rpc.WriteUnauthorized(w)
		return
	}

	a, err := h.uc.CreateAnswer(r.Context(), qID, userID, req.Text)
	if err != nil {
		switch {
		case errors.Is(err, answer.ErrRequestedQuestionNotFound):
			rpc.WriteNotFound(w, "question_not_found")
			return
		default:
			rpc.WriteUnexpectedError(w, err)
			return
		}
	}

	rpc.WriteJSON(w, http.StatusCreated, CreateAnswerResponse{
		ID:         a.ID,
		Text:       a.Text,
		UserID:     a.UserID,
		QuestionID: a.QuestionID,
	})
}
