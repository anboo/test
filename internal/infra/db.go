package infra

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (r *Resources) initDb(ctx context.Context) error {
	db, err := gorm.Open(postgres.Open(r.Env.DbDSN), &gorm.Config{
		Logger: logger.Default.LogMode(parseGormLevel(r.Env.LogLevelGorm)),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return errors.Wrap(err, "try init db")
	}
	r.DB = db

	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "try get db")
	}

	err = goose.UpContext(ctx, sqlDB, r.Env.MigrationPath)
	if err != nil {
		return errors.Wrap(err, "run migrations")
	}

	return nil
}

func parseGormLevel(lvl string) logger.LogLevel {
	switch strings.ToLower(lvl) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	case "info", "debug":
		return logger.Info
	default:
		return logger.Error
	}
}
