package changelog

type ChangeLog struct {
	Id       uint64  `json:"id"`
	Text     string  `json:"text"`
	TodoId   *uint64 `json:"todo_id,omitempty"`
	SchemeId uint64  `json:"scheme_id"`
}
