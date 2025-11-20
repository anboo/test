package infra

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Env struct {
	LogLevelGorm string `env:"LOG_LEVEL_GORM" envDefault:"error"`

	LogLevel  string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"text"`

	ListenPort    string `env:"LISTEN_PORT" envDefault:":8000"`
	DbDSN         string `env:"DB_DSN,required"`
	MigrationPath string `env:"MIGRATION_PATH" envDefault:"./migration"`
}

func (r *Resources) initEnv() error {
	if err := godotenv.Load(".env.local", ".env"); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("dot env load: %w", err)
	}

	if err := env.Parse(&r.Env); err != nil {
		return fmt.Errorf("env parse: %w", err)
	}

	return nil
}
