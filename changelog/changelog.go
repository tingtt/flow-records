package changelog

import "time"

type ChangeLog struct {
	Id       uint64    `json:"id"`
	Text     string    `json:"text"`
	Datetime time.Time `json:"datetime"`
	TodoId   *uint64   `json:"todo_id,omitempty"`
	SchemeId uint64    `json:"scheme_id"`
}
