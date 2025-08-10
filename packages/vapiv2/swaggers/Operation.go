package swaggers

type Operation struct {
	IgnoreBasePath bool                  `json:"x-ignoreBasePath"`
	Consumes       []string              `json:"consumes"`
	Produces       []string              `json:"produces"`
	Tags           []string              `json:"tags"`
	Summary        string                `json:"summary"`
	Parameters     []Parameter           `json:"parameters"`
	Responses      map[string]Response   `json:"responses"`
	Security       []map[string][]string `json:"security"`
}
