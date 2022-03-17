package record

import (
	"flow-records/mysql"
	"strings"
	"time"
)

type PatchBody struct {
	Value    *int       `json:"value" validate:"omitempty"`
	Datetime *time.Time `json:"datetime" validate:"omitempty"`
}

func Patch(userId uint64, id uint64, p PatchBody) (r Record, notFound bool, err error) {
	r, notFound, err = Get(userId, id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Generate query
	queryStr := "UPDATE records SET"
	var queryParam []interface{}
	if p.Value != nil {
		queryStr += " value = ?,"
		queryParam = append(queryParam, p.Value)
		r.Value = *p.Value
	}
	if p.Datetime != nil {
		queryStr += " datetime = ?"
		queryParam = append(queryParam, p.Datetime.UTC())
		r.Datetime = *p.Datetime
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParam = append(queryParam, userId, id)

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
	_, err = stmtIns.Exec(queryParam...)
	if err != nil {
		return
	}

	return
}
