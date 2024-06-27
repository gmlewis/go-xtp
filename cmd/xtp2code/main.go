// xtp2code converts an XTP Extension Plugin to Go or MoonBit source code
// for use with XTP's APIs. It can generate simple structs and/or Host SDK code
// and/or Plugin PDK code. For input, it can process either a schema.yaml file
// or it can query the XTP API directly for a given app ID (for the authenticated
// user) and process all extensions.
//
// Note that if `-appid` is provided, the "XTP_TOKEN" environment variable must
// be set for the logged-in XTP user (`xtp auto login`).
//
// Usage:
//
//		xtp2code \
//	     -lang=[go|mbt] \
//	     [-appid=<id>] \
//	     [-host=<filename>] \
//	     [-plugin=<filename>] \
//	     [-structs=<filename>] \
//	     [-yaml=<filename>]
package main

import (
	"flag"
	"log"
	"os"

	"github.com/gmlewis/go-xtp/api"
)

var (
	appID      = flag.String("appid", "", "XTP App ID to generate code from.")
	lang       = flag.String("lang", "", "Target language for generated code ('go' or 'mbt').")
	hostFile   = flag.String("host", "", "Output filename to generate Host SDK code.")
	pluginFile = flag.String("plugin", "", "Output filename to generate Plugin PDK code.")
	structFile = flag.String("structs", "", "Output filename to generate simple structs code.")
	yamlFile   = flag.String("yaml", "", "Input schema.yaml file to generate code from.")
)

func main() {
	flag.Parse()

	if (*appID == "" && *yamlFile == "") || (*appID != "" && *yamlFile != "") {
		log.Fatalf("Must specify either '-appid=<id>' or '-yaml=<filename>' but not both.")
	}

	if *hostFile == "" && *pluginFile == "" && *structFile == "" {
		log.Fatalf("Must specify at least one of: -host=<filename>, -plugin=<filename>, or -structs=<filename>")
	}

	switch *lang {
	case "go", "Go":
		*lang = "go"
	case "mbt", "moon", "moonbit", "MoonBit":
		*lang = "mbt"
	default:
		log.Fatalf("Must specify either -lang=go or -lang=mbt")
	}

	var schemaYaml string
	switch {
	case *yamlFile != "":
		buf, err := os.ReadFile(*yamlFile)
		if err != nil {
			log.Fatal(err)
		}
		schemaYaml = string(buf)
	case *appID != "":
		c := api.New()
		resp, err := c.GetAppsExtensionPoints(*appID)
		if err != nil {
			log.Fatal(err)
		}
		schemaYaml = resp.SchemaYaml
	}

	if err := processFile(schemaYaml); err != nil {
		log.Fatalf("Error processing yaml file %q: %v", schemaYaml, err)
	}

	log.Printf("Done.")
}

func processFile(yamlFilename string) error {
	return nil
}
