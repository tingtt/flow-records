package record

import "time"

type PostBody struct {
	Value     int       `json:"value" validate:"required"`
	ChangeLog *string   `json:"changelog" validate:"omitempty"`
	Datetime  time.Time `json:"datetime" validate:"required"`
	TodoId    *uint64   `json:"todo_id" validate:"omitempty,gte=1"`
	SchemeId  uint64    `json:"scheme_id" validate:"required,gte=1"`
}

func Post(userId uint64, p PostBody) (r Record, err error) {
	return
}
