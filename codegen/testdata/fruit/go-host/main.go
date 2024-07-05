// fruit is a simple program that uses the XTP API and the Extism Go Host SDK to
// load and communicate with plugins defined by the XTP Extension Plugin mechanism.
//
// It requires the "XTP_TOKEN" to be set to read extensions from the XTP API.
package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/gmlewis/go-xtp/api"
	jsoniter "github.com/json-iterator/go"
)

var jsoncomp = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	appID = "app_01j1b1mek5frq9x7ymk52m7bw5"
)

func main() {
	c := api.New()

	resp, err := c.GetAppsExtensionPoints(appID)
	if err != nil {
		log.Fatalf("GetAppsExtensionPoints(%q): %v", appID, err)
	}

	allBindings := api.BindingsMap{}
	var sortedBindings []string
	for _, ep := range resp.ExtensionPoints {
		bindings, err := c.GetExtensionPointBindings(ep)
		if err != nil {
			log.Fatalf("GetExtensionPointBindings(): %v", err)
		}

		for name, binding := range bindings {
			fmt.Printf("Got binding %v: %v\n", name, binding.ID)
			allBindings[name] = binding
			sortedBindings = append(sortedBindings, name)
		}
	}
	sort.Strings(sortedBindings)

	// Now, download and call each binding.
	for _, name := range sortedBindings {
		log.Printf("Downloading plugin %v", name)
		b, _ := allBindings[name]
		if _, err := c.GetContent(b.ContentAddress); err != nil {
			log.Fatalf("GetContent(): %v", err)
		}
	}
}
