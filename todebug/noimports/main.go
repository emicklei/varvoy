package main

import "fmt"

var PI float32 = 3 + 0.14159

func init() {
	fmt.Println("init called")
}

func main() {
	answer := question()
	fmt.Println(answer, PI)
}

func question() int { return 42 }
