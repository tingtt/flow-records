package scheme

import (
	"flow-records/mysql"
	"flow-records/record"
	"time"
)

type GetQuery struct {
	Embed *string    `query:"embed" validate:"omitempty,oneof=records record.changelog"`
	Start *time.Time `query:"start" validate:"omitempty"`
	End   *time.Time `query:"end" validate:"omitempty"`
}

func Get(userId uint64, id uint64, q GetQuery) (s Scheme, notFound bool, err error) {
	// Generate query
	queryStr := "SELECT name, sum_graph, project_id FROM schemes WHERE user_id = ? AND id = ?"
	queryParams := []interface{}{userId, id}

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

	rows, err := stmtOut.Query(queryParams)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&s.Name, &s.SumGraph, &s.ProjectId)
	if err != nil {
		return
	}

	if q.Embed != nil {
		if *q.Embed == "records" {
			s.Records, err = record.GetList(userId, record.GetListQuery{SchemeId: &id, Start: q.Start, End: q.End})
			if err != nil {
				return
			}
		} else if *q.Embed == "record.changelog" {
			embed := "changelog"
			s.Records, err = record.GetList(userId, record.GetListQuery{SchemeId: &id, Start: q.Start, End: q.End, Embed: &embed})
			if err != nil {
				return
			}
		}
	}

	s.Id = id
	return
}
