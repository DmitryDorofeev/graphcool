package models

import (
	"context"
	"fmt"

	"github.com/DmitryDorofeev/graphcool/graphql"
)

type Schema struct {
	Query Query `graphql:"query"`
}

type Query struct {
	Todo Todo `graphql:"todo"`
}

type Mutation struct {
}

type Todo struct {
	Title graphql.String `graphql:"title"`
}

func (t *Todo) Resolve(ctx context.Context) error {
	fmt.Println("hello, I am todo resolver")
	t.Title = "ba"
	return nil
}
