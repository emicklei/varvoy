package internal

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/traefik/yaegi/extract"
)

// for yaegi Extracter to work, we need to be in the source dir
func yaegiExtractTo(require, targetDir string) error {
	wd, _ := os.Getwd()
	slog.Debug("yaegi extract", "require", require, "wd", wd, "targetdir", targetDir)
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
	f, err := os.Create(filepath.Join(targetDir, oFile))
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
