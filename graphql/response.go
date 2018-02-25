package graphql

import (
	"encoding/json"

	"github.com/DmitryDorofeev/graphcool"
)

type Request struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type Response struct {
	Data   json.RawMessage         `json:"data,omitempty"`
	Errors []*graphcool.QueryError `json:"errors,omitempty"`
}
