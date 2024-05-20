package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/emicklei/varvoy/internal"
	"github.com/lmittmann/tint"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

const Version = "0.0.1"

// This program accepts the Delve (dlv) flags and args because that is hardcoded in the vscode-go plugin. Example is:
//
//	dap --listen=127.0.0.1:52950 --log-dest=3 --log
func main() {
	var (
		addr    string
		verbose bool
		logdest int
	)
	flag.StringVar(&addr, "listen", "", "host:port to listen on")
	flag.IntVar(&logdest, "log-dest", 0, "file descriptor to log to")
	flag.BoolVar(&verbose, "log", false, "Verbose logging")
	flag.Parse()

	// process log flags
	var lf = os.NewFile(uintptr(logdest), "yaegi-dap-vscode-logs")
	lvl := slog.LevelInfo
	if verbose {
		lvl = slog.LevelDebug
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(lf, &tint.Options{
			Level:      lvl,
			TimeFormat: time.Kitchen,
		}),
	))

	// connect
	var l net.Listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	adp := new(internal.ProxyAdapter)
	srv := dap.NewServer(l, adp)

	// Line must start with "DAP server listening at:"
	// see https://github.com/golang/vscode-go/blob/f907536117c3e9fc731be9277e992b8cc7cd74f1/extension/src/goDebugFactory.ts#L558
	fmt.Println("DAP server listening at:", addr, fmt.Sprintf("(yaegi-dap-vscode:%s)", Version))

	// single session
	slog.Debug("accepting...")
	s, c, err := srv.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	slog.Debug("running...")
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}
