package graphql

import (
	"encoding/json"

	"github.com/DmitryDorofeev/graphcool/errors"
)

type Request struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type Response struct {
	Data   json.RawMessage      `json:"data,omitempty"`
	Errors []*errors.QueryError `json:"errors,omitempty"`
}
