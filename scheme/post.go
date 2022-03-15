package scheme

type PostBody struct {
	Name      string  `json:"name"`
	SumGraph  bool    `json:"sum_graph"`
	ProjectId *uint64 `json:"project_id,omitempty"`
}

func Post(userId uint64, p PostBody) (s Scheme, err error) {
	return
}
