package string

import "fmt"

func StringsMethod() {
	ss := [...]int{1, 2, 3, 4, 5}
	s := ss[:]
	s = append(s[:1], s[3:]...)
	fmt.Println(s)
	fmt.Println(ss)

}
