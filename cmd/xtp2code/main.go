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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gmlewis/go-xtp/api"
	"github.com/gmlewis/go-xtp/schema"
)

var (
	appID      = flag.String("appid", "", "XTP App ID to generate code from.")
	lang       = flag.String("lang", "", "Target language for generated code ('go' or 'mbt').")
	hostFile   = flag.String("host", "", "Output filename to generate Host SDK code.")
	pluginFile = flag.String("plugin", "", "Output filename to generate Plugin PDK code.")
	typesFile  = flag.String("types", "", "Output filename to generate simple types code.")
	yamlFile   = flag.String("yaml", "", "Input schema.yaml file to generate code from.")
)

func main() {
	flag.Parse()

	if (*appID == "" && *yamlFile == "") || (*appID != "" && *yamlFile != "") {
		log.Fatalf("Must specify either '-appid=<id>' or '-yaml=<filename>' but not both.")
	}

	if *hostFile == "" && *pluginFile == "" && *typesFile == "" {
		log.Fatalf("Must specify at least one of: -host=<filename>, -plugin=<filename>, or -types=<filename>")
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

func genPkgName(filename string) string {
	baseName := strings.Replace(filepath.Base(filename), "."+*lang, "", 1)
	return strings.ToLower(strings.Replace(baseName, "-", "_", -1))
}

func processPlugin(plugin *schema.Plugin) error {
	plugin.Lang = *lang

	custTypes, custTypesTests, err := plugin.GenCustomTypes()
	if err != nil {
		return err
	}

	if *typesFile != "" {
		pkgName := genPkgName(*typesFile)
		fullSrc := fmt.Sprintf("// Package %v represents the custom datatypes for an XTP Extension Plugin.\npackage %[1]v\n\n%v\n", pkgName, custTypes)
		if err := os.WriteFile(*typesFile, []byte(fullSrc), 0644); err != nil {
			return err
		}

		if custTypesTests != "" {
			testFilename := strings.Replace(*typesFile, "."+*lang, "_test."+*lang, 1)
			testSrc := fmt.Sprintf("package %v\n\n%v\n", pkgName, custTypesTests)
			if err := os.WriteFile(testFilename, []byte(testSrc), 0644); err != nil {
				return err
			}
		}
	}

	if *hostFile != "" {
		pkgName := genPkgName(*hostFile)
		hostSrc, err := plugin.GenHostSDK(custTypes, pkgName)
		if err != nil {
			return err
		}
		if err := os.WriteFile(*hostFile, []byte(hostSrc), 0644); err != nil {
			return err
		}
	}

	if *pluginFile != "" {
		pkgName := genPkgName(*pluginFile)
		pluginSrc, err := plugin.GenPluginPDK(custTypes, pkgName)
		if err != nil {
			return err
		}
		if err := os.WriteFile(*pluginFile, []byte(pluginSrc), 0644); err != nil {
			return err
		}
	}

	return nil
}
