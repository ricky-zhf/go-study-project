package main

import "fmt"

func main() {
	c := make(chan int, 1)
	x, ok := <-c
	if ok {
		fmt.Println(x)
	} else {
		fmt.Println("--")
	}
}
