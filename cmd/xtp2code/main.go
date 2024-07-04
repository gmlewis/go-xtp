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
//	 [-pkg=<packageName>] \
//	 [-q ] \
//	 [-appid=<id> | -yaml=<filename>] \
//	 [-force] \
//	 [-host=<filename>] \
//	 [-plugin=<filename>] \
//	 [-types=<filename>]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gmlewis/go-xtp/api"
	"github.com/gmlewis/go-xtp/codegen"
	"github.com/gmlewis/go-xtp/schema"
)

var (
	// Required:
	lang    = flag.String("lang", "", "Target language for generated code ('go' or 'mbt').")
	pkgName = flag.String("pkg", "", "Set name of generated package code when using -yaml option.")
	// Optional:
	appID     = flag.String("appid", "", "XTP App ID to generate code from.")
	force     = flag.Bool("force", false, "Force overwrite of any existing files.")
	hostDir   = flag.String("host", "", "Output dirname to generate Host SDK code.")
	pluginDir = flag.String("plugin", "", "Output dirname to generate Plugin PDK code.")
	quiet     = flag.Bool("q", false, "Do not print warnings.")
	typesDir  = flag.String("types", "", "Output dirname to generate simple types code.")
	version   = flag.Bool("v", false, "Print version and quit.")
	yamlFile  = flag.String("yaml", "", "Input schema.yaml file to generate code from. (Must also provide -pkg with this option.)")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("xtp2code version v%v\n", codegen.VERSION)
		return
	}

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

	if *pkgName == "" && *yamlFile != "" {
		log.Fatal("Must specify -pkg=<packageName> when using -yaml option")
	}

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
		p.PkgName = *pkgName

		if p.Version == "v0" {
			if !*quiet {
				log.Printf("Skipping v0 plugin")
			}
		} else {
			if err := processPlugin("", p); err != nil {
				log.Fatalf("processPlugin: %v", err)
			}
		}
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
			p.PkgName = strings.TrimSuffix(ep.Name, ".yaml")
			if *pkgName != "" {
				log.Printf("WARNING: Overriding PkgName=%q from API name %q", *pkgName, p.PkgName)
				p.PkgName = *pkgName
			}

			if p.Version == "v0" {
				if !*quiet {
					log.Printf("Skipping v0 plugin")
				}
				continue
			}

			if err := processPlugin(p.PkgName, p); err != nil {
				log.Fatalf("processPlugin: %v", err)
			}
		}
	}

	if !*quiet {
		log.Printf("Done.")
	}
}

func processPlugin(rootDir string, plugin *schema.Plugin) error {
	opts := &codegen.ClientOpts{Force: *force, Quiet: *quiet}
	c, err := codegen.New(*lang, plugin, opts)
	if err != nil {
		return err
	}

	if *typesDir != "" {
		dirName := *typesDir
		if rootDir != "" {
			dirName = filepath.Join(rootDir, dirName)
		}
		if err := c.GenTypesDir(dirName); err != nil {
			return err
		}
	}

	if *hostDir != "" {
		dirName := *hostDir
		if rootDir != "" {
			dirName = filepath.Join(rootDir, dirName)
		}
		if err := c.GenHostDir(dirName); err != nil {
			return err
		}
	}

	if *pluginDir != "" {
		dirName := *pluginDir
		if rootDir != "" {
			dirName = filepath.Join(rootDir, dirName)
		}
		if err := c.GenPluginDir(dirName); err != nil {
			return err
		}
	}

	return nil
}
