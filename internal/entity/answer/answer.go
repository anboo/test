package answer

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrAnswerNotFound            = errors.New("answer not found")
	ErrRequestedQuestionNotFound = errors.New("requested question not found")
	ErrAccessDenied              = errors.New("access denied")
)

type Answer struct {
	ID         int
	QuestionID int
	UserID     string
	Text       string
	CreatedAt  time.Time
}
