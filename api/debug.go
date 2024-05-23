package api

import (
	"context"
	"reflect"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// Debug is called from the augmented debugging target binary
// It compiles a complete packages and starts a DAP listener.
func Debug(mainDir string, symbols map[string]map[string]reflect.Value) {
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
