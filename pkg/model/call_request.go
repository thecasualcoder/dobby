package model

// CallRequest godoc
type CallRequest struct {
	URL    string      `json:"url"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}
