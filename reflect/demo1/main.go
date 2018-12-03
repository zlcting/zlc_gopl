package main

import (
	"fmt"
	"strconv"
)

type Stringer interface {
	String() string
}

func ToString(any interface{}) string {
	if v, ok := any.(Stringer); ok {
		return v.String()
	}
	switch v := any.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return "???"
}
func main() {
	value := ToString(123.345)
	fmt.Println(value)
	value = ToString(123456)
	fmt.Println(value)
}
