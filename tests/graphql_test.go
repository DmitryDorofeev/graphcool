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
	Variables map[string]string
}{
	{
		Request:  `query{todo{title}}`,
		Response: `{"data":{"todo":{"title":"test task"}}}`,
	},
	{
		Request:  `query{todo{title(name:$name)}}`,
		Response: `{"data":{"todo":{"title":"test task"}}}`,
		Variables: map[string]string{
			"name": "test",
		},
	},
	{
		Request:  `query{todo{teetle}}`,
		Response: `{"errors":[{"message":"unknown field teetle"}]}`,
	},
	{
		Request:  `query{getUser(name:$name){name}}`,
		Response: `{"data":{"getUser":{"name":"test"}}}`,
		Variables: map[string]string{
			"name": "test",
		},
	},
	{
		Request:  `{todo{title}}`,
		Response: `{"data":{"todo":{"title":"test task"}}}`,
	},
	{
		Request:  `query{todo}`,
		Response: `{"errors":[{"message":"Objects must have selections (field Task has no selections)"}]}`,
	},
	{
		Request:  `query{todo{user{name}}}`,
		Response: `{"data":{"todo":{"user":{"name":"Dmitry Dorofeev"}}}}`,
	},
	{
		Request:  `query{todo{user{name,friends{name}}}}`,
		Response: `{"data":{"todo":{"user":{"name":"Dmitry Dorofeev","friends":[{"name":"First Friend of Dmitry Dorofeev"},{"name":"Second Friend of Dmitry Dorofeev"}]}}}}`,
	},
}

func TestGetQuery(t *testing.T) {
	ts := httptest.NewServer(GQLHandler{})

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
	ts := httptest.NewServer(GQLHandler{})
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

func TestMain(m *testing.M) {
	err := codegen.Generate("./models.go", "./generated.go")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
