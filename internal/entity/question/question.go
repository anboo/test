package question

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrQuestionNotFound = errors.New("question not found")
	ErrAccessDenied     = errors.New("access denied")
)

type Question struct {
	ID        int
	Text      string
	UserID    string
	CreatedAt time.Time
}
