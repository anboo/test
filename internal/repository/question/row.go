package question

import (
	"time"

	"test-question/internal/entity/question"

	"gorm.io/gorm"
)

type questionRow struct {
	ID        int64          `gorm:"primaryKey;column:id"`
	Text      string         `gorm:"column:text;type:text;not null"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);not null;index"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (questionRow) TableName() string {
	return "questions"
}

func toEntityQuestion(q *questionRow) *question.Question {
	if q == nil {
		return nil
	}
	return &question.Question{
		ID:        int(q.ID),
		Text:      q.Text,
		UserID:    q.UserID,
		CreatedAt: q.CreatedAt,
	}
}

func fromEntityQuestion(e *question.Question) *questionRow {
	if e == nil {
		return nil
	}
	return &questionRow{
		ID:        int64(e.ID),
		Text:      e.Text,
		UserID:    e.UserID,
		CreatedAt: e.CreatedAt,
	}
}
