package main

import (
	"log/slog"
	"os"

	"github.com/emicklei/dot"
	"github.com/emicklei/htmlslog"
)

func main() {
	g := dot.NewGraph(dot.Directed)
	l := htmlslog.New(os.Stdout, htmlslog.Options{})
	slog.SetDefault(slog.New(l))
	defer l.Close()
	slog.Info("done", "graph", g.String())
}
