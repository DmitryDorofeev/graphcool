package tests

import (
	"context"

	"github.com/DmitryDorofeev/graphcool/errors"
	"github.com/DmitryDorofeev/graphcool/graphql"
)

type Query struct {
	Task    Task           `graphql:"todo:Task"`
	GetUser graphql.Getter `graphql:"getUser(name:String!):User"`
}

type Mutation struct {
}

type Task struct {
	Title       string `graphql:"title:String"`
	Description string `graphql:"description:String"`
	Done        bool   `graphql:"done:Bool"`
	User        User   `graphql:"user:User"`
}

func (t *Task) Resolve(ctx context.Context, obj interface{}, args graphql.Arguments) *errors.QueryError {
	t.Title = "test task"
	return nil
}

type User struct {
	Name    string        `graphql:"name:String"`
	Friends graphql.Field `graphql:"friends:[Friends]"`
}

func (u *User) Resolve(ctx context.Context, obj interface{}, args graphql.Arguments) *errors.QueryError {
	name, _ := args.GetString("name")
	if name == "test" {
		u.Name = name
		return nil
	}
	u.Name = "Dmitry Dorofeev"
	return nil
}

type Friends []User

func (f *Friends) Resolve(ctx context.Context, obj interface{}, args graphql.Arguments) *errors.QueryError {
	name := obj.(User).Name
	*f = append(*f, User{Name: "First Friend of " + name}, User{Name: "Second Friend of " + name})
	return nil
}
