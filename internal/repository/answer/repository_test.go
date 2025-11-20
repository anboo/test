//go:build integration
// +build integration

package answer

import (
	"context"
	"testing"
	"time"

	ent "test-question/internal/entity/answer"
	entq "test-question/internal/entity/question"
	"test-question/internal/repository/question"
	"test-question/internal/tests/dbsuite"

	"github.com/stretchr/testify/suite"
)

type AnswerRepoInfraSuite struct {
	dbsuite.DBSuite
	repo     *Repository
	quesRepo *question.Repository
	question *entq.Question
}

func (s *AnswerRepoInfraSuite) SetupTest() {
	s.repo = &Repository{db: s.DB}
	s.quesRepo = question.NewRepository(s.DB)

	s.ResetTables("answers", "questions")

	q := &entq.Question{
		Text:      "Test Question",
		UserID:    "11111111-1111-1111-1111-111111111111",
		CreatedAt: time.Now(),
	}
	createdQ, err := s.quesRepo.Create(context.Background(), q)
	s.Require().NoError(err)
	s.question = createdQ
}

func (s *AnswerRepoInfraSuite) TestCreate() {
	in := &ent.Answer{
		QuestionID: s.question.ID,
		UserID:     "user-1",
		Text:       "Answer text",
		CreatedAt:  time.Now(),
	}

	out, err := s.repo.Create(context.Background(), in)
	s.Require().NoError(err)
	s.Require().Greater(out.ID, 0) //nolint:testifylint

	var row answerRow
	s.NoError(s.DB.First(&row, out.ID).Error) //nolint:testifylint

	s.Equal(out.Text, row.Text)
	s.Equal(out.UserID, row.UserID)
	s.Equal(out.QuestionID, int(row.QuestionID))
}

func (s *AnswerRepoInfraSuite) TestGetByID() {
	ar := &answerRow{
		QuestionID: int64(s.question.ID),
		UserID:     "u2",
		Text:       "hi",
		CreatedAt:  time.Now(),
	}
	s.Require().NoError(s.DB.Create(&ar).Error)

	out, err := s.repo.GetByID(context.Background(), int(ar.ID))
	s.Require().NoError(err)

	s.Equal(int(ar.ID), out.ID)
	s.Equal("u2", out.UserID)
	s.Equal("hi", out.Text)
}

func (s *AnswerRepoInfraSuite) TestDeleteSoft() {
	ar := &answerRow{
		QuestionID: int64(s.question.ID),
		UserID:     "user3",
		Text:       "delete me",
		CreatedAt:  time.Now(),
	}
	s.Require().NoError(s.DB.Create(ar).Error)

	err := s.repo.Delete(context.Background(), int(ar.ID))
	s.Require().NoError(err)

	var row answerRow

	err = s.DB.First(&row, ar.ID).Error
	s.Require().Error(err)

	err = s.DB.Unscoped().First(&row, ar.ID).Error
	s.Require().NoError(err)
	s.True(row.DeletedAt.Valid)
}

func (s *AnswerRepoInfraSuite) TestListByQuestionID() {
	now := time.Now()

	a1 := &answerRow{
		QuestionID: int64(s.question.ID),
		UserID:     "u1",
		Text:       "A1",
		CreatedAt:  now.Add(-time.Minute),
	}
	a2 := &answerRow{
		QuestionID: int64(s.question.ID),
		UserID:     "u2",
		Text:       "A2",
		CreatedAt:  now,
	}

	s.NoError(s.DB.Create(a1).Error)
	s.NoError(s.DB.Create(a2).Error)

	list, err := s.repo.ListByQuestionID(context.Background(), s.question.ID)
	s.Require().NoError(err)
	s.Require().Len(list, 2)

	s.Equal("A1", list[0].Text)
	s.Equal("A2", list[1].Text)
}

func TestAnswerRepoInfraSuite(t *testing.T) {
	s := &AnswerRepoInfraSuite{}
	suite.Run(t, s)
}
