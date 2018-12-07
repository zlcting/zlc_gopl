package main

import "fmt"

type Describer interface {
	Describe()
}
type St string

func (s St) Describe() {
	fmt.Println("被调用le!")
}

func findType(i interface{}) {

	switch v := i.(type) {
	case Describer:
		v.Describe()
	case string:
		fmt.Println("String 变量")
	default:
		fmt.Printf("unknown type\n")
	}
}

func main() {
	findType("Naveen")
	st := St("我的字符串")
	findType(st)
}
