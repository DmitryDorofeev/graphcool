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
	"bytes"
	"strconv"
	"fmt"
	"github.com/DmitryDorofeev/graphcool"
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

	handlerTmpl += fmt.Sprintf(`
		func (h GraphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			c := &graphcool.Context{
				Request: r,
				Keys: nil,
			}

			parseError := graphcool.Errorf("Cannot parse request")
			req := graphql.Request{}

			var errs []*graphcool.QueryError
			var vars map[string]interface{}
			switch r.Method {
				case http.MethodGet:
					req.Query = r.URL.Query().Get("query")
					v := r.URL.Query().Get("variables")
					if v == "" {
						break
					}
					json.Unmarshal([]byte(v), &vars)
				case http.MethodPost:
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						errBytes, _ := json.Marshal(parseError)
						w.Write(errBytes)
						return
					}

					err = json.Unmarshal(body, &req)
					if err != nil {
						errBytes, _ := json.Marshal(graphcool.Errorf(err.Error()))
						w.Write(errBytes)
						return
					}

					vars, _ = req.Variables.(map[string]interface{})

				default:
					w.WriteHeader(http.StatusMethodNotAllowed)
			}
			res, err := query.Parse(req.Query)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("{\"errors\": [{\"message\":\"%%s\"}]}", err)))
				return
			}

			for _, o := range res.Operations {
				switch(o.Type) {
					case query.Query:
						%s
					case query.Mutation:
						%s
				}
			}
		}
	`, handleQuery(data), handleMutation(data))

	b := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(handlerTmpl)
	if err != nil {
		return "", err
	}

	tmpl.Execute(b, data)
	return b.String(), nil
}

func handleQuery(structs parser.ParsedStructs) string {
	s, ok := structs["Query"]
	if !ok {
		return `
			errs = []*graphcool.QueryError{
				{
					Message: "query resolvers are not present",
				},
			}

			resp, _ := json.Marshal(graphql.Response{
				Errors: errs,
			})

			w.Write(resp)
			return
		`
	}

	var resolve string
	var hasResolve bool
	for _, method := range s.Methods {
		if method.Name == "Resolve" {
			hasResolve = true
		}
	}
	if hasResolve {
		resolve = "q.Value.Resolve(c, nil, graphql.Arguments{})"
	}

	return fmt.Sprintf(`
		q := QueryMeta{}
		%s
		fields, err := q.Lookup(c, o.Selections, vars)
		if err != nil {
			errs = []*graphcool.QueryError{
				err,
			}
		}
		resp, _ := json.Marshal(graphql.Response{
			Data: fields,
			Errors: errs,
		})

		w.Write(resp)
		return
	`, resolve)
}

func handleMutation(structs parser.ParsedStructs) string {
	s, ok := structs["Mutation"]
	if !ok {
		return `
			errs = []*graphcool.QueryError{
				{
					Message: "mutation resolvers are not present",
				},
			}

			resp, _ := json.Marshal(graphql.Response{
				Errors: errs,
			})

			w.Write(resp)
			return
		`
	}

	var resolve string
	var hasResolve bool
	for _, method := range s.Methods {
		if method.Name == "Resolve" {
			hasResolve = true
		}
	}
	if hasResolve {
		resolve = "m.Value.Resolve(c, nil, graphql.Arguments{})"
	}

	return fmt.Sprintf(`
		m := MutationMeta{}
		%s
		fields, err := m.Lookup(c, o.Selections, vars)
		if err != nil {
			errs = []*graphcool.QueryError{
				err,
			}
		}
		resp, _ := json.Marshal(graphql.Response{
			Data: fields,
			Errors: errs,
		})

		w.Write(resp)
		return
	`, resolve)
}
