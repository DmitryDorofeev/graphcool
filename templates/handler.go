package templates

import (
	"bytes"
	"fmt"
	"go/types"
	"html/template"
	"strings"

	"github.com/DmitryDorofeev/graphcool/common"
	"github.com/DmitryDorofeev/graphcool/parser"
)

func generateLookups(fields []parser.StructField, s parser.ParsedStructs) (lookups string) {
	for _, field := range fields {
		t := field.Type()
		n := field.Name()

		var typeName string
		var fieldName string

		tag := field.Tag().Get("graphql")
		if strings.Contains(tag, ":") {
			fieldInfo := parser.ParseField(tag)
			fieldName = fieldInfo.Name
			switch fieldInfo.Type.(type) {
			case *common.TypeName:
				typeName = fieldInfo.Type.(*common.TypeName).Name
			case *common.List:
				typeName = fieldInfo.Type.(*common.List).OfType.(*common.TypeName).Name
			}
		} else {
			fieldName = tag
		}

		switch t.(type) {
		case *types.Named:
			if typeName == "" {
				typeName = t.(*types.Named).Obj().Name()
			}

			lookups += fmt.Sprintf(`
				case "%s":
					f := %sMeta{
					}

					err := f.Value.Resolve(ctx)
					if err != nil {
						return nil, err
					}

					innerField, err := f.Lookup(ctx, field.Selections)
					if err != nil {
						return nil, err
					}

					res = append(res, innerField)

			`, fieldName, typeName)
		case *types.Basic:
			if typeName == "" {
				typeName = strings.Title(t.(*types.Basic).Name())
			}

			lookups += fmt.Sprintf(`
				case "%s":
					fmt.Println("touch %s")
					f := %sMeta{
						Value: s.Value.%s,
					}

					res = append(res, f.Marshal())
			`, fieldName, fieldName, typeName, n)
		}

	}
	return
}

func generateTypes(structs parser.ParsedStructs) (types string) {
	for name, s := range structs {
		types += fmt.Sprintf(`
			type %sMeta struct {
				Meta
				Value %s
			}

			func (s *%sMeta) Lookup(ctx context.Context, selections []query.Selection) ([]byte, error) {
				fmt.Println("Looking up %s")
				res := make([][]byte, 0)
				for _, selection := range selections {
					field, ok := selection.(*query.Field)
					if !ok {
						fmt.Println("cannot cast to Field")
						continue
					}
					switch(field.Name.Name) {
						%s
						default:
							return nil, fmt.Errorf("unknown field " + field.Name.Name)
					}
				}
				return s.Marshal(ctx, selections, res)
			}

			func (s *%sMeta) Marshal(ctx context.Context, selections []query.Selection, fields [][]byte) ([]byte, error) {
				fmt.Println("Marshal %s")
				buf := bytes.Buffer{}
				buf.WriteString("{")
				for i, value := range fields {
					field, ok := selections[i].(*query.Field)
					if !ok {
						continue
					}

					if (i != 0) {
						buf.WriteString(",")
					}

					buf.WriteString("\"" + field.Name.Name + "\"")
					buf.WriteString(":")
					buf.Write(value)
				}
				buf.WriteString("}")
				return buf.Bytes(), nil
			}
		`, name, name, name, name, generateLookups(s.Fields, structs), name, name)
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
					w.Write([]byte("pzdc: " + err.Error()))
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
