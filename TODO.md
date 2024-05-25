
- how to pass command line arguments ?


```
// Start a process:
cmd := exec.Command("sleep", "5")
if err := cmd.Start(); err != nil {
    log.Fatal(err)
}

// Kill it:
if err := cmd.Process.Kill(); err != nil {
    log.Fatal("failed to kill process: ", err)
}
```

ps aux | grep varvoy | awk '{print $2}' | xargs kill



package main

import (
	"regexp"
	"testing"
)

func main() {
	testing.Init()
	testing.Main(regexp.MatchString, []testing.InternalTest{
		{
			Name: "TestThis",
			F:    TestThis,
		},
	}, []testing.InternalBenchmark{}, []testing.InternalExample{})
}

func TestThis(t *testing.T) {
	t.Log("this")
}
go run main.go -test.v


Passing FD to child process
https://groups.google.com/g/golang-nuts/c/Ws09uN64I80