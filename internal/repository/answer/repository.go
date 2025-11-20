package answer

import (
	"context"
	"errors"

	ent "test-question/internal/entity/answer"
	"test-question/internal/pkg/uow"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, e *ent.Answer) (*ent.Answer, error) {
	row := fromEntityAnswer(e)

	if err := uow.GetTx(ctx, r.db).WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}

	return toEntityAnswer(row), nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*ent.Answer, error) {
	var row answerRow

	err := r.db.WithContext(ctx).First(&row, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ent.ErrAnswerNotFound
		}
		return nil, err
	}

	return toEntityAnswer(&row), nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	return uow.GetTx(ctx, r.db).WithContext(ctx).Delete(&answerRow{}, id).Error
}

func (r *Repository) DeleteByQuestionID(ctx context.Context, questionID int) error {
	err := uow.GetTx(ctx, r.db).WithContext(ctx).
		Where("question_id = ?", questionID).
		Delete(&answerRow{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByQuestionID(ctx context.Context, questionID int) ([]*ent.Answer, error) {
	var rows []answerRow

	err := r.db.WithContext(ctx).
		Where("question_id = ?", questionID).
		Order("created_at ASC").
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	out := make([]*ent.Answer, 0, len(rows))
	for i := range rows {
		out = append(out, toEntityAnswer(&rows[i]))
	}

	return out, nil
}
