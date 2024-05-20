package internal

import "testing"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestCompose(t *testing.T) {
	c := NewComposer("/Users/emicklei/Projects/github.com/emicklei/varvoy/todebug")
	err := c.Compose()
	check(err)
}
