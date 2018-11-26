package main

import "fmt"

func main() {
	resule := getNum()
	fmt.Println(resule())
	fmt.Println(resule())
}
func getNum() func() int {
	i := 1
	return func() int {
		i += 1
		return i
	}
}
