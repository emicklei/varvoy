package internal

import (
	"html/template"
	"log/slog"
	"os"
	"os/exec"
	"path"

	_ "embed"

	"golang.org/x/mod/modfile"
)

type Composer struct {
	workDir  string
	mainPath string
}

func NewComposer(mainPath string) *Composer {
	// dir := os.TempDir()
	dir := "/Users/emicklei/Projects/github.com/emicklei/varvoy/tmp"

	return &Composer{
		workDir:  path.Join(dir, "varvoy"),
		mainPath: mainPath,
	}
}

func (c *Composer) Compose() error {
	if err := osEnsureDir(c.workDir); err != nil {
		return err
	}

	modPath := path.Join(c.mainPath, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return err
	}
	mod, err := modfile.ParseLax(modPath, data, nil)
	if err != nil {
		return err
	}
	slog.Info("module", "path", mod.Module.Mod.Path)

	err = osCopy(modPath, path.Join(c.workDir, "go.mod"))
	if err != nil {
		wd, _ := os.Getwd()
		slog.Error("copy failed", "wd", wd, "err", err)
		return err
	}
	err = osCopy(path.Join(c.mainPath, "go.sum"), path.Join(c.workDir, "go.sum"))
	if err != nil {
		wd, _ := os.Getwd()
		slog.Error("copy failed", "wd", wd, "err", err)
		return err
	}
	err = genMain(c.workDir, mod.Module.Mod.Path)
	if err != nil {
		return err
	}

	// create imports folder
	importsDir := path.Join(c.workDir, "imports")
	if err := osEnsureDir(importsDir); err != nil {
		return err
	}

	// copy shared definition for Symbols
	err = os.WriteFile(path.Join(importsDir, "symbols.go"), symbolsTmpl, os.ModePerm)
	if err != nil {
		wd, _ := os.Getwd()
		slog.Error("copy failed", "wd", wd, "err", err)
		return err
	}

	// for yaegi Extracter to work, we need to be in the imports dir
	os.Chdir(importsDir)
	for _, each := range mod.Require {
		if err := yaegiExtract(each.Mod.Path); err != nil {
			return err
		}
	}
	os.Chdir(c.workDir)

	// add dependencies for the interpreter and varvoy
	err = goModTidy()
	if err != nil {
		return err
	}

	// build binary to connect and run
	err = goBuild("_debug_bin_varvoy")
	if err != nil {
		return err
	}

	return nil
}

//go:embed templates/debugbin.tmpl
var debugbinTmpl []byte

//go:embed templates/symbols.tmpl
var symbolsTmpl []byte

type debugbinData struct {
	ModPath string
}

func genMain(dir string, modpath string) error {
	tmpl, err := template.New("debugbin").Parse(string(debugbinTmpl))
	if err != nil {
		return err
	}
	out, err := os.Create(path.Join(dir, "main.go"))
	if err != nil {
		return err
	}
	defer out.Close()
	return tmpl.Execute(out, debugbinData{ModPath: modpath})
}

func goModTidy() error {
	slog.Debug("tidy go modules")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func goBuild(target string) error {
	slog.Debug("go build", "-o", target)
	cmd := exec.Command("go", "build", "-o", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
