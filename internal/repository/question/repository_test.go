//go:build integration
// +build integration

package question

import (
	"context"
	"testing"
	"time"

	ent "test-question/internal/entity/question"
	"test-question/internal/tests/dbsuite"

	"github.com/stretchr/testify/suite"
)

type QuestionRepoInfraSuite struct {
	dbsuite.DBSuite
	repo *Repository
}

func (s *QuestionRepoInfraSuite) SetupTest() {
	s.repo = &Repository{db: s.DB}
	s.ResetTables("answers", "questions")
}

func (s *QuestionRepoInfraSuite) TestCreate() {
	in := &ent.Question{
		Text:      "hello world",
		UserID:    "11111111-1111-1111-1111-111111111111",
		CreatedAt: time.Now(),
	}

	out, err := s.repo.Create(context.Background(), in)
	s.Require().NoError(err)
	s.Require().NotNil(out)
	s.Require().Greater(out.ID, 0) //nolint:testifylint

	var row questionRow
	err = s.DB.First(&row, out.ID).Error
	s.Require().NoError(err)

	s.Equal(out.ID, int(row.ID))
	s.Equal(out.Text, row.Text)
	s.Equal(out.UserID, row.UserID)
	s.WithinDuration(out.CreatedAt, row.CreatedAt, time.Second)
}

func (s *QuestionRepoInfraSuite) TestGetByID() {
	now := time.Now()

	q := &questionRow{
		Text:      "test fetch",
		UserID:    "11111111-1111-1111-1111-111111111111",
		CreatedAt: now,
	}
	s.Require().NoError(s.DB.Create(q).Error)

	out, err := s.repo.GetByID(context.Background(), int(q.ID))
	s.Require().NoError(err)

	s.Equal(int(q.ID), out.ID)
	s.Equal(q.Text, out.Text)
	s.Equal(q.UserID, out.UserID)
	s.WithinDuration(q.CreatedAt, out.CreatedAt, time.Second)
}

func (s *QuestionRepoInfraSuite) TestDelete() {
	q := &questionRow{
		Text:      "to delete",
		UserID:    "11111111-1111-1111-1111-111111111111",
		CreatedAt: time.Now(),
	}
	s.Require().NoError(s.DB.Create(q).Error)

	err := s.repo.Delete(context.Background(), int(q.ID))
	s.Require().NoError(err)

	var row questionRow
	err = s.DB.First(&row, q.ID).Error
	s.Require().Error(err)

	err = s.DB.Unscoped().First(&row, q.ID).Error
	s.Require().NoError(err)
	s.NotZero(row.DeletedAt.Time)
	s.True(row.DeletedAt.Valid)
}

func TestQuestionRepoInfraSuite(t *testing.T) {
	s := &QuestionRepoInfraSuite{}
	suite.Run(t, s)
}
