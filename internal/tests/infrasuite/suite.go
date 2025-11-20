package infrasuite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InfraSuite struct {
	suite.Suite

	Ctx       context.Context //nolint:containedctx
	Container tc.Container
	DB        *gorm.DB
}

func (s *InfraSuite) SetupSuite() {
	s.Ctx = context.Background()

	req := tc.ContainerRequest{
		Image: "postgres:16",
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "testdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForListeningPort("5432/tcp").
			WithStartupTimeout(time.Second * 30),
	}

	c, err := tc.GenericContainer(s.Ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)
	s.Container = c

	host, err := c.Host(s.Ctx)
	s.Require().NoError(err)

	port, err := c.MappedPort(s.Ctx, "5432/tcp")
	s.Require().NoError(err)

	dsn := fmt.Sprintf(
		"postgres://user:pass@%s:%d/testdb?sslmode=disable",
		host, port.Int(),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	s.Require().NoError(err)
	s.DB = db

	sqlDB, err := db.DB()
	s.Require().NoError(err)

	mPath := resolveMigrationsPath(s.T())
	err = goose.Up(sqlDB, mPath)
	s.Require().NoError(err)
}

func (s *InfraSuite) TearDownSuite() {
	if s.Container != nil {
		_ = s.Container.Terminate(s.Ctx)
	}
}

func resolveMigrationsPath(t testing.TB) string { //nolint:thelper
	_, currFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot get runtime.Caller")
	}

	base := filepath.Dir(currFile) // .../internal/tests/infrasuite
	base = filepath.Dir(base)      // .../internal/tests
	base = filepath.Dir(base)      // .../internal
	base = filepath.Dir(base)      // .../

	mPath := filepath.Join(base, "migration")
	if _, err := os.Stat(mPath); err != nil {
		t.Fatalf("migration directory not found: %s", mPath)
	}

	return mPath
}
