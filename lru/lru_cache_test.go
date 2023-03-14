package lru

import (
	"fmt"
	"testing"
)

func TestInitLRU(t *testing.T) {
	l := InitLRU(3)
	l.Put(1, 1)
	l.Put(2, 2)
	l.Put(3, 3)
	temp := l.head.Next
	fmt.Println("测试put")
	for temp != nil {
		fmt.Println(temp.Key, "-", temp.Value)
		temp = temp.Next
	}
	l.Get(2)
	temp = l.head.Next
	fmt.Println("测试get")
	for temp != nil {
		fmt.Println(temp.Key, "-", temp.Value)
		temp = temp.Next
	}

	l.Put(4, 4)
	temp = l.head.Next
	fmt.Println("测试超限")
	for temp != nil {
		fmt.Println(temp.Key, "-", temp.Value)
		temp = temp.Next
	}

}
