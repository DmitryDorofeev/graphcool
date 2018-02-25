package main

import (
	"log"
	"os"

	"github.com/DmitryDorofeev/graphcool/codegen"
)

func main() {
	file := os.Args[1]
	err := codegen.Generate(file, "")
	if err != nil {
		log.Fatal(err)
	}
}
