package scheme

import (
	"database/sql"
	"flow-records/mysql"
)

type GetListQuery struct {
	ProjectId *uint64 `query:"project_id" validate:"omitempty,gte=1"`
	// Start     *time.Time `query:"start" validate:"omitempty"`
	// End       *time.Time `query:"end" validate:"omitempty"`
	// Embed     *string    `query:"embed" validate:"omitempty,oneof=records record.changelog"`
}

func GetList(userId uint64, q GetListQuery) (schemes []Scheme, err error) {
	// TODO: Embedding records
	// Generate query
	queryStr := "SELECT id, name, sum_graph, project_id FROM schemes WHERE user_id = ?"
	if q.ProjectId != nil {
		queryStr += " AND project_id = ?"
	}

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

	var rows *sql.Rows
	if q.ProjectId == nil {
		rows, err = stmtOut.Query(userId)
	} else {
		rows, err = stmtOut.Query(userId, *q.ProjectId)
	}
	if err != nil {
		return
	}

	for rows.Next() {
		// TODO: uint64に対応 (projectId)
		var (
			id        uint64
			name      string
			sumGraph  bool
			projectId sql.NullInt64
		)
		err = rows.Scan(&id, &name, &sumGraph, &projectId)
		if err != nil {
			return
		}

		s := Scheme{Id: id, Name: name, SumGraph: sumGraph}
		if projectId.Valid {
			projectIdTmp := uint64(projectId.Int64)
			s.ProjectId = &projectIdTmp
		}

		schemes = append(schemes, s)
	}

	return
}
