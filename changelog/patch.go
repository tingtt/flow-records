package changelog

import (
	"flow-records/mysql"
	"strings"
	"time"
)

type PatchBody struct {
	Text     *string    `json:"text" validate:"required,gte=1"`
	Datetime *time.Time `json:"datetime" validate:"omitempty"`
	SchemeId *uint64    `json:"scheme_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, p PatchBody) (c ChangeLog, notFound bool, err error) {
	c, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Generate query
	queryStr := "UPDATE changelogs SET"
	var queryParams []interface{}
	if p.Text != nil {
		queryStr += " text = ?,"
		queryParams = append(queryParams, p.Text)
		c.Text = *p.Text
	}
	if p.Datetime != nil {
		queryStr += " datetime = ?,"
		queryParams = append(queryParams, p.Datetime.UTC())
		c.Datetime = *p.Datetime
	}
	if p.SchemeId != nil {
		queryStr += " scheme_id = ?"
		queryParams = append(queryParams, p.SchemeId)
		c.SchemeId = *p.SchemeId
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

	// Update row
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
	_, err = stmtIns.Exec(queryParams...)
	if err != nil {
		return
	}

	return
}
