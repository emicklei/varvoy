package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/emicklei/varvoy/internal"
	"github.com/lmittmann/tint"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

func main() {
	var (
		addr    string
		verbose bool
	)
	flag.StringVar(&addr, "addr", "", "Net address to listen on, must be a TCP or Unix socket URL")
	flag.BoolVar(&verbose, "log", false, "Verbose logging")
	flag.Parse()

	// setuplogging
	lvl := slog.LevelInfo
	if verbose {
		lvl = slog.LevelDebug
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      lvl,
			TimeFormat: time.Kitchen,
		}),
	))

	// parse addr flag
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}

	// get addrs
	if u.Scheme == "unix" {
		addr = u.Path
		if _, err = os.Stat(addr); err == nil {
			// Remove any pre-existing connection
			_ = os.Remove(addr)
		}

		// Remove when done
		defer func() { _ = os.Remove(addr) }()
	} else {
		addr = u.Host
	}

	// connect
	var l net.Listener
	l, err = net.Listen(u.Scheme, addr)
	if err != nil {
		log.Fatal(err)
	}
	adp := new(internal.ProxyAdapter)
	srv := dap.NewServer(l, adp)

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
