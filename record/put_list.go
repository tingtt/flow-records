package record

import "flow-records/mysql"

type PutBody []postBodyWithoutTodoId

func Put(userId uint64, todoId uint64, p PutBody) (records []Record, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()

	// Delete old records and changelogs
	stmtDel1, err := db.Prepare("DELETE FROM records WHERE user_id = ? AND todo_id = ?")
	if err != nil {
		return
	}
	defer stmtDel1.Close()
	_, err = stmtDel1.Exec(userId, todoId)
	if err != nil {
		return
	}
	stmtDel2, err := db.Prepare("DELETE FROM changelogs WHERE user_id = ? AND todo_id = ?")
	if err != nil {
		return
	}
	defer stmtDel2.Close()
	_, err = stmtDel2.Exec(userId, todoId)
	if err != nil {
		return
	}

	records, err = PostMultiple(userId, MultiplePostBody{p, &todoId})

	return
}
