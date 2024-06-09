package main

import (
	"fmt"
	"log"
)

var PI float32 = 3 + 0.14159

func init() {
	fmt.Println("init called")
}

func main() {
	answer := question()
	log.Println(answer, PI)
}

func question() int { return 42 }
