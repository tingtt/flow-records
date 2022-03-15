package changelog

type PostBody struct {
	Name      string  `json:"name" validate:"required,gte=1"`
	SumGraph  bool    `json:"sum_graph" validate:"required"`
	ProjectId *uint64 `json:"project_id,omitempty" validate:"omitempty,gte=1"`
}

func Post(userId uint64, p PostBody) (c ChangeLog, err error) {
	return
}
