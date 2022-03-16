package changelog

import (
	"flow-records/mysql"
	"time"
)

type PatchBody struct {
	Text     *string    `json:"text" validate:"required,gte=1"`
	Datetime *time.Time `json:"datetime" validate:"omitempty"`
	SchemeId *uint64    `json:"scheme_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, p PatchBody) (c ChangeLog, notFound bool, err error) {
	old, notFound, err := Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	if p.Text == nil {
		p.Text = &old.Text
	}
	if p.Datetime == nil {
		p.Datetime = &old.Datetime
	}
	if p.SchemeId == nil {
		p.SchemeId = &old.SchemeId
	}

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE changelogs SET text = ?, datetime = ?, scheme_id = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(p.Text, p.Datetime, p.SchemeId, userId, id)
	if err != nil {
		return
	}

	c.Id = id
	c.Text = *p.Text
	c.Datetime = *p.Datetime
	c.TodoId = old.TodoId
	c.SchemeId = *p.SchemeId
	return
}
