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
	workDir string
	mainDir string
}

func NewComposer(mainDir string) *Composer {
	// dir := os.TempDir()
	dir := "/Users/emicklei/Projects/github.com/emicklei/varvoy/tmp"

	return &Composer{
		workDir: path.Join(dir, "varvoy"),
		mainDir: mainDir,
	}
}

func (c *Composer) Compose() error {
	if err := osEnsureDir(c.workDir); err != nil {
		return err
	}

	modPath := path.Join(c.mainDir, "go.mod")
	data, err := os.ReadFile(modPath)
	if err != nil {
		return err
	}
	mod, err := modfile.ParseLax(modPath, data, nil)
	if err != nil {
		return err
	}
	slog.Info("module", "path", mod.Module.Mod.Path)

	// modify module to add replace
	// replace github.com/traefik/yaegi => ../../../yaegi
	// TODO
	// mod.AddRequire("github.com/traefik/yaegi", "v0.16.1")
	if err := mod.AddReplace("github.com/traefik/yaegi", "", "../../../yaegi", ""); err != nil {
		return err
	}
	// mod.AddRequire("github.com/emicklei/varvoy", "v0.0.0")
	if err := mod.AddReplace("github.com/emicklei/varvoy", "", "../../../varvoy", ""); err != nil {
		return err
	}
	// write mod
	modContent, err := mod.Format()
	if err != nil {
		return err
	}
	os.WriteFile(path.Join(c.workDir, "go.mod"), modContent, os.ModePerm)

	err = genMain(c.mainDir, c.workDir, mod.Module.Mod.Path)
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
	os.Chdir(c.mainDir)
	for _, each := range mod.Require {
		if err := yaegiExtractTo(each.Mod.Path, importsDir); err != nil {
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
	MainDir string
}

func genMain(mainDir, targetDir string, modpath string) error {
	slog.Debug("write main.go", "maindir", mainDir, "targetdir", targetDir, "modpath", modpath)
	tmpl, err := template.New("debugbin").Parse(string(debugbinTmpl))
	if err != nil {
		return err
	}
	out, err := os.Create(path.Join(targetDir, "main.go"))
	if err != nil {
		return err
	}
	defer out.Close()
	return tmpl.Execute(out, debugbinData{ModPath: modpath, MainDir: mainDir})
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
