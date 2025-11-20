package question

import (
	"testing"
	"time"

	ent "test-question/internal/entity/question"

	"github.com/stretchr/testify/require"
)

func TestQuestionConverters(t *testing.T) {
	now := time.Now()

	testsToEntity := []struct {
		name   string
		row    *questionRow
		entity *ent.Question
	}{
		{
			name: "row_to_entity",
			row: &questionRow{
				ID:        1,
				Text:      "hi",
				UserID:    "1",
				CreatedAt: now,
			},
			entity: &ent.Question{
				ID:        1,
				Text:      "hi",
				UserID:    "1",
				CreatedAt: now,
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
				require.Nil(t, toEntityQuestion(tt.row))
				return
			}
			e := toEntityQuestion(tt.row)
			require.Equal(t, tt.entity, e)
		})
	}

	testsToRow := []struct {
		name   string
		entity *ent.Question
		row    *questionRow
	}{
		{
			name: "entity_to_row",
			entity: &ent.Question{
				ID:        2,
				Text:      "yo",
				UserID:    "1",
				CreatedAt: now,
			},
			row: &questionRow{
				ID:        2,
				Text:      "yo",
				UserID:    "1",
				CreatedAt: now,
			},
		},
		{
			name:   "nil_entity",
			entity: nil,
			row:    nil,
		},
	}

	for _, tt := range testsToRow {
		t.Run("fromEntity_"+tt.name, func(t *testing.T) {
			if tt.entity == nil {
				require.Nil(t, fromEntityQuestion(tt.entity))
				return
			}
			r := fromEntityQuestion(tt.entity)
			require.Equal(t, tt.row, r)
		})
	}
}
