package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

const (
	eventInitialized = "initialized"
	respondSuccess   = "Success"
)

type ProxyAdapter struct {
	session *dap.Session

	// for the target program to debug
	proxySession *ProxySession
	debugProcess *os.Process
	debugConn    net.Conn
	debugBin     string
}

// Initialize implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Initialize(s *dap.Session, ccaps *dap.InitializeRequestArguments) (*dap.Capabilities, error) {
	slog.Debug("Initialize")
	a.session = s

	// create temporary binary
	wd, _ := os.Getwd()
	opts := ComposeOptions{
		TempDir: os.TempDir(),
		MainDir: wd,
	}
	comp := NewExecutableComposer(opts)
	if err := comp.Compose(); err != nil {
		slog.Error("unable to create debug binary", "err", err)
		return nil, err
	}
	a.debugBin = comp.FullExecName()

	//a.debugBin := "simdap"

	// fire up debug process
	port, err := getFreePort()
	if err != nil {
		slog.Error("unable to get free tcp port", "err", err)
		return nil, err
	}
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	// TODO get the 3 from flag
	debugArgs := []string{
		"--log-dest=3", "--log", fmt.Sprintf("--listen=%s", addr),
	}

	slog.Debug("start exec", "bin", a.debugBin, "args", debugArgs)
	cmd := exec.Command(a.debugBin, debugArgs...)
	if err := cmd.Start(); err != nil {
		slog.Error("unable to start program to debug", "err", err)
		return nil, err
	}
	a.debugProcess = cmd.Process
	// connect to debugProcess
	attempts := 5
	var conn net.Conn
	for {
		if attempts == 0 {
			slog.Error("unable to get tcp connect to debug process", "err", err)
			return nil, err
		}

		conn, err = net.Dial("tcp", addr)
		if err == nil {
			break
		}
		slog.Debug("failed to dial", "err", err)
		slog.Info("waiting for the debug process...", "addr", addr)
		time.Sleep(1 * time.Second)
		attempts--
	}
	a.proxySession = NewProxySession(a, conn)

	slog.Debug("simulate an initialize that was received by the session")
	dapRequest := new(dap.Request)
	dapRequest.Seq = 1 // guess
	dapRequest.Type = dap.ProtocolMessageType_Request
	dapRequest.Command = "initialize"
	dapRequest.Arguments = ccaps

	if err := a.proxySession.Forward(dapRequest); err != nil {
		slog.Error("unable to forward initialize request to debug process", "err", err)
		return nil, err
	}

	caps, err := a.proxySession.ReceiveInitializeResponse()
	if err != nil {
		return nil, err
	}

	go func() {
		if err := a.proxySession.Run(); err != nil {
			slog.Error("Run failed", "err", err)
		}
	}()

	return caps, nil
}

// Process implements dap.Handler and should not be called directly.
func (a *ProxyAdapter) Process(pm dap.IProtocolMessage) error {
	m, ok := pm.(*dap.Request)
	if !ok {
		return nil
	}
	slog.Debug("Process", "command", m.Command)
	var stop bool
	success := false

	switch m.Command {

	case "disconnect":
		if a.debugConn != nil {
			a.debugConn.Close()
			a.debugConn = nil
		}
		stop = true // only at disconnect
		success = true
	case "terminate":
		if err := a.tearDown(); err != nil {
			return err
		}
		success = true
	default:
		if err := a.proxySession.Forward(m); err != nil {
			slog.Error("unable to forward request to debug process", "err", err)
			return err
		}
		return nil
	}
	err := a.session.Respond(m, success, respondSuccess, nil)
	if err != nil {
		slog.Error("unable to respong to Process", "command", m.Command, "err", err)
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
	a.tearDown()
}

func (a *ProxyAdapter) tearDown() error {
	slog.Debug("close connection")
	if a.debugConn != nil {
		a.debugConn.Close()
	}
	slog.Debug("kill the debug process")
	if a.debugProcess != nil {
		if err := a.debugProcess.Kill(); err != nil {
			slog.Error("unable to kill debug process", "err", err)
			return err
		}
		a.debugProcess = nil
	}
	slog.Debug("remove the temporary binary")
	if a.debugBin != "" {
		// keep binary around when debugging
		if !slog.Default().Enabled(context.Background(), slog.LevelDebug) {
			_ = os.Remove(a.debugBin)
		} else {
			fmt.Println("VARVOY_RUN=true", a.debugBin)
			fmt.Println("VARVOY_EXEC=true", a.debugBin)
		}
	}
	return nil
}
