package templates

import (
	"bytes"
	"html/template"
)

var schemaTmpl = `
schema {
	query: Query
}
{{range .Types}}
type {{.Name}} {
	{{range .Fields}}{{.Name}}: {{.Type}}{{end}}
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
