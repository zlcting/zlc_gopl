package main

import "fmt"

func main() {
	var a = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var b = a[4:7]

	c := map[string]int{
		"alice":   31,
		"charlie": 34,
	}
	d := make(map[string]int)

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}
	