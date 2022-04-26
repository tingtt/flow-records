package record

import (
	"flow-records/mysql"
)

func Get(userId uint64, id uint64) (r Record, notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT value, datetime, scheme_id, todo_id FROM records WHERE user_id = ? AND id = ?")
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
	err = rows.Scan(&r.Value, &r.Datetime, &r.SchemeId, &r.TodoId)
	if err != nil {
		return
	}

	r.Id = id
	return
}
