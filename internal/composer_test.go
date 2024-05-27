package internal

import (
	"fmt"
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
	os.Mkdir("tmp", os.ModePerm)
	opts := ComposeOptions{
		MainDir: "/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug/withimports",
		TempDir: "/Users/emicklei/Projects/github.com/emicklei/varvoy/internal/tmp",
	}
	c := NewExecutableComposer(opts)
	err := c.Compose()
	check(err)
	fmt.Println("VARVOY_RUN=true ", c.FullExecName())
}
