package api

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/emicklei/structexplorer"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func TestExecNoImports(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	exec("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/noimports", nil)
}

func TestProgramExplore(t *testing.T) {
	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		fmt.Println("use stdlib failed:", err)
		os.Exit(1)
	}

	prog, err := i.CompilePackage("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/noimports")
	if err != nil {
		fmt.Println("compile package failed:", err)
		os.Exit(1)
	}
	structexplorer.NewService("interpreter", i, "prog", prog).Start()
}
