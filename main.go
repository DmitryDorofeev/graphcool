package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/DmitryDorofeev/graphcool/parser"
	"github.com/DmitryDorofeev/graphcool/templates"
)

func main() {
	file := os.Args[1]
	parts := strings.Split(file, "/")
	filePath := strings.Join(parts[:len(parts)-1], "/")
	resultFile := path.Join(filePath, "handler.go")

	os.Remove(resultFile)

	pkgInfo, structs, err := parser.GetStructsInFile(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", structs)

	schema, err := templates.ProcessSchema(structs)
	if err != nil {
		log.Println(err)
	}
	handler, err := templates.ProcessHandler(pkgInfo.Pkg.Name(), structs)
	if err != nil {
		log.Println(err)
	}
	ioutil.WriteFile("schema.graphql", []byte(schema), os.ModePerm)
	ioutil.WriteFile(resultFile, []byte(handler), os.ModePerm)

	exec.Command("go", "fmt", resultFile).Run()
}
