package record

import (
	"flow-records/changelog"
	"flow-records/mysql"
	"time"
)

type GetListQuery struct {
	TodoId   *uint64    `query:"todo_id" validate:"omitempty"`
	SchemeId *uint64    `query:"scheme_id" validate:"omitempty"`
	Start    *time.Time `query:"start" validate:"omitempty"`
	End      *time.Time `query:"end" validate:"omitempty"`
	Embed    *string    `query:"embed" validate:"omitempty,oneof=changelog"`
}

type embedChangeLog struct {
	Id       *uint64
	Text     *string
	Datetime *time.Time
}

func GetList(userId uint64, q GetListQuery) (records []Record, err error) {
	// Generate query
	var queryStr string
	queryParam := []interface{}{userId}
	if q.Embed != nil {
		queryStr =
			`SELECT records.id, records.value, changelogs.id, changelogs.text, changelogs.datetime, records.datetime, records.todo_id, records.scheme_id
			FROM records
			LEFT JOIN changelogs
			ON records.todo_id = changelogs.todo_id
			AND records.scheme_id = changelogs.scheme_id
			WHERE records.user_id = ?`
	} else {
		queryStr = "SELECT id, value, datetime, todo_id, scheme_id FROM records WHERE user_id = ?"
	}
	if q.TodoId != nil {
		queryStr += " AND records.todo_id = ?"
		queryParam = append(queryParam, *q.TodoId)
	}
	if q.SchemeId != nil {
		queryStr += " AND records.scheme_id = ?"
		queryParam = append(queryParam, *q.SchemeId)
	}
	if q.Start != nil {
		queryStr += " AND records.datetime >= ?"
		queryParam = append(queryParam, q.Start.UTC())
	}
	if q.End != nil {
		queryStr += " AND records.datetime <= ?"
		queryParam = append(queryParam, q.End.UTC())
	}
	queryStr += " ORDER BY records.datetime"

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

	rows, err := stmtOut.Query(queryParam...)
	if err != nil {
		return
	}

	for rows.Next() {
		r := Record{}

		if q.Embed != nil {
			tmpCl := embedChangeLog{}
			err = rows.Scan(&r.Id, &r.Value, &tmpCl.Id, &tmpCl.Text, &tmpCl.Datetime, &r.Datetime, &r.TodoId, &r.SchemeId)
			if err != nil {
				return
			}
			if tmpCl.Id != nil && tmpCl.Text != nil && tmpCl.Datetime != nil {
				r.ChangeLog = &changelog.ChangeLog{Id: *tmpCl.Id, Text: *tmpCl.Text, Datetime: *tmpCl.Datetime, TodoId: r.TodoId, SchemeId: r.SchemeId}
			}
		} else {
			err = rows.Scan(&r.Id, &r.Value, &r.Datetime, &r.TodoId, &r.SchemeId)
			if err != nil {
				return
			}
		}

		records = append(records, r)
	}

	return
}
