package changelog

import (
	"flow-records/mysql"
	"time"
)

type PostBody struct {
	Text     string    `json:"text" validate:"required,gte=1"`
	Datetime time.Time `json:"datetime" validate:"required"`
	TodoId   *uint64   `json:"todo_id,omitempty" validate:"omitempty,gte=1"`
	SchemeId uint64    `json:"scheme_id" validate:"required,gte=1"`
}

func Post(userId uint64, p PostBody) (c ChangeLog, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO changelogs (user_id, text, datetime, scheme_id, todo_id) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, p.Text, p.Datetime.UTC(), p.SchemeId, p.TodoId)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	c.Id = uint64(id)
	c.Text = p.Text
	c.Datetime = p.Datetime.UTC()
	c.SchemeId = p.SchemeId
	if p.TodoId != nil {
		c.TodoId = p.TodoId
	}
	return
}
