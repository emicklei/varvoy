package main

import (
	"log"

	"github.com/emicklei/varvoy/api"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

func main() {
	opts := api.ListenOptions{
		BeforeAccept: func(addr string) {
			log.Println("tcp listening on", addr)
		},
	}
	api.ListenAndHandle(new(mockAdapter), opts)
}

type mockAdapter struct{}

func (a *mockAdapter) Initialize(s *dap.Session, ccaps *dap.InitializeRequestArguments) (*dap.Capabilities, error) {
	return &dap.Capabilities{
		SupportsConfigurationDoneRequest: dap.Bool(true),
		SupportsFunctionBreakpoints:      dap.Bool(true),
	}, nil
}

func (a *mockAdapter) Process(pm dap.IProtocolMessage) error {
	return nil
}

func (a *mockAdapter) Terminate() {
}
