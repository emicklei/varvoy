package api

import (
	"log/slog"
	"os"
	"testing"
)

func TestExecNoImports(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	exec("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/noimports", nil)
}
