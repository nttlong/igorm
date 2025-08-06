package swaggers

type Operation struct {
	Consumes   []string              `json:"consumes"`
	Produces   []string              `json:"produces"`
	Tags       []string              `json:"tags"`
	Summary    string                `json:"summary"`
	Parameters []Parameter           `json:"parameters"`
	Responses  map[string]Response   `json:"responses"`
	Security   []map[string][]string `json:"security"`
}
