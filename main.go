package main

import (
	"fmt"
	"sync"
)

type s struct {
	i int
}

func main() {
	var wg sync.WaitGroup
	sli := []string{"a", "b", "c", "d"}
	for _, v := range sli {
		go func() {
			ii := v + "f"
			wg.Add(1)
			fmt.Println("===", ii)
			wg.Done()
		}()
	}
	wg.Wait()
}

func getNum(b byte) int {
	return int(b - '0')
}
