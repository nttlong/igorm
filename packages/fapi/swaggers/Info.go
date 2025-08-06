package swaggers

// Info chứa các thông tin về API
type Info struct {
	Description string   `json:"description"`
	Title       string   `json:"title"`
	Contact     struct{} `json:"contact"`
	Version     string   `json:"version"`
}
