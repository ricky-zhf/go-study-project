package main

import "math"

func twoEggDrop(n int) int {
	m := make(map[s]int, 0)

	var f func(k, n int) int
	f = func(k, n int) int {
		if k == 1 {
			return n
		}
		if n == 0 {
			return 0
		}

		if _, ok := m[s{k, n}]; ok {
			return m[s{k, n}]
		}

		res := math.MaxInt64
		for i := 0; i <= n; i++ {
			res = min(res, max(f(k, n-i), f(k-1, i-1))+1)
		}
		m[s{k, n}] = res
		return res
	}

	return f(2, n)
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type s struct {
	k int
	v int
}
