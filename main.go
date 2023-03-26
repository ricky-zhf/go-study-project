package main

import (
	"fmt"
	"sync"
	"time"
)

type sayer interface { //define a interface
	say()
}

//define a struct realizing interface sayer
type cat struct{}
type dog struct{}

func (c *cat) say() {

}
func (d dog) say() {

}

//define a function that can accept the variable implementing the interface sayer
func hit(arg sayer) {
	arg.say()
}

var s1 ss

type ss struct {
	m map[int]int
}

var mmm sync.Map

var Mm = map[int]int{
	1: 1,
}

func main() {
	ch := make(chan int, 5)
	go test(ch)

	for {
		if v, ok := <-ch; ok {
			fmt.Println("get val: ", v, ok)
		} else {
			fmt.Println("done")
			break
		}
	}
}

func test(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}

	time.Sleep(10 * time.Second)
	close(ch)
	fmt.Println("put in done")
}

func t() {
	for i := 0; i < 10; i++ {
		go fmt.Println(i)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("---")
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}

	time.Sleep(3 * time.Second)
	fmt.Println("---")
	for i := 0; i < 20; i++ {
		go func() {
			j := i
			fmt.Println(j)
		}()
	}
}
