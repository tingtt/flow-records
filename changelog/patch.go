package changelog

type PatchBody struct {
	Name      *string `json:"name" validate:"omitempty,gte=1"`
	SumGraph  *bool   `json:"sum_graph" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, p PatchBody) (c ChangeLog, notFound bool, err error) {
	return
}
