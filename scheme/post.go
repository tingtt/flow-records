package scheme

import "flow-records/mysql"

type PostBody struct {
	Name      string  `json:"name" validate:"required,gte=1"`
	SumGraph  bool    `json:"sum_graph" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id,omitempty" validate:"omitempty,gte=1"`
}

func Post(userId uint64, p PostBody) (s Scheme, err error) {
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare("INSERT INTO schemes (user_id, name, sum_graph, project_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(userId, p.Name, p.SumGraph, p.ProjectId)
	if err != nil {
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	s.Id = uint64(id)
	s.Name = p.Name
	s.SumGraph = p.SumGraph
	s.ProjectId = p.ProjectId
	return
}
