package record

import "flow-records/mysql"

func Delete(userId uint64, id uint64) (notFound bool, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	stmtIns, err := db.Prepare("DELETE FROM records WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(userId, id)
	if err != nil {
		return
	}

	// Check affected row count
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return
	}
	if affectedRowCount == 0 {
		// Not found
		notFound = true
		return
	}
	return
}
