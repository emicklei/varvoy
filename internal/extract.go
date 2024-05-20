package internal

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/traefik/yaegi/extract"
)

// for yaegi Extracter to work, we need to be in the imports dir
func yaegiExtract(require string) error {
	// TODO license,tag,exclude,include
	ext := extract.Extractor{
		Dest:    "imports",
		License: "",
	}
	repl := strings.NewReplacer("/", "-", ".", "_", "~", "_")
	var buf bytes.Buffer
	importPath, err := ext.Extract(require, "imports", &buf)
	if err != nil {
		return err
	}
	oFile := repl.Replace(importPath) + ".go"
	f, err := os.Create(oFile)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, &buf); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
