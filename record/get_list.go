package record

import "time"

type GetListQuery struct {
	TodoId   *uint64    `query:"todo_id" validate:"omitempty"`
	SchemeId *uint64    `query:"scheme_id" validate:"omitempty"`
	Start    *time.Time `query:"start" validate:"omitempty"`
	End      *time.Time `query:"end" validate:"omitempty"`
	Embed    *string    `query:"embed" validate:"omitempty,oneof=changelog"`
}

func GetList(userId uint64, q GetListQuery) (records []Record, err error) {
	return
}
