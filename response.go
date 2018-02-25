package graphcool

import (
	"encoding/json"
)

type Request struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type Response struct {
	Data   json.RawMessage `json:"data,omitempty"`
	Errors []*QueryError   `json:"errors,omitempty"`
}
