package codegen

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) GenTypesDir(dirName string) error {
	typesSrc, err := c.GenCustomTypes()
	if err != nil {
		return err
	}

	return c.writeSrcFiles(dirName, typesSrc)
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
