package scheme

import (
	"flow-records/mysql"
	"strings"
)

type PatchBody struct {
	Name      *string `json:"name" validate:"omitempty,gte=1"`
	SumGraph  *bool   `json:"sum_graph" validate:"omitempty"`
	ProjectId *uint64 `json:"project_id" validate:"omitempty,gte=1"`
}

func Patch(userId uint64, id uint64, p PatchBody) (s Scheme, notFound bool, err error) {
	s, notFound, err = Get(userId, id, GetQuery{})
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// Generate query
	queryStr := "UPDATE schemes SET"
	var queryParams []interface{}
	if p.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, p.Name)
		s.Name = *p.Name
	}
	if p.SumGraph != nil {
		queryStr += " sum_graph = ?,"
		queryParams = append(queryParams, p.SumGraph)
		s.SumGraph = *p.SumGraph
	}
	if p.ProjectId != nil {
		queryStr += " project_id = ?"
		queryParams = append(queryParams, p.ProjectId)
		s.ProjectId = p.ProjectId
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE user_id = ? AND id = ?"
	queryParams = append(queryParams, userId, id)

	// Update row
	db, err := mysql.Open()
	if err != nil {
		return
	}
	defer db.Close()
	stmtIns, err := db.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(queryParams...)
	if err != nil {
		return
	}

	return
}
