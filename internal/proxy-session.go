package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

// ErrStop is the error returned by a handler to indicate that the session
// should terminate.
var ErrStop = errors.New("stop")

// ProxySession forwards DAP requests to the debug process
// and passes DAP response back to the varvoy consumer (e.g. VS code)
type ProxySession struct {
	adapter *ProxyAdapter
	dec     *dap.Decoder
	enc     *dap.Encoder
}

func NewProxySession(adapter *ProxyAdapter, conn net.Conn) *ProxySession {
	return &ProxySession{
		adapter: adapter,
		dec:     dap.NewDecoder(conn),
		enc:     dap.NewEncoder(conn),
	}
}

func (p *ProxySession) ForwardAndRespond(dapRequest *dap.Request) error {
	// send to downstream
	err := p.enc.Encode(dapRequest)
	if err != nil {
		slog.Error("failed to forward request", "err", err, "r", dapRequest)
		return err
	}
	slog.Debug("forward", "seq", dapRequest.ProtocolMessage.Seq, "command", dapRequest.Command)

	// receive from downstream
	pm, err := p.dec.Decode()
	if err != nil {
		slog.Error("failed to receive response", "err", err, "r", dapRequest)
		return err
	}
	dapResponse, ok := pm.(*dap.Response)
	if !ok {
		return fmt.Errorf("loop: response expected, got: %[1]v(%[1]T)", pm)
	}
	slog.Debug("received", "seq", dapResponse.Seq, "success", dapResponse.Success, "command", dapResponse.Command)

	// respond back to upstream
	err = p.adapter.session.Respond(dapRequest, true, dapResponse.Message.Get(), dapResponse.Body)
	if errors.Is(err, ErrStop) {
		return nil
	} else if err != nil {
		slog.Error("unable to respond", "err", err, "response-seq", dapResponse.Seq)
		return err
	}
	return nil
}
