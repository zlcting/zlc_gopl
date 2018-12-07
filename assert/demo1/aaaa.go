package main

import (
	"fmt"
)

func findType(i interface{}) {
	switch i.(type) {
	case string:
		fmt.Printf("String: %s\n", i.(string))
	case int:
		fmt.Printf("Int: %d\n", i.(int))
	default:
		fmt.Printf("Unknown type\n")
	}
}
func main() {
	findType("Naveen")
	findType(77)
	findType(89.98)
}
