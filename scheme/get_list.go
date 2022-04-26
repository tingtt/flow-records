package scheme

import (
	"database/sql"
	"flow-records/mysql"
	"flow-records/record"
	"time"
)

type GetListQuery struct {
	ProjectId *uint64    `query:"project_id" validate:"omitempty,gte=1"`
	Start     *time.Time `query:"start" validate:"omitempty"`
	End       *time.Time `query:"end" validate:"omitempty"`
	Embed     *string    `query:"embed" validate:"omitempty,oneof=records record.changelog"`
}

func GetList(userId uint64, q GetListQuery) (schemes []Scheme, err error) {
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
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		s := Scheme{}
		err = rows.Scan(&s.Id, &s.Name, &s.SumGraph, &s.ProjectId)
		if err != nil {
			return
		}
		if q.Embed != nil {
			if *q.Embed == "records" {
				s.Records, err = record.GetList(userId, record.GetListQuery{SchemeId: &s.Id, Start: q.Start, End: q.End})
				if err != nil {
					return
				}
			} else if *q.Embed == "record.changelog" {
				embed := "changelog"
				s.Records, err = record.GetList(userId, record.GetListQuery{SchemeId: &s.Id, Start: q.Start, End: q.End, Embed: &embed})
				if err != nil {
					return
				}
			}
		}
		schemes = append(schemes, s)
	}

	return
}
