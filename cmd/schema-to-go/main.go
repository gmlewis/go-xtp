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
)

var (
	outDir = flag.String("outdir", ".", "Directory to write the resulting Go files into.")
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
}
