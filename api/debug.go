package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	dbg "github.com/traefik-contrib/yaegi-debug-adapter"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Debug is called from the augmented debugging target binary
// It compiles a complete package and starts a DAP listener.
// For debugging,
// when the binary is started with VARVOY_EXEC=true then execute the program with an interpreter.
// when the binary is started with VARVOY_RUN=true then continue the program with a debugger.
func Debug(mainDir string, symbols map[string]map[string]reflect.Value) {
	if os.Getenv("VARVOY_EXEC") != "" {
		exec(mainDir, symbols)
		return
	}
	if os.Getenv("VARVOY_RUN") != "" {
		run(mainDir, symbols)
		return
	}
	newInterp := func(opts interp.Options) (*interp.Interpreter, error) {
		i := interp.New(opts)
		if err := i.Use(stdlib.Symbols); err != nil {
			return nil, err
		}
		if err := i.Use(symbols); err != nil {
			return nil, err
		}
		return i, nil
	}
	errch := make(chan error)
	go func() {
		for err := range errch {
			fmt.Printf("ERR %v\n", err)
		}
	}()
	defer close(errch)

	debugOpts := &dbg.Options{
		StopAtEntry:    false,
		NewInterpreter: newInterp,
		Errors:         errch,
		SrcPath:        filepath.Join(mainDir, "main.go"), // TODO
	}
	adp := dbg.NewAdapter((*interp.Interpreter).CompilePackage, mainDir, debugOpts)
	ListenAndHandle(adp, ListenOptions{})
}

func exec(mainDir string, symbols map[string]map[string]reflect.Value) {
	i := interp.New(interp.Options{})
	if err := i.Use(stdlib.Symbols); err != nil {
		fmt.Println("use stdlib failed:", err)
		os.Exit(1)
	}
	if err := i.Use(symbols); err != nil {
		fmt.Println("use compiled symbols failed:", err)
		os.Exit(1)
	}
	prog, err := i.CompilePackage(mainDir)
	if err != nil {
		fmt.Println("compile package failed:", err)
		os.Exit(1)
	}
	res, err := i.ExecuteWithContext(context.Background(), prog)
	if err != nil {
		fmt.Println("execute with context failed:", err, "result", res)
		os.Exit(1)
	}
}

func run(mainDir string, symbols map[string]map[string]reflect.Value) {
	i := interp.New(interp.Options{})
	_ = i.Use(stdlib.Symbols)
	_ = i.Use(symbols)
	prog, err := i.CompilePackage(mainDir)
	if err != nil {
		fmt.Println("compile package failed:", err)
		os.Exit(1)
	}
	dbg := i.Debug(context.Background(), prog, func(de *interp.DebugEvent) {
		slog.Debug("handle", "event", de)
	}, &interp.DebugOptions{GoRoutineStartAt1: true})
	if err := dbg.Continue(1); err != nil {
		fmt.Println("cannot continue go-routine 1:", err)
	}
	_, err = dbg.Wait()
	if err != nil {
		fmt.Println("cannot wait:", err)
	}
}
