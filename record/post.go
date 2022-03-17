package record

import (
	"flow-records/changelog"
	"flow-records/mysql"
	"time"
)

type PostBody struct {
	Value     *int      `json:"value" validate:"required"`
	ChangeLog *string   `json:"changelog" validate:"omitempty,gte=1"`
	Datetime  time.Time `json:"datetime" validate:"required"`
	TodoId    *uint64   `json:"todo_id" validate:"omitempty,gte=1"`
	SchemeId  uint64    `json:"scheme_id" validate:"required,gte=1"`
}

func Post(userId uint64, p PostBody) (r Record, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO records (user_id, value, datetime, scheme_id, todo_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, p.Value, p.Datetime.UTC(), p.SchemeId, p.TodoId)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	r.Id = uint64(id)
	r.Value = *p.Value
	r.Datetime = p.Datetime
	r.SchemeId = p.SchemeId
	r.TodoId = p.TodoId

	// Post changelog
	if p.ChangeLog != nil {
		var cl changelog.ChangeLog
		cl, err = changelog.Post(userId, changelog.PostBody{Text: *p.ChangeLog, Datetime: p.Datetime, TodoId: p.TodoId, SchemeId: p.SchemeId})
		if err != nil {
			return
		}
		r.ChangeLog = &cl
	}

	return
}
