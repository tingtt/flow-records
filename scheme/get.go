package scheme

import (
	"database/sql"
	"flow-records/mysql"
)

type GetQuery struct {
	// Start *time.Time `query:"start" validate:"omitempty"`
	// End   *time.Time `query:"end" validate:"omitempty"`
	// Embed *string    `query:"embed" validate:"omitempty,oneof=records record.changelog"`
}

func Get(userId uint64, id uint64, q GetQuery) (s Scheme, notFound bool, err error) {
	// TODO: Embedding records
	// Generate query
	queryStr := "SELECT name, sum_graph, project_id FROM schemes WHERE user_id = ? AND id = ?"

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return
	}

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	// TODO: uint64に対応 (projectId)
	var (
		name      string
		sumGraph  bool
		projectId sql.NullInt64
	)
	err = rows.Scan(&name, &sumGraph, &projectId)
	if err != nil {
		return
	}

	s.Id = id
	s.Name = name
	s.SumGraph = sumGraph
	if projectId.Valid {
		projectIdTmp := uint64(projectId.Int64)
		s.ProjectId = &projectIdTmp
	}
	return
}
