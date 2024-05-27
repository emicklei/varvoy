package main

import (
	"fmt"

	"github.com/emicklei/varvoy/api"
	"github.com/emicklei/varvoy/internal"
)

const Version = "0.0.3"

// This program accepts the Delve (dlv) flags and args because that is hardcoded in the vscode-go plugin. Example is:
//
//	dap --listen=127.0.0.1:52950 --log-dest=3 --log
func main() {
	adp := new(internal.ProxyAdapter)
	opts := api.ListenOptions{
		BeforeAccept: func(addr string) {
			// Line must start with "DAP server listening at:"
			// see https://github.com/golang/vscode-go/blob/f907536117c3e9fc731be9277e992b8cc7cd74f1/extension/src/goDebugFactory.ts#L558
			fmt.Println("DAP server listening at:", addr, fmt.Sprintf("(varvoy:%s)", Version))
		},
	}
	api.ListenAndHandle(adp, opts)

}
