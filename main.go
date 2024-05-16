package main

import (
	"fmt"
	"sort"
)

func main() {
	frequent := topKFrequent([]int{1, 2, 2, 3, 2, 12, 1, 11, 2, 2, 1, 2, 3, 3, 12, 31, 3, 312, 1, 31, 1, 3, 3, 2, 2, 2}, 5)
	fmt.Println(frequent)
}

type freq struct {
	num      int
	frequent int
}

func topKFrequent(nums []int, k int) []int {
	freqMap := make(map[int]int) // num : freq
	for _, v := range nums {
		freqMap[v]++
	}
	fmt.Println(freqMap)
	var res []freq
	for k, v := range freqMap {
		res = append(res, freq{num: k, frequent: v})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].frequent > res[j].frequent
	})
	fmt.Println(res)
	var resNum []int
	for _, v := range res {
		resNum = append(resNum, v.num)
	}
	return resNum[:k+1]
}
