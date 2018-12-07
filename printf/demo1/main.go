package main

import "fmt"

func main() {
	val := []int{1, 2, 3}

	fmt.Printf("%v, %T\n", val, val) //output: [1,2,3], []int

	fmt.Printf("|%05d|\n", 1)
	fmt.Printf("|%-5d|\n", 1)
	fmt.Printf("|%5d|\n", 1234567111111)
}
