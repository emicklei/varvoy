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

func main() {
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
	defer close(errch)

	srcPath := "/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/hello"
	breakOnLine := 6

	opts := &dbg.Options{
		StopAtEntry:    false, // true then stop at first statement of main
		NewInterpreter: newInterp,
		Errors:         errch,
		SrcPath:        srcPath,
	}

	adp := dbg.NewAdapter((*interp.Interpreter).CompilePackage, srcPath, opts)

	ses := dap.NewSession(os.Stdin, os.Stdout, adp)
	ses.Debug(os.Stderr)

	_, _ = adp.Initialize(ses, new(dap.InitializeRequestArguments))

	r0 := &dap.Request{
		ProtocolMessage: dap.ProtocolMessage{Seq: 1},
		Command:         "launch",
		Arguments:       &dap.LaunchRequestArguments{},
	}
	if err := adp.Process(r0); err != nil {
		fmt.Println("process error:", err)
	}

	r1 := &dap.Request{
		ProtocolMessage: dap.ProtocolMessage{Seq: 2},
		Command:         "setBreakpoints",
		Arguments: &dap.SetBreakpointsArguments{
			Source: dap.Source{
				Path: dap.Str(srcPath),
			},
			Breakpoints: []*dap.SourceBreakpoint{{Line: breakOnLine}},
			Lines:       []int{breakOnLine},
		},
	}

	if err := adp.Process(r1); err != nil {
		fmt.Println("process error:", err)
	}

	r2 := &dap.Request{
		ProtocolMessage: dap.ProtocolMessage{Seq: 3},
		Command:         "configurationDone",
		Arguments:       &dap.ConfigurationDoneArguments{},
	}

	if err := adp.Process(r2); err != nil {
		fmt.Println("process error:", err)
	}

	r3 := &dap.Request{
		ProtocolMessage: dap.ProtocolMessage{Seq: 4},
		Command:         "continue",
		Arguments: &dap.ContinueArguments{
			ThreadId: 1,
		},
	}

	if err := adp.Process(r3); err != nil {
		fmt.Println("process error:", err)
	}
}
