package internal

import (
	"log/slog"
	"os"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func enableDebugLog() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func TestCompose(t *testing.T) {
	enableDebugLog()
	c := NewComposer("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/withimports")
	err := c.Compose()
	check(err)
}
