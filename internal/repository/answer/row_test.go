package answer

import (
	"testing"
	"time"

	ent "test-question/internal/entity/answer"

	"github.com/stretchr/testify/require"
)

func TestAnswerConverters(t *testing.T) {
	now := time.Now()

	testsToEntity := []struct {
		name   string
		row    *answerRow
		entity *ent.Answer
	}{
		{
			name: "row_to_entity",
			row: &answerRow{
				ID:         10,
				QuestionID: 3,
				UserID:     "u1",
				Text:       "hello",
				CreatedAt:  now,
			},
			entity: &ent.Answer{
				ID:         10,
				QuestionID: 3,
				UserID:     "u1",
				Text:       "hello",
				CreatedAt:  now,
			},
		},
		{
			name:   "nil_row",
			row:    nil,
			entity: nil,
		},
	}

	for _, tt := range testsToEntity {
		t.Run("toEntity_"+tt.name, func(t *testing.T) {
			if tt.row == nil {
				require.Nil(t, toEntityAnswer(tt.row))
				return
			}
			e := toEntityAnswer(tt.row)
			require.Equal(t, tt.entity, e)
		})
	}

	testsTowRow := []struct {
		name   string
		entity *ent.Answer
		row    *answerRow
	}{
		{
			name: "entity_to_row",
			entity: &ent.Answer{
				ID:         20,
				QuestionID: 7,
				UserID:     "u2",
				Text:       "ok",
				CreatedAt:  now,
			},
			row: &answerRow{
				ID:         20,
				QuestionID: 7,
				UserID:     "u2",
				Text:       "ok",
				CreatedAt:  now,
			},
		},
		{
			name:   "nil_entity",
			entity: nil,
			row:    nil,
		},
	}

	for _, tt := range testsTowRow {
		t.Run("fromEntity_"+tt.name, func(t *testing.T) {
			if tt.entity == nil {
				require.Nil(t, fromEntityAnswer(tt.entity))
				return
			}
			r := fromEntityAnswer(tt.entity)
			require.Equal(t, tt.row, r)
		})
	}
}
