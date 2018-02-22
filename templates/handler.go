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
			types += generateLookup(name, generateCases(s.Fields, structs))
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
)

type GQLHandler struct {
}

type Request struct {
	Query string ` + "`json:\"query\"`" + `
}

func NewHandler() GQLHandler {
	return GQLHandler{}
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
func (h GQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("error"))
		return
	}

	req := Request{}

	json.Unmarshal(body, &req)

	res, err := query.Parse(req.Query)

	for _, o := range res.Operations {
		switch(o.Type) {
			case query.Query:
				q := QueryMeta{}
				fields, err := q.Lookup(ctx, o.Selections)
				if err != nil {
					errBytes, _ := json.Marshal(err)
					w.Write(errBytes)
					return
				}

				w.Write(fields)
				return
		}
	}

	w.Write([]byte("ba"))
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
