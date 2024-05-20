package internal

import "github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"

type ProxyAdapter struct{}

// Initialize implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Initialize(s *dap.Session, ccaps *dap.InitializeRequestArguments) (*dap.Capabilities, error) {
	return &dap.Capabilities{
		SupportsConfigurationDoneRequest: dap.Bool(true),
		SupportsFunctionBreakpoints:      dap.Bool(true),
	}, nil
}

// Process implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Process(pm dap.IProtocolMessage) error {
	return nil
}

// Terminate implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Terminate() {}
