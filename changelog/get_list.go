package changelog

import (
	"flow-records/mysql"
	"time"
)

type GetListQuery struct {
	Start     *time.Time `query:"start" validate:"omitempty"`
	End       *time.Time `query:"end" validate:"omitempty"`
	TodoId    *uint64    `query:"todo_id" validate:"omitempty,gte=1"`
	SchemeId  *uint64    `query:"scheme_id" validate:"omitempty,gte=1"`
	ProjectId *uint64    `query:"project_id" validate:"omitempty,gte=1"`
}

func GetList(userId uint64, q GetListQuery) (changeLogs []ChangeLog, err error) {
	// Generate query
	queryStr := "SELECT id, text, datetime, todo_id, scheme_id FROM changelogs WHERE user_id = ?"
	queryParams := []interface{}{userId}
	if q.Start != nil {
		queryStr += " AND datetime >= ?"
		queryParams = append(queryParams, q.Start.UTC())
	}
	if q.End != nil {
		queryStr += " AND datetime >= ?"
		queryParams = append(queryParams, q.End.UTC())
	}
	if q.TodoId != nil {
		queryStr += " AND todo_id = ?"
		queryParams = append(queryParams, *q.TodoId)
	}
	if q.SchemeId != nil {
		queryStr += " AND scheme_id = ?"
		queryParams = append(queryParams, *q.SchemeId)
	}
	if q.ProjectId != nil {
		queryStr += ` AND scheme_id IN (
			SELECT id FROM schemes WHERE project_id = ?
		)`
		queryParams = append(queryParams, *q.ProjectId)
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

	rows, err := stmtOut.Query(queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		c := ChangeLog{}
		err = rows.Scan(&c.Id, &c.Text, &c.Datetime, &c.TodoId, &c.SchemeId)
		if err != nil {
			return
		}

		changeLogs = append(changeLogs, c)
	}
	return
}
