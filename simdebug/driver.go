package main

import (
	"fmt"
	"go/build"
	"os"

	dbg "github.com/traefik-contrib/yaegi-debug-adapter"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type driver struct {
	seq     int
	adp     *dbg.Adapter
	srcPath string
}

func (d *driver) Cmd(command string, args dap.RequestArguments) error {
	d.seq++
	r := &dap.Request{
		ProtocolMessage: dap.ProtocolMessage{Seq: d.seq},
		Command:         command,
		Arguments:       args,
	}
	err := d.adp.Process(r)
	if err != nil {
		fmt.Println("process error:", err)
	}
	return err
}

func newDriver(srcPath string) *driver {
	newInterp := func(opts interp.Options) (*interp.Interpreter, error) {
		opts.GoPath = build.Default.GOPATH
		opts.Stderr = os.Stderr
		opts.Stdout = os.Stdout
		i := interp.New(opts)
		if err := i.Use(stdlib.Symbols); err != nil {
			return nil, err
		}
		if err := i.Use(interp.Symbols); err != nil {
			return nil, err
		}
		i.ImportUsed()
		return i, nil
	}
	// capture errors
	errch := make(chan error)
	go func() {
		for err := range errch {
			fmt.Printf("ERR %v\n", err)
		}
	}()
	// defer close(errch)

	opts := &dbg.Options{
		StopAtEntry:    false, // true then stop at first statement of main
		NewInterpreter: newInterp,
		Errors:         errch,
		SrcPath:        srcPath,
	}

	adp := dbg.NewAdapter(srcPath, (*interp.Interpreter).CompilePackage, opts)

	ses := dap.NewSession(os.Stdin, os.Stdout, adp)
	ses.Debug(os.Stderr)
	_, _ = adp.Initialize(ses, new(dap.InitializeRequestArguments))
	return &driver{adp: adp, srcPath: srcPath}
}
