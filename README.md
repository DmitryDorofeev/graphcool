# Golang GraphQL codegen library [![Build Status](https://travis-ci.org/DmitryDorofeev/graphcool.svg?branch=master)](https://travis-ci.org/DmitryDorofeev/graphcool) [![Maintainability](https://api.codeclimate.com/v1/badges/c890cd27321257d0c116/maintainability)](https://codeclimate.com/github/DmitryDorofeev/graphcool/maintainability)

## Как тебе такое, Илон Маск?

#### Usage

First things first, install `graphcool`

`go get github.com/DmitryDorofeev/graphcool/...`

Define `Query` and  `Mutation` structs:

```go
type Query struct {
	Task    Task           `graphql:"todo:Task"`
	GetUser graphcool.Getter `graphql:"getUser(name:String!):User"`
}

type Mutation struct {
}

type Task struct {
	Title       string `graphql:"title:String"`
	Description string `graphql:"description:String"`
	Done        bool   `graphql:"done:Bool"`
	User        User   `graphql:"user:User"`
}

type User struct {
	Name string `graphql:"name:String"`
}

```

Put your field schema into `graphql` tag

Define Resolvers of your data structs:

```go
func (u *User) Resolve(c *graphcool.Context, obj interface{}, args graphcool.Arguments) *errors.QueryError {
	u.Name = "Dmitry Dorofeev"
	return nil
}

func (t *Task) Resolve(c *graphcool.Context, obj interface{}, args graphcool.Arguments) *errors.QueryError {
	t.Title = "Awesome task"
	return nil
}
```

Define methods for `Mutation`

```go
// updateUser(name:String!):User
func (m *Mutation) UpdateUser(c *graphcool.Context, obj interface{}, args graphcool.Arguments) (User, *errors.QueryError) {
	name, _ := args.GetString("name")
	return User{
		Name: name,
	}, nil
}
```
⚠️ WARNING: comment with a method schema is necessary!

Generate code for your data `graphcool ./filewithstructs.go`

Use `NewHandler()` method

```go
http.Handle("/graphql", models.NewHandler())
```

### Roadmap
- [x] Mutations
- [x] Queries with params
- [x] Custom context
- [ ] Pass query name to resolvers
- [ ] More docs :)
- [ ] Generate file with graphql schema
- [ ] Logging
- [ ] Schema queries
- [ ] Nullable types
- [ ] Directives
- [ ] Fragments
- [ ] Subscriptions
