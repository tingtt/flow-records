package record

import "time"

type PatchBody struct {
	Value    *int       `json:"value" validate:"omitempty"`
	Datetime *time.Time `json:"datetime" validate:"omitempty"`
}

func Patch(userId uint64, id uint64, p PatchBody) (r Record, notFound bool, err error) {
	return
}
