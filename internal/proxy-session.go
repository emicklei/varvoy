package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

// ErrStop is the error returned by a handler to indicate that the session
// should terminate.
var ErrStop = errors.New("stop")

// ProxySession forwards DAP requests to the debug process
// and passes DAP response back to the varvoy consumer (e.g. VS code)
type ProxySession struct {
	adapter *ProxyAdapter

	// for communicating to the debug process
	mux      *sync.Mutex
	inFlight map[int]*dap.Request

	dec *dap.Decoder
	enc *dap.Encoder
}

func NewProxySession(adapter *ProxyAdapter, conn net.Conn) *ProxySession {
	return &ProxySession{
		adapter:  adapter,
		mux:      new(sync.Mutex),
		inFlight: map[int]*dap.Request{},
		dec:      dap.NewDecoder(conn),
		enc:      dap.NewEncoder(conn),
	}
}

func (p *ProxySession) ForwardAndRespond(dapRequest *dap.Request, skipInFlight bool) error {
	slog.Debug("forward process", "command", dapRequest.Command, "skip", skipInFlight)
	if !skipInFlight {
		p.mux.Lock()
		p.inFlight[dapRequest.Seq] = dapRequest
		p.mux.Unlock()
	}

	// send
	err := p.enc.Encode(dapRequest)
	if err != nil {
		slog.Error("failed to forward request", "err", err, "r", dapRequest)
		return err
	}

	// receive
	pm, err := p.dec.Decode()
	if err != nil {
		slog.Error("failed to receive response", "err", err, "r", dapRequest)
		return err
	}
	dapResponse, ok := pm.(*dap.Response)
	if !ok {
		return fmt.Errorf("loop: response expected, got: %[1]v(%[1]T)", pm)
	}
	slog.Debug("recv", "response-seq", dapResponse.Seq, "success", dapResponse.Success, "command", dapResponse.Command)

	err = p.adapter.session.Respond(dapRequest, true, dapResponse.Message.Get(), dapResponse.Body)
	if errors.Is(err, ErrStop) {
		return nil
	} else if err != nil {
		slog.Error("unable to respond", "err", err, "response-seq", dapResponse.Seq)
		return err
	}
	return nil
}

// Run starts the session. Run blocks until the session is terminated.
func (p *ProxySession) Run() error {
	// forever { receive response from debugProcess, respond response to proxy consumer }
	for {
		pm, err := p.recv()
		if err != nil {
			return fmt.Errorf("loop: decode: %w", err)
		}
		dapResponse, ok := pm.(*dap.Response)
		if !ok {
			return fmt.Errorf("loop: response expected, got: %[1]v(%[1]T)", pm)
		}
		slog.Debug("recv", "response-seq", dapResponse.Seq, "success", dapResponse.Success, "command", dapResponse.Command)

		p.mux.Lock()
		dapRequest, ok := p.inFlight[dapResponse.Seq]
		if !ok {
			p.mux.Unlock()
			slog.Warn("no matching request found", "response-seq", dapResponse.Seq)
			continue
		}
		delete(p.inFlight, dapRequest.Seq)
		p.mux.Unlock()

		err = p.adapter.session.Respond(dapRequest, true, dapResponse.Message.Get(), dapResponse.Body)
		if errors.Is(err, ErrStop) {
			break
		} else if err != nil {
			slog.Error("unable to respond", "err", err, "response-seq", dapResponse.Seq)
			return err
		}
	}
	// call terminate
	return nil
}

// recv receives from debugProcess
func (p *ProxySession) recv() (dap.IProtocolMessage, error) {
	m, err := p.dec.Decode()
	if err == nil {
		return m, nil
	}
	return nil, err
}
