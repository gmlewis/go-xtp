package codegen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) GenTypesDir(dirName string) error {
	fullSrc := c.CustTypes
	if c.Lang == "go" {
		fullSrc = fmt.Sprintf("// Package %v represents the custom datatypes for an XTP Extension Plugin.\npackage %[1]v\n\n%v", c.PkgName, c.CustTypes)
	}

	dirFile := filepath.Join(dirName, fmt.Sprintf("%v.%v", c.PkgName, c.Lang))
	if err := c.maybeWriteFile(dirFile, fullSrc); err != nil {
		return err
	}

	if c.CustTypesTests != "" {
		testFilename := strings.Replace(dirFile, "."+c.Lang, "_test."+c.Lang, 1)
		testSrc := c.CustTypesTests
		if c.Lang == "go" {
			testSrc = fmt.Sprintf("package %v\n\n%v", c.PkgName, c.CustTypesTests)
		}
		if err := c.maybeWriteFile(testFilename, testSrc); err != nil {
			return err
		}
	}

	if c.Lang == "mbt" {
		if err := c.genMoonPkgJsonFileIfNeeded(filepath.Dir(dirFile)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GenHostDir(dirName string) error {
	hostSrc, err := c.GenHostSDK()
	if err != nil {
		return err
	}

	dirFile := filepath.Join(dirName, fmt.Sprintf("%v.%v", c.PkgName, c.Lang))
	if err := c.maybeWriteFile(dirFile, hostSrc); err != nil {
		return err
	}

	if c.Lang == "mbt" {
		if err := c.genMoonPkgJsonFileIfNeeded(dirName); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GenPluginDir(dirName string) error {
	pluginSrc, err := c.GenPluginPDK()
	if err != nil {
		return err
	}

	dirFile := filepath.Join(dirName, fmt.Sprintf("%v.%v", c.PkgName, c.Lang))
	if err := c.maybeWriteFile(dirFile, pluginSrc); err != nil {
		return err
	}

	if c.Lang == "mbt" {
		if err := c.genMoonPkgJsonFileIfNeeded(dirName); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) genMoonPkgJsonFileIfNeeded(dirName string) error {
	filename := filepath.Join(dirName, "moon.pkg.json")
	if err := c.maybeWriteFile(filename, moonPkgJSONFile); err != nil {
		return err
	}

	return nil
}

func (c *Client) maybeWriteFile(path, buf string) error {
	parent := filepath.Dir(path)
	if _, err := os.Stat(parent); parent != "." && err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(parent, 0755); err != nil {
			return err
		}
	}

	if !c.force {
		if _, err := os.Stat(path); err == nil {
			log.Printf("WARNING: not writing file %q - add -force to overwrite", path)
			return nil
		}
	}

	if err := os.WriteFile(path, []byte(buf), 0644); err != nil {
		return err
	}

	return nil
}

const moonPkgJSONFile = `{
  "import": [
    {
      "path": "gmlewis/json",
      "alias": "jsonutil"
    }
  ]
}`
