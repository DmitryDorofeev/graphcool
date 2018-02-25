package tests

import (
	"github.com/DmitryDorofeev/graphcool"
)

type Query struct {
	Task    Task             `graphql:"todo:Task"`
	GetUser graphcool.Getter `graphql:"getUser(name:String!):User"`
}

type Mutation struct {
}

// updateUser(name:String!):User
func (m *Mutation) UpdateUser(c *graphcool.Context, obj interface{}, args graphcool.Arguments) (User, *graphcool.QueryError) {
	name, _ := args.GetString("name")
	return User{
		Name: name + "_updated",
	}, nil
}

type Task struct {
	Title       string `graphql:"title:String"`
	Description string `graphql:"description:String"`
	Done        bool   `graphql:"done:Bool"`
	User        User   `graphql:"user:User"`
}

func (t *Task) Resolve(c *graphcool.Context, obj interface{}, args graphcool.Arguments) *graphcool.QueryError {
	t.Title = "test task"
	return nil
}

type User struct {
	Name    string          `graphql:"name:String"`
	Friends graphcool.Field `graphql:"friends:[Friends]"`
}

func (u *User) Resolve(c *graphcool.Context, obj interface{}, args graphcool.Arguments) *graphcool.QueryError {
	name, err := args.GetString("name")
	if err != nil {
		u.Name = "Dmitry Dorofeev"
		return nil
	}

	u.Name = name
	return nil
}

type Friends []User

func (f *Friends) Resolve(c *graphcool.Context, obj interface{}, args graphcool.Arguments) *graphcool.QueryError {
	name := obj.(User).Name
	*f = append(*f, User{Name: "First Friend of " + name}, User{Name: "Second Friend of " + name})
	return nil
}
