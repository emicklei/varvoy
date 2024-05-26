package api

import (
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

type ListenOptions struct {
	BeforeAccept func(addr string)
}

func ListenAndHandle(adp dap.Handler, opts ListenOptions) {
	addr := flagValueString(getListenFlag())
	verbose := getLogFlag()
	logdest := flagValueInt(getLogDestFlag())

	// process log flags
	var lf = os.NewFile(uintptr(logdest), "varvoy-logs")
	lvl := slog.LevelInfo
	if verbose {
		lvl = slog.LevelDebug
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(lf, &tint.Options{
			Level:      lvl,
			TimeFormat: time.Kitchen,
			AddSource:  true,
		}),
	))
	slog.Debug("flags", "addr", addr, "logdest", logdest, "verbose", verbose)

	// connect
	var l net.Listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	srv := dap.NewServer(l, adp)

	if opts.BeforeAccept != nil {
		opts.BeforeAccept(addr)
	}

	// single session
	s, c, err := srv.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}
