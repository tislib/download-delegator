package model

type Model struct {
	Name       string       `json:"name"`
	Source     string       `json:"source"`
	Examples   []Example    `json:"examples"`
	Markers    []MarkerData `json:"markers"`
	Schema     string
	ObjectType string
	UrlCheck   string
	Ref        string
}

type Example struct {
	Url string `json:"url"`
}

type MarkerData struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	ParentName string                 `json:"parentName"`
	Parameters map[string]interface{} `json:"parameters"`
}
