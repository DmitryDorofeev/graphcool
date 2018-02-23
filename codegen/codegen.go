package codegen

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/DmitryDorofeev/graphcool/codegen/templates"
	"github.com/DmitryDorofeev/graphcool/parser"
)

func Generate(in, out string) error {
	parts := strings.Split(in, "/")
	filePath := strings.Join(parts[:len(parts)-1], "/")
	if out == "" {
		out = path.Join(filePath, "generated.go")
	}

	os.Remove(out)

	pkgInfo, structs, err := parser.GetStructsInFile(in)

	if err != nil {
		return err
	}

	schema, err := templates.ProcessSchema(structs)
	if err != nil {
		return err
	}
	handler, err := templates.ProcessHandler(pkgInfo.Pkg.Name(), structs)
	if err != nil {
		return err
	}
	ioutil.WriteFile(path.Join(filePath, "schema.graphql"), []byte(schema), os.ModePerm)
	ioutil.WriteFile(out, []byte(handler), os.ModePerm)

	return exec.Command("go", "fmt", out).Run()
}
