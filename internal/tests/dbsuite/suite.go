package dbsuite

import (
	"test-question/internal/tests/infrasuite"
)

type DBSuite struct {
	infrasuite.InfraSuite
}

func (s *DBSuite) ResetTables(tables ...string) {
	for _, tbl := range tables {
		s.Require().NoError(
			s.DB.Exec("TRUNCATE TABLE " + tbl + " RESTART IDENTITY CASCADE").Error,
		)
	}
}
