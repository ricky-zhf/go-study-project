package main

import "fmt"

type st struct {
	a  string
	am map[string]string
}

func main() {
	am := map[string]string{
		"1": "1",
	}
	s1 := st{
		a:  "123",
		am: am,
	}
	am2 := map[string]string{
		"2": "2",
	}
	s2 := st{
		a:  "456",
		am: am2,
	}
	fmt.Printf("1:%p\n", &s1)
	fmt.Printf("2:%p\n", &s2)
	s1 = s2
	fmt.Printf("1:%p\n", &s1)
	fmt.Printf("2:%p\n", &s2)

	fmt.Printf("1:%p\n", s1.am)
	fmt.Printf("2:%p\n", s2.am)
}
