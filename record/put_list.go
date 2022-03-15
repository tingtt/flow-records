package record

type PutBody []postBodyWithoutTodoId

func Put(userId uint64, todoId uint64, p PutBody) (records []Record, err error) {
	return
}
