package user

import (
	"context"
	"errors"

	ent "test-question/internal/entity/user"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, u *ent.User) (*ent.User, error) {
	row := fromEntityUser(u)

	if err := r.db.WithContext(ctx).Create(row).Error; err != nil {
		return nil, err
	}

	return toEntityUser(row), nil
}

func (r *Repository) GetUserByUsernamePassword(ctx context.Context, username, password string) (*ent.User, error) {
	var row userRow

	err := r.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password).
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return toEntityUser(&row), nil
}
