package user

import (
	"time"

	ent "test-question/internal/entity/user"

	"gorm.io/gorm"
)

type userRow struct {
	ID        string         `gorm:"primaryKey;column:id"`
	Username  string         `gorm:"column:username;unique"`
	Password  string         `gorm:"column:password"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (userRow) TableName() string {
	return "users"
}

func toEntityUser(r *userRow) *ent.User {
	if r == nil {
		return nil
	}

	return &ent.User{
		ID:        r.ID,
		Username:  r.Username,
		Password:  r.Password,
		CreatedAt: r.CreatedAt,
	}
}

func fromEntityUser(e *ent.User) *userRow {
	if e == nil {
		return nil
	}

	return &userRow{
		ID:        e.ID,
		Username:  e.Username,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
	}
}
