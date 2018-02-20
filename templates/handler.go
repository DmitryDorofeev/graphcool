package templates

import (
	"bytes"
	"html/template"
)

var handlerTmpl = `// hello
package {{.Pkg}}
{{$types := .Types}}
import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
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
		for _, s := range o.Selections {
			if ident, ok := s.(*query.Field); ok {
				{{range .Types}}
					{{if eq .Name "Query"}}
						switch(ident.Name.Name) {
							{{range .Fields}}
								case "{{.Name}}":
									{{.Name}} := &{{.Type}}{}
									err := {{.Name}}.Resolve(ctx)
									if err != nil {
										w.Write([]byte("error"))
										return
									}
									fmt.Println(ident.Name.Name)
							{{end}}
								default:
									fmt.Println("wrong!")
						}
					{{end}}
				{{end}}
			}
		}
	}

	w.Write([]byte("ba"))
}

`

func ProcessHandler(data interface{}) (string, error) {
	b := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(handlerTmpl)
	if err != nil {
		return "", err
	}

	tmpl.Execute(b, data)
	return b.String(), nil
}
