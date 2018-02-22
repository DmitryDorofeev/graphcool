package models

import (
	"context"
	"fmt"

	"github.com/DmitryDorofeev/graphcool/errors"
	"github.com/DmitryDorofeev/graphcool/graphql"
)

type Query struct {
	Task Task `graphql:"todo:Task"`
}

type Task struct {
	Title       string `graphql:"title:String"`
	Description string `graphql:"description:String"`
	Done        bool   `graphql:"done:Bool"`
	User        User   `graphql:"user:User"`
}

type User struct {
	Name    string        `graphql:"name:String"`
	Friends graphql.Field `graphql:"friends:[Friends]"`
}

func (u *User) Resolve(ctx context.Context, obj interface{}) *errors.QueryError {
	fmt.Println("hello, I am user resolver: ", obj)
	u.Name = "Dmitry Dorofeev"
	return nil
}

type Friends []User

func (f *Friends) Resolve(ctx context.Context, obj interface{}) *errors.QueryError {
	name := obj.(User).Name
	*f = append(*f, User{Name: "First Friend of " + name}, User{Name: "Second Friend of " + name})
	return nil
}

func (t *Task) Resolve(ctx context.Context, obj interface{}) *errors.QueryError {
	fmt.Println("hello, I am todo resolver")
	t.Title = "Schlafen"
	t.Description = ""
	t.Done = true
	return nil
}
