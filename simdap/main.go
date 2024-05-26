package main

import (
	"log"
	"log/slog"
	"os"

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
	var lf = os.NewFile(uintptr(3), "varvoy-logs")
	a.session.Debug(lf)
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

	var body dap.ResponseBody
	switch m.Command {
	case "launch":
		slog.Debug("Event", "name", "initialized")
		a.session.Event("initialized", nil)
	case "setBreakpoints":
		body = &dap.SetBreakpointsResponseBody{}
	case "setFunctionBreakpoints":
		body = &dap.SetFunctionBreakpointsResponseBody{}
	}

	return a.session.Respond(m, true, "Success", body)
}

func (a *mockAdapter) Terminate() {
	slog.Debug("Terminate")
}
