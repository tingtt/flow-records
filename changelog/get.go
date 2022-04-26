package changelog

import (
	"flow-records/mysql"
)

func Get(userId uint64, id uint64) (c ChangeLog, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT text, datetime, todo_id, scheme_id FROM changelogs WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}

	c = ChangeLog{Id: id}
	err = rows.Scan(&c.Text, &c.Datetime, &c.TodoId, &c.SchemeId)
	if err != nil {
		return
	}

	return
}
