package scheme

import "time"

type GetListQuery struct {
	ProjectId *uint64    `query:"project_id" validate:"omitempty,gte=1"`
	Start     *time.Time `query:"start" validate:"omitempty"`
	End       *time.Time `query:"end" validate:"omitempty"`
	Embed     *string    `query:"embed" validate:"omitempty,oneof=records record.changelog"`
}

func GetList(userId uint64, q GetListQuery) (schemes []Scheme, err error) {
	return
}
