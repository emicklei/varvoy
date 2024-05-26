package main

import (
	"log"
	"log/slog"

	"github.com/emicklei/varvoy/api"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

func main() {
	opts := api.ListenOptions{
		BeforeAccept: func(addr string) {
			log.Println("simdap listening at:", addr)
		},
	}
	api.ListenAndHandle(new(mockAdapter), opts)
}

type mockAdapter struct {
	session *dap.Session
}

func (a *mockAdapter) Initialize(s *dap.Session, ccaps *dap.InitializeRequestArguments) (*dap.Capabilities, error) {
	a.session = s
	return &dap.Capabilities{
		SupportsConfigurationDoneRequest: dap.Bool(true),
		SupportsFunctionBreakpoints:      dap.Bool(true),
	}, nil
}

func (a *mockAdapter) Process(pm dap.IProtocolMessage) error {
	m, ok := pm.(*dap.Request)
	if !ok {
		return nil
	}
	slog.Debug("Process", "command", m.Command, "seq", m.Seq)
	return a.session.Respond(m, true, "Success", nil)
}

func (a *mockAdapter) Terminate() {
	slog.Debug("Terminate")
}
