package record

import (
	"flow-records/changelog"
	"flow-records/mysql"
	"strings"
	"time"
)

type postBodyWithoutTodoId struct {
	Value     *int      `json:"value" validate:"required"`
	ChangeLog *string   `json:"changelog" validate:"omitempty"`
	Datetime  time.Time `json:"datetime" validate:"required"`
	SchemeId  uint64    `json:"scheme_id" validate:"required,gte=1"`
}

type MultiplePostBody struct {
	Records []postBodyWithoutTodoId `json:"records" validate:"required,gte=1,dive,required"`
	TodoId  *uint64                 `json:"todo_id" validate:"omitempty,gte=1"`
}

func PostMultiple(userId uint64, p MultiplePostBody) (records []Record, err error) {
	// Generate query
	queryStr := "INSERT INTO records (user_id, value, datetime, scheme_id, todo_id) VALUES"
	queryParams := []interface{}{}
	for _, r := range p.Records {
		queryStr += " (?, ?, ?, ?, ?),"
		queryParams = append(queryParams, userId, r.Value, r.Datetime.UTC(), r.SchemeId, p.TodoId)
	}
	queryStr = strings.TrimRight(queryStr, ",")

	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(queryParams...)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	for i, r := range p.Records {
		r2 := Record{Id: uint64(id) + uint64(i), Value: *r.Value, Datetime: r.Datetime, SchemeId: r.SchemeId, TodoId: p.TodoId}

		// Post changelog
		if r.ChangeLog != nil {
			var cl changelog.ChangeLog
			cl, err = changelog.Post(userId, changelog.PostBody{Text: *r.ChangeLog, Datetime: r.Datetime, TodoId: p.TodoId, SchemeId: r.SchemeId})
			if err != nil {
				return
			}
			r2.ChangeLog = &cl
		}

		records = append(records, r2)
	}
	return
}
