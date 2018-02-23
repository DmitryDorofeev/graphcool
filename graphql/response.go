package graphql

import (
	"encoding/json"

	"github.com/DmitryDorofeev/graphcool/errors"
)

type Response struct {
	Data   json.RawMessage     `json:"data"`
	Errors []errors.QueryError `json:"errors"`
}
