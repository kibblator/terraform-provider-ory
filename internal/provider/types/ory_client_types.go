package orytypes

type Hook struct {
	Config map[string]interface{} `json:"config,omitempty"`
	Hook   string                 `json:"hook"`
}
