package templates

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/DmitryDorofeev/graphcool/parser"
)

func generateTypes(structs parser.ParsedStructs) (types string) {
	for name, s := range structs {
		types += fmt.Sprintf(`
				type %sMeta struct {
					Meta
					Value %s
				}
			`, name, name)

		if name != s.TypeName {
			types += generateListLookup(name, s.TypeName)
			types += generateListMarshaler(name)
		} else {
			types += generateLookup(name, rangeFields(s, structs))
			types += generateObjectMarshaler(name)
		}
	}
	return
}

func ProcessHandler(pkg string, data parser.ParsedStructs) (string, error) {

	var handlerTmpl = "// auto generated file \n"
	handlerTmpl += fmt.Sprintf(`package %s`, pkg)
	handlerTmpl += `
import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"context"
	"bytes"
	"strconv"
	"fmt"
	"github.com/DmitryDorofeev/graphcool/errors"
	"github.com/DmitryDorofeev/graphcool/query"
	"github.com/DmitryDorofeev/graphcool/graphql"
)

type GraphqlHandler struct {
}

func NewHandler() GraphqlHandler {
	return GraphqlHandler{}
}

type Meta struct {
	Nullable bool
	List bool
}

type StringMeta struct {
	Meta
	Value string
}

func (s *StringMeta) Marshal() ([]byte) {
	return []byte("\""+s.Value+"\"")
}

type IntMeta struct {
	Meta
	Value int
}

func (s *IntMeta) Marshal() ([]byte) {
	return []byte(strconv.Itoa(s.Value))
}

type BoolMeta struct {
	Meta
	Value bool
}

func (s *BoolMeta) Marshal() ([]byte) {
	return []byte(strconv.FormatBool(s.Value))
}
`

	handlerTmpl += generateTypes(data)

	handlerTmpl += `
func (h GraphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parseError := errors.Errorf("Cannot parse request")
	req := graphql.Request{}

	var errs []*errors.QueryError
	switch r.Method {
		case http.MethodGet:
			req.Query = r.URL.Query().Get("query")
			vars := r.URL.Query().Get("variables")
			if vars == "" {
				break
			}
			json.Unmarshal([]byte(vars), &req.Variables)
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errBytes, _ := json.Marshal(parseError)
				w.Write(errBytes)
				return
			}

			err = json.Unmarshal(body, &req)
			if err != nil {
				errBytes, _ := json.Marshal(parseError)
				w.Write(errBytes)
				return
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
	res, err := query.Parse(req.Query)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{\"errors\": [{\"message\":\"%s\"}]}", err)))
		return
	}

	for _, o := range res.Operations {
		switch(o.Type) {
			case query.Query:
				q := QueryMeta{}
				fields, err := q.Lookup(ctx, o.Selections, req.Variables)
				if err != nil {
					errs = []*errors.QueryError{
						err,
					}
				}
				resp, _ := json.Marshal(graphql.Response{
					Data: fields,
					Errors: errs,
				})

				w.Write(resp)
				return
			case query.Mutation:
				m := MutationMeta{}
				fields, err := m.Lookup(ctx, o.Selections, req.Variables)
				if err != nil {
					errs = []*errors.QueryError{
						err,
					}
				}
				resp, _ := json.Marshal(graphql.Response{
					Data: fields,
					Errors: errs,
				})

				w.Write(resp)
				return
		}
	}
}

`

	b := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(handlerTmpl)
	if err != nil {
		return "", err
	}

	tmpl.Execute(b, data)
	return b.String(), nil
}
