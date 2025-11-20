package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"test-question/cmd"
	"test-question/internal/infra"
	"test-question/internal/tests/dbsuite"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type E2ESuite struct {
	dbsuite.DBSuite

	Ctx       context.Context //nolint:containedctx
	Resources *infra.Resources
	Server    *httptest.Server
	Client    *http.Client

	// roles
	currentUser *AuthUser
	Users       map[string]*AuthUser
}

// user representation for auth
type AuthUser struct {
	Username string
	Password string
	UserID   string
}

// ==========================
//   AUTH ROLE SWITCHERS
// ==========================

func (s *E2ESuite) IAmBob() *E2ESuite {
	s.currentUser = s.Users["bob"]
	return s
}

func (s *E2ESuite) IAmAlice() *E2ESuite {
	s.currentUser = s.Users["alice"]
	return s
}

func (s *E2ESuite) IAmNobody() *E2ESuite {
	s.currentUser = nil
	return s
}

// ==========================
//   POSTGRES CONTAINER
// ==========================

func (s *E2ESuite) startPostgres() (tc.Container, string) {
	ctx := context.Background()

	req := tc.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	host, _ := c.Host(ctx)
	mapped, _ := c.MappedPort(ctx, "5432/tcp")

	dsn := "postgres://user:pass@" + host + ":" + mapped.Port() + "/testdb?sslmode=disable"
	return c, dsn
}

func (s *E2ESuite) SetupSuite() {
	s.Ctx = context.Background()

	// --- database
	pg, dsn := s.startPostgres()
	s.T().Cleanup(func() { pg.Terminate(context.Background()) }) //nolint:errcheck,gosec

	// --- env
	os.Setenv("DB_DSN", dsn)                             //nolint:errcheck,gosec
	os.Setenv("LISTEN_PORT", ":9999")                    //nolint:errcheck,gosec
	os.Setenv("MIGRATION_PATH", resolveMigrationsPath()) //nolint:errcheck,gosec

	// --- init resources
	res, err := infra.Init(s.Ctx)
	s.Require().NoError(err)
	s.Resources = res
	s.DB = res.DB

	// roles created by migration add_test_users.sql
	s.Users = map[string]*AuthUser{
		"bob": {
			Username: "bob",
			Password: "bob123",
			UserID:   "11111111-1111-1111-1111-111111111111",
		},
		"alice": {
			Username: "alice",
			Password: "alice123",
			UserID:   "22222222-2222-2222-2222-222222222222",
		},
	}

	// default — no one is logged in
	s.IAmNobody()

	// --- start HTTP server
	srv := &http.Server{ //nolint:gosec
		Addr:    ":9999",
		Handler: cmd.SetupRouter(res),
	}
	s.Server = httptest.NewServer(srv.Handler)

	s.Client = &http.Client{Timeout: 5 * time.Second}
}

func (s *E2ESuite) TearDownSuite() {
	s.Server.Close()
}

// ==========================
//         HTTP HELPERS
// ==========================

func (s *E2ESuite) request(method, path string, body any) *http.Response {
	var r io.Reader
	if body != nil {
		b, _ := json.Marshal(body) //nolint:errchkjson
		r = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, s.Server.URL+path, r) //nolint:noctx
	s.Require().NoError(err)

	req.Header.Set("Content-Type", "application/json")

	// Apply BasicAuth based on role
	if s.currentUser != nil {
		req.SetBasicAuth(s.currentUser.Username, s.currentUser.Password)
	}

	resp, err := s.Client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *E2ESuite) GET(path string) *http.Response {
	return s.request("GET", path, nil)
}

func (s *E2ESuite) POST(path string, body any) *http.Response {
	return s.request("POST", path, body)
}

func (s *E2ESuite) DELETE(path string) *http.Response {
	return s.request("DELETE", path, nil)
}

// ==========================
//     MIGRATIONS PATH
// ==========================

func resolveMigrationsPath() string {
	_, currFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get runtime.Caller")
	}

	base := filepath.Dir(currFile)
	base = filepath.Dir(base)
	base = filepath.Dir(base)
	base = filepath.Dir(base)

	mPath := filepath.Join(base, "migration")
	if _, err := os.Stat(mPath); err != nil {
		panic(fmt.Sprintf("migration directory not found: %s", mPath))
	}

	return mPath
}
