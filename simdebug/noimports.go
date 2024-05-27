package main

import "github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"

func noimports() {
	drv := newDriver("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/noimports")
	drv.Cmd("launch", &dap.LaunchRequestArguments{})
}
