package api

import (
	"context"
	"fmt"
	"reflect"

	dbg "github.com/traefik-contrib/yaegi-debug-adapter"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Debug is called from the augmented debugging target binary
// It compiles a complete package and starts a DAP listener.
func Debug(mainDir string, symbols map[string]map[string]reflect.Value) {
	newInterp := func(opts interp.Options) (*interp.Interpreter, error) {
		i := interp.New(opts)
		i.Use(stdlib.Symbols)
		i.Use(symbols)
		i.ImportUsed()
		_, err := i.CompilePackage(mainDir)
		if err != nil {
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
		SrcPath:        mainDir,
	}
	adp := dbg.NewAdapter(mainDir, (*interp.Interpreter).CompilePackage, debugOpts)
	ListenAndHandle(adp, ListenOptions{})
}

func Run(mainDir string, symbols map[string]map[string]reflect.Value) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)
	i.Use(symbols)
	i.ImportUsed()
	prog, err := i.CompilePackage(mainDir)
	if err != nil {
		// TODO
		panic(err)
	}
	_, err = i.ExecuteWithContext(context.Background(), prog)
	if err != nil {
		// TODO
		panic(err)
	}
}
