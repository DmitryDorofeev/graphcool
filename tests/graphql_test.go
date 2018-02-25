package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/DmitryDorofeev/graphcool/codegen"
)

var cases = []struct {
	Request   string
	Response  string
	Variables interface{}
}{
	// Simple query
	{
		Request:  `query{todo{title}}`,
		Response: `{"data":{"todo":{"title":"test task"}}}`,
	},

	// Simple query with parameter
	{
		Request:  `query{todo{user(name:$name){name}}}`,
		Response: `{"data":{"todo":{"user":{"name":"test"}}}}`,
		Variables: map[string]string{
			"name": "test",
		},
	},

	// Query with unknown field
	{
		Request:  `query{todo{teetle}}`,
		Response: `{"errors":[{"message":"unknown field teetle"}]}`,
	},

	// Query with func
	{
		Request:  `query{getUser(name:$name){name}}`,
		Response: `{"data":{"getUser":{"name":"test"}}}`,
		Variables: map[string]string{
			"name": "test",
		},
	},

	// Query without 'query' word
	{
		Request:  `{todo{title}}`,
		Response: `{"data":{"todo":{"title":"test task"}}}`,
	},

	// Query without selection
	{
		Request:  `query{todo}`,
		Response: `{"errors":[{"message":"Objects must have selections (field Task has no selections)"}]}`,
	},

	// Query with nested selection
	{
		Request:  `query{todo{done,user{name}}}`,
		Response: `{"data":{"todo":{"done":false,"user":{"name":"Dmitry Dorofeev"}}}}`,
	},

	// Query with list selection
	{
		Request:  `query{todo{user{name,friends{name}}}}`,
		Response: `{"data":{"todo":{"user":{"name":"Dmitry Dorofeev","friends":[{"name":"First Friend of Dmitry Dorofeev"},{"name":"Second Friend of Dmitry Dorofeev"}]}}}}`,
	},

	// Named query
	{
		Request:  `query TodoUserWithFriends {todo{user{name,friends{name}}}}`,
		Response: `{"data":{"todo":{"user":{"name":"Dmitry Dorofeev","friends":[{"name":"First Friend of Dmitry Dorofeev"},{"name":"Second Friend of Dmitry Dorofeev"}]}}}}`,
	},

	// Named query with params
	{
		Request:  `query GetUserWithFriends($name: String) {getUser(name:$name){name}}`,
		Response: `{"data":{"getUser":{"name":"Dmitry Dorofeev"}}}`,
	},

	// Variables as empty string (Graphiql)
	{
		Request:   `query GetUserWithFriends($name: String) {getUser(name:$name){name}}`,
		Response:  `{"data":{"getUser":{"name":"Dmitry Dorofeev"}}}`,
		Variables: "",
	},

	// Mutation
	{
		Request:  `mutation{updateUser(name:$name){name}}`,
		Response: `{"data":{"updateUser":{"name":"Vlad_updated"}}}`,
		Variables: map[string]string{
			"name": "Vlad",
		},
	},
}

func TestGetQuery(t *testing.T) {
	ts := httptest.NewServer(NewHandler())

	for _, testCase := range cases {
		vars, _ := json.Marshal(testCase.Variables)
		resp, err := http.Get(ts.URL + fmt.Sprintf("?query=%s&variables=%s", url.QueryEscape(testCase.Request), vars))
		if err != nil {
			t.Error("error response")
		}

		data, _ := ioutil.ReadAll(resp.Body)

		if string(data) != testCase.Response {
			t.Errorf("expected %s, received %s for request %s", testCase.Response, data, testCase.Request)
		}
		resp.Body.Close()
	}
}

func TestPostQuery(t *testing.T) {
	ts := httptest.NewServer(NewHandler())
	for _, testCase := range cases {
		vars, _ := json.Marshal(testCase.Variables)
		body := []byte(fmt.Sprintf(`{"query":"%s","variables":%s}`, testCase.Request, vars))
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Error("error response")
		}

		data, _ := ioutil.ReadAll(resp.Body)

		if string(data) != testCase.Response {
			t.Errorf("expected %s, received %s for request %s", testCase.Response, data, testCase.Request)
		}
		resp.Body.Close()
	}
}

func BenchmarkPostQuery(b *testing.B) {
	ts := httptest.NewServer(NewHandler())
	for i := 0; i < b.N; i++ {
		vars, _ := json.Marshal(cases[0].Variables)
		body := []byte(fmt.Sprintf(`{"query":"%s","variables":%s}`, cases[0].Request, vars))
		resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			b.Error("error response")
		}

		data, _ := ioutil.ReadAll(resp.Body)

		if string(data) != cases[0].Response {
			b.Errorf("expected %s, received %s for request %s", cases[0].Response, data, cases[0].Request)
		}
		resp.Body.Close()
	}
}

func TestMain(m *testing.M) {
	err := codegen.Generate("./models.go", "./generated.go")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
