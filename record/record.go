package record

import (
	"flow-records/changelog"
	"time"
)

type Record struct {
	Id        uint64               `json:"id"`
	Value     int                  `json:"value"`
	ChangeLog *changelog.ChangeLog `json:"changelog,omitempty"`
	Datetime  time.Time            `json:"datetime"`
	TodoId    *uint64              `json:"todo_id,omitempty"`
	SchemeId  uint64               `json:"scheme_id"`
}
