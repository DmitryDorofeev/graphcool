package templates

import (
	"bytes"
	"html/template"
)

var schemaTmpl = `{{range $name, $struct := .}}
type {{$name}} {
	{{range $struct.Fields}}{{.Tag.Get "graphql"}}: {{.Type.Obj.Name}}{{end}}
}
{{end}}
`

func ProcessSchema(data interface{}) (string, error) {
	b := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(schemaTmpl)
	if err != nil {
		return "", err
	}

	tmpl.Execute(b, data)
	return b.String(), nil
}
