package graphql

import (
	"encoding/json"

	"github.com/DmitryDorofeev/graphcool/errors"
)

type Response struct {
	Data   json.RawMessage      `json:"data,omitempty"`
	Errors []*errors.QueryError `json:"errors,omitempty"`
}
