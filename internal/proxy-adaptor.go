package internal

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

type ProxyAdapter struct {
	debugProcess *os.Process
	session      *dap.Session
	debugPort    int
	debugConn    net.Conn
}

// Initialize implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Initialize(s *dap.Session, ccaps *dap.InitializeRequestArguments) (*dap.Capabilities, error) {
	slog.Debug("Initialize", "ccaps", ccaps)
	a.session = s
	return &dap.Capabilities{
		SupportsConfigurationDoneRequest: dap.Bool(true),
		SupportsFunctionBreakpoints:      dap.Bool(true),
	}, nil
}

// Process implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Process(pm dap.IProtocolMessage) error {
	m, ok := pm.(*dap.Request)
	if !ok {
		return nil
	}
	var stop bool
	success := false
	var message string
	var body dap.ResponseBody

	switch m.Command {
	case "launch":
		args := m.Arguments.(*dap.LaunchRequestArguments)
		port, err := getFreePort()
		if err != nil {
			slog.Error("unable to get free tcp port", "err", err)
			return err
		}
		a.debugPort = port
		debugArgs := []string{
			"--log-dest=3", "--log", fmt.Sprintf("--listen=127.0.0.1:%d", port),
		}
		slog.Debug("launch", "args", args, "exec", "simdap", "arg", debugArgs)
		cmd := exec.Command("simdap", debugArgs...)
		if err := cmd.Start(); err != nil {
			slog.Error("unable to start program to debug", "err", err)
			return err
		}
		a.debugProcess = cmd.Process
		// connect to debugProces
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 5*time.Second)
		if err != nil {
			slog.Error("unable to get tcp connect to debug process", "err", err)
			return err
		}
		a.debugConn = conn
		// forward the launch request

		success = true

	case "disconnect":
		if err := a.killDebug(); err != nil {
			return err
		}
		success = true
	case "terminate":
		if err := a.killDebug(); err != nil {
			return err
		}
		success = true
	default:
		slog.Debug("forward process", "command", m.Command, "args", pm)
		// pass through
	}
	err := a.session.Respond(m, success, message, body)
	if err != nil {
		return err
	}

	if stop {
		return dap.ErrStop
	}
	return nil
}

// Terminate implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Terminate() {
	slog.Debug("Terminate")
	a.killDebug()
}

func (a *ProxyAdapter) killDebug() error {
	slog.Debug("disconnect and kill the debug process")
	if a.debugConn != nil {
		a.debugConn.Close()
	}
	if a.debugProcess != nil {
		if err := a.debugProcess.Kill(); err != nil {
			slog.Error("unable to kill debug process", "err", err)
			return err
		}
	}
	return nil
}
