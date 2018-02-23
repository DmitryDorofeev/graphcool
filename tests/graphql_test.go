package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DmitryDorofeev/graphcool/codegen"
)

func TestMain(m *testing.M) {
	err := codegen.Generate("./models.go", "./generated.go")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestGetQuery(t *testing.T) {
	ts := httptest.NewServer(GQLHandler{})
	query := "query{todo{title}}"
	resp, err := http.Get(ts.URL + fmt.Sprintf("?query=%s", query))
	defer resp.Body.Close()
	if err != nil {
		t.Error("error response")
	}

	data, _ := ioutil.ReadAll(resp.Body)

	t.Log(string(data))
}

func TestPostQuery(t *testing.T) {
	ts := httptest.NewServer(GQLHandler{})
	query := "query{todo{title}}"
	body := []byte(fmt.Sprintf(`{"query":"%s"}`, query))
	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
	defer resp.Body.Close()
	if err != nil {
		t.Error("error response")
	}

	data, _ := ioutil.ReadAll(resp.Body)

	t.Log(string(data))
}
