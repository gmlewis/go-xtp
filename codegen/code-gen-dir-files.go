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
	if err := c.maybeWriteSourceFile(dirFile, fullSrc); err != nil {
		return err
	}

	if c.CustTypesTests != "" {
		testFilename := strings.Replace(dirFile, "."+c.Lang, "_test."+c.Lang, 1)
		testSrc := c.CustTypesTests
		if c.Lang == "go" {
			testSrc = fmt.Sprintf("package %v\n\n%v", c.PkgName, c.CustTypesTests)
		}
		if err := c.maybeWriteSourceFile(testFilename, testSrc); err != nil {
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

	return c.writeSrcFiles(dirName, hostSrc)
}

func (c *Client) GenPluginDir(dirName string) error {
	pluginSrc, err := c.GenPluginPDK()
	if err != nil {
		return err
	}

	return c.writeSrcFiles(dirName, pluginSrc)
}

func (c *Client) writeSrcFiles(dirName string, srcFiles GeneratedFiles) error {
	for filename, src := range srcFiles {
		dirFile := filepath.Join(dirName, filename)
		fileWriter := c.maybeWriteSourceFile
		if strings.HasSuffix(filename, ".sh") {
			fileWriter = c.maybeWriteScriptFile
		}
		if err := fileWriter(dirFile, src); err != nil {
			return err
		}
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
	if err := c.maybeWriteSourceFile(filename, moonPkgJSONFile); err != nil {
		return err
	}

	return nil
}

func (c *Client) maybeWriteScriptFile(path, buf string) error {
	return c.maybeWriteFile(path, buf, 0755)
}

func (c *Client) maybeWriteSourceFile(path, buf string) error {
	return c.maybeWriteFile(path, buf, 0644)
}

func (c *Client) maybeWriteFile(path, buf string, perm os.FileMode) error {
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

	if err := os.WriteFile(path, []byte(buf), perm); err != nil {
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
