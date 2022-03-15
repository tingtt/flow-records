package scheme

import "time"

type GetQuery struct {
	Start *time.Time `query:"start" validate:"omitempty"`
	End   *time.Time `query:"end" validate:"omitempty"`
	Embed *string    `query:"embed" validate:"omitempty,oneof=records record.changelog"`
}

func Get(userId uint64, id uint64, q GetQuery) (s Scheme, notFound bool, err error) {
	return
}
