package main

import (
	"fmt"
	"time"
)

func main() {
	m := map[bool]bool{true: false}
	i := 42
	iptr := &i
	f := float64(3.14)
	type A struct{ B string }
	var a A
	ar := [2]int{-1, 1}
	s := "hello varvoy"
	fmt.Println(s, m, i, f, a, ar)
	ml := map[string]A{
		"a long key for a map that may fit in the debugger variables window pane": A{B: "a long value for a key in a map that may fit in the debugger variables window pane"},
	}
	now := time.Now()
	fmt.Println(s, m, i, iptr, f, a, ar, ml, now)

	// this part is commented because it break yaegi somehow

	// mt := map[string]*time.Time{
	// 	"now":  &now,
	// 	"then": nil,
	// }
	// fmt.Println(s, m, i, f, a, ar, ml, now, mt)
}
