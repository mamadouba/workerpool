package tasks

import (
	"errors"
	"sort"
)

func fib(n int) int {
	if n < 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
func Fibn(arg ...interface{}) (interface{}, error) {
	n, _ := arg[0].(int)
	r := []int{}
	for i := 0; i < n; i++ {
		r = append(r, fib(i))
	}
	return r, nil
}

func Fibc(arg ...interface{}) interface{} {
	n, _ := arg[0].(int)
	c := make(chan int, n)
	go func(n int, c chan int) {
		x, y := 1, 1
		for i := 0; i < n; i++ {
			c <- x
			x, y = y, x+y
		}
		close(c)
	}(n, c)
	r := []int{}
	for v := range c {
		r = append(r, v)
	}
	return r
}

func SortList(arg ...interface{}) interface{} {
	list, ok := arg[0].([]int)
	if !ok {
		return errors.New("list of integers expected")
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return list
}
