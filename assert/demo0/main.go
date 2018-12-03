package main

import (
	"fmt"
)

func assert(i interface{}) {
	v, ok := i.(int)
	fmt.Println(v, ok)
}
func main() {
	var s interface{} = 56
	assert(s)
	var i interface{} = "Steven Paul"
	assert(i)
}
