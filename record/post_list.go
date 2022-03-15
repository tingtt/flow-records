package record

import "time"

type postBodyWithoutTodoId struct {
	Value     int       `json:"value" validate:"required"`
	ChangeLog *string   `json:"changelog" validate:"omitempty"`
	Datetime  time.Time `json:"datetime" validate:"required"`
	SchemeId  uint64    `json:"scheme_id" validate:"required,gte=1"`
}

type MultiplePostBody struct {
	Records []postBodyWithoutTodoId `json:"records" validate:"required"`
	TodoId  *uint64                 `json:"todo_id" validate:"omitempty,gte=1"`
}

func PostMultiple(userId uint64, p MultiplePostBody) (records []Record, err error) {
	return
}
