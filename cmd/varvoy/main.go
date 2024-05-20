package main

import (
	"flag"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/emicklei/varvoy/internal"
	"github.com/traefik-contrib/yaegi-debug-adapter/pkg/dap"
)

func main() {
	var (
		addr string
	)
	flag.StringVar(&addr, "addr", "tcp://localhost:16348", "Net address to listen on, must be a TCP or Unix socket URL")

	// parse flag
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}

	// get addrs
	var addrs string
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
	l, err = net.Listen(u.Scheme, addrs)
	if err != nil {
		log.Fatal(err)
	}
	adp := new(internal.ProxyAdapter)
	srv := dap.NewServer(l, adp)

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
