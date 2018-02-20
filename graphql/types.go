package graphql

import (
	"context"
	"log"

	"github.com/DmitryDorofeev/graphcool/errors"
)

type String string
type Int string

func (s String) Resolve(ctx context.Context, rootObject interface{}) (String, *errors.QueryError) {
	log.Println("String resolved")
	return s, nil
}

func (i Int) Resolve(ctx context.Context, rootObject interface{}) (Int, *errors.QueryError) {
	log.Println("Int resolved")
	return i, nil
}
