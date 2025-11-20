package uow

import (
	"context"

	"gorm.io/gorm"
)

type GormUnitOfWork struct {
	db *gorm.DB
}

func NewGormUoW(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{db: db}
}

func (u *GormUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, "tx", tx) //nolint:fatcontext,staticcheck
		return fn(ctx)
	})
}

func GetTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
