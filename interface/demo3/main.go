package main

import (
	"fmt"
	"strconv" //for conversions to and from string
)

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
	//fmt.Println("This Human is : ", Bob)
	fmt.Printf("This Human is : %v\n", Bob)
}
