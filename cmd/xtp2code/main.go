// xtp2code converts an XTP Extension Plugin to Go or MoonBit source code
// for use with XTP's APIs. It can generate simple custom datatypes and/or Host SDK code
// and/or Plugin PDK code. For input, it can process either a schema.yaml file
// or it can query the XTP API directly for a given app ID (for the authenticated
// user) and process all extension plugin definitions.
//
// Note that if `-appid` is provided, the "XTP_TOKEN" environment variable must
// be set for the logged-in XTP user (`xtp auto login`).
//
// Usage:
//
//	xtp2code \
//	 -lang=[go|mbt] \
//	 -pkg=<packageName> \
//	 [-q ] \
//	 [-appid=<id> | -yaml=<filename>] \
//	 [-force] \
//	 [-host=<filename>] \
//	 [-plugin=<filename>] \
//	 [-types=<filename>]
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
	// Required:
	lang    = flag.String("lang", "", "Target language for generated code ('go' or 'mbt').")
	pkgName = flag.String("pkg", "", "Set name of generated package code.")
	// Optional:
	appID     = flag.String("appid", "", "XTP App ID to generate code from.")
	force     = flag.Bool("force", false, "Force overwrite of any existing files.")
	hostDir   = flag.String("host", "", "Output dirname to generate Host SDK code.")
	pluginDir = flag.String("plugin", "", "Output dirname to generate Plugin PDK code.")
	quiet     = flag.Bool("q", false, "Do not print warnings.")
	typesDir  = flag.String("types", "", "Output dirname to generate simple types code.")
	yamlFile  = flag.String("yaml", "", "Input schema.yaml file to generate code from.")
)

func main() {
	flag.Parse()

	if (*appID == "" && *yamlFile == "") || (*appID != "" && *yamlFile != "") {
		log.Fatal("Must specify either '-appid=<id>' or '-yaml=<filename>' but not both.")
	}

	if *hostDir == "" && *pluginDir == "" && *typesDir == "" {
		log.Fatal("Must specify at least one of: -host=<dirname>, -plugin=<dirname>, or -types=<dirname>")
	}

	switch *lang {
	case "go", "Go":
		*lang = "go"
	case "mbt", "moon", "moonbit", "MoonBit":
		*lang = "mbt"
	default:
		log.Fatal("Must specify either -lang=go or -lang=mbt")
	}

	if *pkgName == "" {
		log.Fatal("Must specify -pkg=<packageName>")
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

	if !*quiet {
		log.Printf("Done.")
	}
}

func processPlugin(plugin *schema.Plugin) error {
	opts := &codegen.ClientOpts{Force: *force, Quiet: *quiet}
	c, err := codegen.New(*lang, *pkgName, plugin, opts)
	if err != nil {
		return err
	}

	if *typesDir != "" {
		if err := c.GenTypesDir(*typesDir); err != nil {
			return err
		}
	}

	if *hostDir != "" {
		if err := c.GenHostDir(*hostDir); err != nil {
			return err
		}
	}

	if *pluginDir != "" {
		if err := c.GenPluginDir(*pluginDir); err != nil {
			return err
		}
	}

	return nil
}
