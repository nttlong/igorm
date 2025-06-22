package datasourceservice

import "dynacall"

type QueryPager struct {
	Index int `json:"index"`
	Size  int `json:"size"`
}
type QueryObject struct {
	Source  string                  `json:"dataSource"`
	Fields  map[string]bool         `json:"fields"`
	Filter  *map[string]interface{} `json:"filter"`
	OrderBy *map[string]bool        `json:"orderBy"`
	Pager   *QueryPager             `json:"pager"`
}

func (ds *DataSource) Query(queryInfo QueryObject) (interface{}, error) {
	return queryInfo.Source, nil
}

func init() {
	dynacall.RegisterCaller(&DataSource{})
}
