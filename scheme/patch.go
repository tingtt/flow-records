package scheme

import "flow-records/mysql"

type PatchBody struct {
	Name      *string `json:"name" validate:"omitempty,gte=1"`
	SumGraph  *bool   `json:"sum_graph" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, p PatchBody) (s Scheme, notFound bool, err error) {
	old, notFound, err := Get(userId, id, GetQuery{})
	if err != nil {
		return
	}
	if notFound {
		return
	}

	if p.Name == nil {
		p.Name = &old.Name
	}
	if p.SumGraph == nil {
		p.SumGraph = &old.SumGraph
	}
	if p.ProjectId == nil {
		p.ProjectId = old.ProjectId
	}

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("UPDATE schemes SET name = ?, sum_graph = ?, project_id = ? WHERE user_id = ? AND id = ?")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(p.Name, p.SumGraph, p.ProjectId, userId, id)
	if err != nil {
		return
	}

	s.Id = id
	s.Name = *p.Name
	s.SumGraph = *p.SumGraph
	s.ProjectId = p.ProjectId
	return
}
