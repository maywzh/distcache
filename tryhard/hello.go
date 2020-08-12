package main

import (
	"fmt"
	"sort"
)

func main() {
	s := [6]int {2,3,5,7,11,13}
	idx := sort.Search(6, func(i int) bool{
		return s[i] == 13
	})
	fmt.Print(idx)
}