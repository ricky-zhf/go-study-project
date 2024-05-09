package main

import (
	"fmt"
)

type st struct {
	name string
}

func main() {
	sli := make([]int, 2)
	sli[0], sli[1] = 1, 2
	fmt.Printf("sli=%v|cap=%v|p=%p\n", sli, cap(sli), sli)

	sli = append(sli, 3)
	fmt.Printf("加3后，新的sli=%v|cap=%v|新的p=%p\n", sli, cap(sli), sli)
	fmt.Println()

	sli2 := append(sli, 4)
	fmt.Printf("加4后，原来sli=%v|原来p=%p\n", sli, sli)
	fmt.Printf("加4后，新的sli2=%v|新的p2=%p\n", sli2, sli2)
	fmt.Println()

	sli3 := append(sli, 5, 6, 7)
	fmt.Printf("加567后，原来sli=%v|cap=%v|原来p=%p\n", sli, cap(sli), sli)
	fmt.Printf("加567后，原来sli3=%v|cap3=%v|原来p3=%p\n", sli3, cap(sli3), sli3)
	fmt.Println()

	sli = append(sli, 8, 9, 0)
	fmt.Printf("加890后，sli=%v|p=%p\n", sli, sli)
	fmt.Printf("加890后，sli3=%v|p=%p\n", sli3, sli3)
}
