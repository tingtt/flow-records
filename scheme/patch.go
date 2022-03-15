package scheme

type PatchBody struct {
	Name      *string `json:"name" validate:"omitempty,gte=1"`
	SumGraph  *bool   `json:"sum_graph" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, p PatchBody) (s Scheme, notFound bool, err error) {
	return
}
