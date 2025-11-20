package question

import (
	"context"

	ent "test-question/internal/entity/question"
	"test-question/internal/pkg/uow"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]*ent.Question, error) {
	var rows []questionRow

	err := r.db.WithContext(ctx).Find(&rows).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}

	res := make([]*ent.Question, 0, len(rows))
	for _, row := range rows {
		res = append(res, toEntityQuestion(&row))
	}

	return res, nil
}

func (r *Repository) Create(ctx context.Context, e *ent.Question) (*ent.Question, error) {
	row := fromEntityQuestion(e)

	if err := uow.GetTx(ctx, r.db).WithContext(ctx).
		Create(row).Error; err != nil {
		return nil, err
	}

	return toEntityQuestion(row), nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	return uow.GetTx(ctx, r.db).WithContext(ctx).
		Where("id = ?", id).
		Delete(&questionRow{}).Error
}

func (r *Repository) GetByID(ctx context.Context, id int) (*ent.Question, error) {
	var row questionRow

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ent.ErrQuestionNotFound
		}
		return nil, err
	}

	return toEntityQuestion(&row), nil
}
