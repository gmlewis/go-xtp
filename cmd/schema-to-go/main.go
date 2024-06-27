// schema-to-go converts an XTP Extension Plugin to Go structs
// for using with XTP's APIs.
//
// Usage:
//
//	schema-to-go [-outdir dirname] *.yaml
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gmlewis/go-xtp/schema"
)

var (
	genMain = flag.Bool("main", false, "Generate main example to call functions")
	outDir  = flag.String("outdir", ".", "Directory to write the resulting Go files into.")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		if err := processFile(arg); err != nil {
			log.Fatalf("Error processing yaml file %q: %v", arg, err)
		}
	}

	log.Printf("Done.")
}

func processFile(yamlFilename string) error {
	buf, err := os.ReadFile(yamlFilename)
	if err != nil {
		return err
	}

	plugin, err := schema.ParseStr(string(buf))
	if err != nil {
		return err
	}

	structs, err := plugin.GenStructs()
	if err != nil {
		return err
	}

	if *genMain {
		mainStr, err := plugin.GenMain(structs)
		if err != nil {
			return err
		}
		return os.WriteFile("main.go", []byte(mainStr), 0644)
	}

	outFilename := strings.Replace(yamlFilename, ".yaml", ".go", 1)
	return os.WriteFile(outFilename, []byte(structs.String()), 0644)
}
