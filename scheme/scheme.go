package scheme

import "flow-records/record"

type Scheme struct {
	Id        uint64          `json:"id"`
	Name      string          `json:"name"`
	SumGraph  bool            `json:"sum_graph"`
	ProjectId *uint64         `json:"project_id,omitempty"`
	Records   []record.Record `json:"records,omitempty"`
}
