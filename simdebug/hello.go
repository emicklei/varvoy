package main

import (
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

func hello() {
	drv := newDriver("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/hello")

	breakOnLine := 6

	drv.Cmd("launch", &dap.LaunchRequestArguments{})

	drv.Cmd("setBreakpoints", &dap.SetBreakpointsArguments{
		Source: dap.Source{
			Path: dap.Str(drv.srcPath),
		},
		Breakpoints: []*dap.SourceBreakpoint{{Line: breakOnLine}},
		Lines:       []int{breakOnLine},
	})

	drv.Cmd("configurationDone", &dap.ConfigurationDoneArguments{})

	drv.Cmd("threads", &dap.ThreadsArguments{})

	drv.Cmd("stackTrace", &dap.StackTraceArguments{ThreadId: 1})

	drv.Cmd("continue", &dap.ContinueArguments{ThreadId: 1})
}
