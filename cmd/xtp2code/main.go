// xtp2code converts an XTP Extension Plugin to Go or MoonBit source code
// for use with XTP's APIs. It can generate simple custom datatypes and/or Host SDK code
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
//	     [-types=<filename>] \
//	     [-yaml=<filename>]
package main

import (
	"flag"
	"log"
	"os"

	"github.com/gmlewis/go-xtp/api"
	"github.com/gmlewis/go-xtp/codegen"
	"github.com/gmlewis/go-xtp/schema"
)

var (
	appID         = flag.String("appid", "", "XTP App ID to generate code from.")
	lang          = flag.String("lang", "", "Target language for generated code ('go' or 'mbt').")
	hostDirFile   = flag.String("host", "", "Output dirname or filename to generate Host SDK code.")
	pluginDirFile = flag.String("plugin", "", "Output dirname or filename to generate Plugin PDK code.")
	typesDirFile  = flag.String("types", "", "Output dirname or filename to generate simple types code.")
	yamlFile      = flag.String("yaml", "", "Input schema.yaml file to generate code from.")
)

func main() {
	flag.Parse()

	if (*appID == "" && *yamlFile == "") || (*appID != "" && *yamlFile != "") {
		log.Fatalf("Must specify either '-appid=<id>' or '-yaml=<filename>' but not both.")
	}

	if *hostDirFile == "" && *pluginDirFile == "" && *typesDirFile == "" {
		log.Fatalf("Must specify at least one of: -host=<dirname>, -plugin=<dirname>, or -types=<dirname/filename>")
	}

	switch *lang {
	case "go", "Go":
		*lang = "go"
	case "mbt", "moon", "moonbit", "MoonBit":
		*lang = "mbt"
	default:
		log.Fatalf("Must specify either -lang=go or -lang=mbt")
	}

	var plugins []*schema.Plugin

	switch {
	case *yamlFile != "":
		buf, err := os.ReadFile(*yamlFile)
		if err != nil {
			log.Fatal(err)
		}
		p, err := schema.ParseStr(string(buf))
		if err != nil {
			log.Fatalf("schema.Parse: %v", err)
		}
		plugins = append(plugins, p)
	case *appID != "":
		c := api.New()
		resp, err := c.GetAppsExtensionPoints(*appID)
		if err != nil {
			log.Fatal(err)
		}
		for _, ep := range resp.ExtensionPoints {
			p, err := schema.ParseStr(ep.SchemaYaml)
			if err != nil {
				log.Fatalf("schema.Parse: %v", err)
			}
			plugins = append(plugins, p)
		}
	}

	for _, plugin := range plugins {
		if err := processPlugin(plugin); err != nil {
			log.Fatalf("processPlugin: %v", err)
		}
	}

	log.Printf("Done.")
}

func processPlugin(plugin *schema.Plugin) error {
	c, err := codegen.New(*lang, plugin)
	if err != nil {
		return err
	}

	if *typesDirFile != "" {
		if err := c.GenTypesDirFiles(*typesDirFile); err != nil {
			return err
		}
	}

	if *hostDirFile != "" {
		if err := c.GenHostDirFiles(*hostDirFile); err != nil {
			return err
		}
	}

	if *pluginDirFile != "" {
		if err := c.GenPluginDirFiles(*pluginDirFile); err != nil {
			return err
		}
	}

	return nil
}
