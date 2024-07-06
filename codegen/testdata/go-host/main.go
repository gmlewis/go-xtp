// fruit is a simple program that uses the XTP API and the Extism Go Host SDK to
// load and communicate with plugins defined by the XTP Extension Plugin mechanism.
//
// It requires the "XTP_TOKEN" to be set to read extensions from the XTP API.
package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	extism "github.com/extism/go-sdk"
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

	ctx := context.Background()

	// Now, download and call each plugin function.
	for _, name := range sortedBindings {
		binding, ok := allBindings[name]
		if !ok {
			log.Fatalf("programming error - lookup of name %v failed", name)
		}
		extID := strings.Split(binding.ID, "/")[0]

		log.Printf("Calling plugin: %v (extension point ID: %v)", name, extID)

		manifest := extism.Manifest{
			Wasm: []extism.Wasm{
				extism.WasmUrl{
					Url: c.GetURL(binding.ContentAddress),
				},
			},
		}

		config := extism.PluginConfig{}
		plugin, err := extism.NewPlugin(ctx, manifest, config, hostFunctions)
		if err != nil {
			log.Fatalf("Failed to initialize plugin %q: %v\n", name, err)
		}

		if err := exercisePlugin(extID, plugin); err != nil {
			log.Fatalf("Error calling plugin %v (extension point ID: %v): %v", name, extID, err)
		}
	}
}

var hostFunctions = []extism.HostFunction{
	eatAFruit,
}

var eatAFruit = extism.NewHostFunctionWithStack(
	"eatAFruit",
	func(ctx context.Context, plugin *extism.CurrentPlugin, stack []uint64) {
		buf, err := plugin.ReadBytes(stack[0])
		if err != nil {
			log.Printf("eatAFruit: ReadBytes err: %v", err)
			return
		}

		var input Fruit
		if err := jsoncomp.Unmarshal(buf, &input); err != nil {
			log.Printf("eatAFruit: Unmarshal err: %v", err)
			return
		}

		result := EatAFruit(input)

		buf, err = jsoncomp.Marshal(result)
		if err != nil {
			log.Printf("eatAFruit: Marshal err: %v", err)
			return
		}

		mem, err := plugin.WriteBytes(buf)
		if err != nil {
			log.Printf("eatAFruit: WriteBytes err: %v", err)
			return
		}

		stack[0] = mem
	},
	[]extism.ValueType{extism.ValueTypeI64},
	[]extism.ValueType{extism.ValueTypeI64},
)

func exercisePlugin(extID string, plugin *extism.Plugin) error {
	switch extID {
	case "ext_01j1gszkhmenz9ecq0cbvmm9mt":
		return exercisePluginFruit(plugin)
	case "ext_01j1gt0mmdez5teegkpwqtbv9q":
		return exercisePluginUser(plugin)
	default:
		return fmt.Errorf("unknown extension plugin id: %v", extID)
	}
	return nil
}

func exercisePluginFruit(plugin *extism.Plugin) error {
	{
		log.Printf("Calling fruit voidFunc()")
		success, outBuf, err := plugin.Call("voidFunc", nil)
		log.Printf("voidFunc returned: (%v, '%s', %v)", success, outBuf, err)
	}

	{
		log.Printf("Calling fruit primitiveTypeFunc(`'yo'`)")
		success, outBuf, err := plugin.Call("primitiveTypeFunc", []byte(`"yo"`))
		log.Printf("primitiveTypeFunc returned: (%v, '%s', %v)", success, outBuf, err)
	}

	{
		fruit := FruitEnumApple
		inBuf, err := jsoncomp.Marshal(fruit)
		if err != nil {
			return err
		}
		log.Printf("Calling fruit referenceTypeFunc(`%s`)", inBuf)
		success, outBuf, err := plugin.Call("referenceTypeFunc", inBuf)
		log.Printf("referenceTypeFunc returned: (%v, '%s', %v)", success, outBuf, err)
	}

	return nil
}

func exercisePluginUser(plugin *extism.Plugin) error {
	{
		user := User{
			Age:   intPtr(0),
			Email: stringPtr("email"),
			Address: &Address{
				Street: "street",
			},
		}
		inBuf, err := jsoncomp.Marshal(user)
		if err != nil {
			return err
		}
		log.Printf("Calling user processUser(`%s`)", inBuf)
		success, buf, err := plugin.Call("processUser", inBuf)
		log.Printf("processUser returned: (%v, '%s', %v)", success, buf, err)
	}

	return nil
}

func boolPtr(b bool) *bool       { return &b }
func intPtr(i int) *int          { return &i }
func stringPtr(s string) *string { return &s }

// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema map[string]string
