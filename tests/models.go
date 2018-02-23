package tests

import (
	"context"

	"github.com/DmitryDorofeev/graphcool/errors"
)

type Query struct {
	Task Task `graphql:"todo:Task"`
}

type Task struct {
	Title       string `graphql:"title:String"`
	Description string `graphql:"description:String"`
	Done        bool   `graphql:"done:Bool"`
}

func (t *Task) Resolve(ctx context.Context, obj interface{}) *errors.QueryError {
	t.Title = "test task"
	return nil
}
