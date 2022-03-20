package scheme

import (
	"encoding/json"
	"flow-records/mysql"
	"strings"
)

type PatchBody struct {
	Name      *string             `json:"name" validate:"omitempty,gte=1"`
	SumGraph  *bool               `json:"sum_graph" validate:"omitempty"`
	ProjectId PatchNullJSONUint64 `json:"project_id" validate:"dive"`
}
type PatchNullJSONUint64 struct {
	UInt64 **uint64 `validate:"omitempty,gte=1"`
}

func (p *PatchNullJSONUint64) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	var valueP *uint64 = nil
	if string(data) == "null" {
		// key exists and value is null
		p.UInt64 = &valueP
		return nil
	}

	var tmp uint64
	tmpP := &tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		// invalid value type
		return err
	}
	// valid value
	p.UInt64 = &tmpP
	return nil
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
	if p.ProjectId.UInt64 != nil {
		if *p.ProjectId.UInt64 != nil {
			queryStr += " project_id = ?"
			queryParams = append(queryParams, **p.ProjectId.UInt64)
			s.ProjectId = *p.ProjectId.UInt64
		} else {
			queryStr += " project_id = ?"
			queryParams = append(queryParams, nil)
			s.ProjectId = nil
		}
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
