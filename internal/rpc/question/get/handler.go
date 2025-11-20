package get

import (
	"context"
	"net/http"
	"strconv"
	"time"

	entQ "test-question/internal/entity/question"
	"test-question/internal/pkg/rpc"
	"test-question/internal/usecase/question/get_with_answers"

	"github.com/pkg/errors"
)

//go:generate mockery --name=useCase --output=mocks --outpkg=mocks --exported
type (
	useCase interface {
		GetQuestionWithAnswers(
			ctx context.Context,
			questionID int,
		) (*get_with_answers.QuestionWithAnswers, error)
	}
)

type Response struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt string    `json:"created_at"`
	UserID    string    `json:"user_id"`
	Answers   []Answers `json:"answers"`
}

type Answers struct {
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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rpc.WriteBadRequest(w, "invalid question id")
		return
	}

	q, err := h.uc.GetQuestionWithAnswers(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, entQ.ErrQuestionNotFound):
			rpc.WriteNotFound(w, "question_not_found")
			return
		default:
			rpc.WriteUnexpectedError(w, err)
			return
		}
	}

	answers := make([]Answers, len(q.Answers))
	for i, a := range q.Answers {
		answers[i] = Answers{
			ID:        a.ID,
			Text:      a.Text,
			UserID:    a.UserID,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
		}
	}

	resp := Response{
		ID:        q.Question.ID,
		Text:      q.Question.Text,
		CreatedAt: q.Question.CreatedAt.Format(time.RFC3339),
		UserID:    q.Question.UserID,
		Answers:   answers,
	}

	rpc.WriteJSON(w, http.StatusOK, resp)
}
