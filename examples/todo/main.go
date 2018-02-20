package main

import (
	"net/http"

	"github.com/DmitryDorofeev/graphcool/examples/todo/models"
	graphiql "github.com/mnmtanish/go-graphiql"
)

func main() {
	http.Handle("/graphql", models.NewHandler())
	http.HandleFunc("/", graphiql.ServeGraphiQL)

	http.ListenAndServe(":8081", nil)
}
