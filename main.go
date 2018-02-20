package main

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/DmitryDorofeev/graphcool/codegen"
	"github.com/DmitryDorofeev/graphcool/templates"
)

func main() {
	pkgDir := os.Args[1]

	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	s := codegen.Schema{
		Types: make(map[string]codegen.Type),
	}

	for _, p := range packages {
		s.Pkg = p.Name

		for _, file := range p.Files {
			for _, f := range file.Decls {
				types, err := codegen.ExtractGQLStructs(f)
				if err != nil {
					continue
				}

				for k, v := range types {
					s.Types[k] = v
				}
			}
		}
	}

	schema, err := templates.ProcessSchema(s)
	if err != nil {
		log.Println(err)
	}
	handler, err := templates.ProcessHandler(s)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile("schema.graphql", []byte(schema), os.ModePerm)
	ioutil.WriteFile(path.Join(pkgDir, "handler.go"), []byte(handler), os.ModePerm)
}
