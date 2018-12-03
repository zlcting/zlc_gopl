package main

import (
	"fmt"
	"strconv"
)

//https://wangzhezhe.github.io/2017/01/25/golang-reflection-model/
type stringer interface {
	String() string
}

func Sprint(x interface{}) string {

	switch x := x.(type) {
	case stringer:
		fmt.Println(111111111)
		fmt.Printf("%+v\n", x.String())

		return x.String()
	case string:
		return x
	case int:
		return strconv.Itoa(x)
	// ...similar cases for int16, uint32, and so on...
	case bool:
		if x {
			return "true"
		}
		return "false"
	default:
		// array, chan, func, map, pointer, slice, struct
		return "???"
	}
}

type Human struct {
	name  string
	age   int
	phone string
}

func (h Human) String() string {

	return "❰" + h.name + " - " + strconv.Itoa(h.age) + " years111 -  ✆ " + h.phone + "❱"
}

func main() {
	Bob := Human{"Bob", 39, "000-7777-XXX"}
	a := Sprint(Bob)

	fmt.Printf("%T \n", a)
	fmt.Printf("%+v\n", a)
}
