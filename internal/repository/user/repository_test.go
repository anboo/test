//go:build integration
// +build integration

package user

import (
	"context"
	"testing"
	"time"

	ent "test-question/internal/entity/user"
	"test-question/internal/tests/dbsuite"

	"github.com/stretchr/testify/suite"
)

type UserRepoInfraSuite struct {
	dbsuite.DBSuite
	repo *Repository
}

func (s *UserRepoInfraSuite) SetupTest() {
	s.repo = &Repository{db: s.DB}
	s.ResetTables("users")
}

func (s *UserRepoInfraSuite) TestCreateUser() {
	in := &ent.User{
		ID:        "11111111-1111-1111-1111-111111111111",
		Username:  "john",
		Password:  "pass123",
		CreatedAt: time.Now(),
	}

	out, err := s.repo.CreateUser(context.Background(), in)
	s.Require().NoError(err)
	s.Require().NotNil(out)
	s.Equal(in.ID, out.ID)
	s.Equal(in.Username, out.Username)
	s.Equal(in.Password, out.Password)

	var row userRow
	err = s.DB.First(&row, "id = ?", in.ID).Error
	s.Require().NoError(err)
	s.Equal("john", row.Username)
	s.Equal("pass123", row.Password)
}

func (s *UserRepoInfraSuite) TestGetUserByUsernamePassword() {
	row := &userRow{
		ID:        "22222222-2222-2222-2222-222222222222",
		Username:  "max",
		Password:  "qwerty",
		CreatedAt: time.Now(),
	}
	s.Require().NoError(s.DB.Create(row).Error)

	out, err := s.repo.GetUserByUsernamePassword(context.Background(), "max", "qwerty")
	s.Require().NoError(err)
	s.Require().NotNil(out)

	s.Equal("max", out.Username)
	s.Equal("qwerty", out.Password)
	s.Equal("22222222-2222-2222-2222-222222222222", out.ID)
}

func (s *UserRepoInfraSuite) TestGetUserByUsernamePassword_NotFound() {
	_, err := s.repo.GetUserByUsernamePassword(context.Background(), "nosuchuser", "badpass")
	s.Require().Error(err)
	s.ErrorIs(err, ErrUserNotFound)
}

func TestUserRepoInfraSuite(t *testing.T) {
	s := &UserRepoInfraSuite{}
	suite.Run(t, s)
}
