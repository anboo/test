package answer

import (
	"time"

	"test-question/internal/entity/answer"

	"gorm.io/gorm"
)

type answerRow struct {
	ID         int64          `gorm:"primaryKey;column:id"`
	QuestionID int64          `gorm:"column:question_id;not null"`
	UserID     string         `gorm:"column:user_id;type:text;not null"`
	Text       string         `gorm:"column:text;type:text;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (answerRow) TableName() string {
	return "answers"
}

func toEntityAnswer(a *answerRow) *answer.Answer {
	if a == nil {
		return nil
	}
	return &answer.Answer{
		ID:         int(a.ID),
		QuestionID: int(a.QuestionID),
		UserID:     a.UserID,
		Text:       a.Text,
		CreatedAt:  a.CreatedAt,
	}
}

func fromEntityAnswer(e *answer.Answer) *answerRow {
	if e == nil {
		return nil
	}
	return &answerRow{
		ID:         int64(e.ID),
		QuestionID: int64(e.QuestionID),
		UserID:     e.UserID,
		Text:       e.Text,
		CreatedAt:  e.CreatedAt,
	}
}
