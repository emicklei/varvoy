package internal

import (
	"errors"
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
	dec     *dap.Decoder
	enc     *dap.Encoder
	sendMux *sync.Mutex
}

func NewProxySession(adapter *ProxyAdapter, conn net.Conn) *ProxySession {
	return &ProxySession{
		adapter: adapter,
		dec:     dap.NewDecoder(conn),
		enc:     dap.NewEncoder(conn),
		sendMux: new(sync.Mutex),
	}
}

func (p *ProxySession) Forward(dapRequest *dap.Request) error {
	p.sendMux.Lock()
	defer p.sendMux.Unlock()

	// send to downstream
	err := p.enc.Encode(dapRequest)
	if err != nil {
		slog.Error("failed to forward request", "err", err, "r", dapRequest)
		return err
	}
	slog.Debug("forwarded", "seq", dapRequest.ProtocolMessage.Seq, "command", dapRequest.Command)
	/**
		// receive from downstream
		pm, err := p.dec.Decode()
		if err != nil {
			slog.Error("failed to receive response", "err", err, "request-seq", dapRequest.ProtocolMessage.Seq)
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
	**/
	return nil
}

// ReceiveAndRespond receives protocolmessages (Response,Event) and respond with them
func (p *ProxySession) ReceiveAndRespond() error {
	slog.Debug("receiving and responding messages...")

	for {
		pm, err := p.dec.Decode()
		if err != nil {
			slog.Error("failed to decode message", "err", err)
			return err
		}
		slog.Debug("received from downstream", "pm", pm)
		dapResponse, ok := pm.(*dap.Response)
		if ok {
			req := new(dap.Request)
			req.Command = dapResponse.Command
			req.Seq = dapResponse.Seq
			if err := p.adapter.session.Respond(req, dapResponse.Success, dapResponse.Message.Get(), dapResponse.Body); err != nil {
				return err
			}
			slog.Debug("responded to upstream", "seq", req.Seq, "command", req.Command)
		} else {
			dapEvent, ok := pm.(*dap.Event)
			if ok {
				slog.Debug("todo event", "event", dapEvent)
			} else {
				slog.Debug("unhandled message", "pm", pm)
			}
		}
	}
}
