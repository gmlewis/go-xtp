package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) GenTypesDirFiles(dirFile string) error {
	pkgName := c.genPkgName(dirFile)
	fullSrc := c.CustTypes
	if c.Lang == "go" {
		fullSrc = fmt.Sprintf("// Package %v represents the custom datatypes for an XTP Extension Plugin.\npackage %[1]v\n\n%v", pkgName, c.CustTypes)
	}
	if err := os.WriteFile(dirFile, []byte(fullSrc), 0644); err != nil {
		return err
	}

	if c.CustTypesTests != "" {
		testFilename := strings.Replace(dirFile, "."+c.Lang, "_test."+c.Lang, 1)
		testSrc := c.CustTypesTests
		if c.Lang == "go" {
			testSrc = fmt.Sprintf("package %v\n\n%v", pkgName, c.CustTypesTests)
		}
		if err := os.WriteFile(testFilename, []byte(testSrc), 0644); err != nil {
			return err
		}
	}

	if c.Lang == "mbt" {
		if err := genMoonPkgJsonFileIfNeeded(filepath.Dir(dirFile)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GenHostDirFiles(dirFile string) error {
	pkgName := c.genPkgName(dirFile)
	hostSrc, err := c.GenHostSDK(pkgName)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dirFile, []byte(hostSrc), 0644); err != nil {
		return err
	}

	if c.Lang == "mbt" {
		if err := genMoonPkgJsonFileIfNeeded(filepath.Dir(dirFile)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GenPluginDirFiles(dirFile string) error {
	pkgName := c.genPkgName(dirFile)
	pluginSrc, err := c.GenPluginPDK(pkgName)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dirFile, []byte(pluginSrc), 0644); err != nil {
		return err
	}

	if c.Lang == "mbt" {
		if err := genMoonPkgJsonFileIfNeeded(filepath.Dir(dirFile)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) genPkgName(filename string) string {
	baseName := strings.Replace(filepath.Base(filename), "."+c.Lang, "", 1)
	return strings.ToLower(strings.Replace(baseName, "-", "_", -1))
}

func genMoonPkgJsonFileIfNeeded(dirname string) error {
	filename := filepath.Join(dirname, "moon.mod.json")
	_, err := os.Stat(filename)
	if err == nil || !os.IsNotExist(err) {
		return err
	}

	// create the file
	return os.WriteFile(filename, []byte(moonModJSONFile), 0644)
}

const moonModJSONFile = `{
  "import": [
    {
      "path": "gmlewis/json",
      "alias": "jsonutil"
    }
  ]
}`
