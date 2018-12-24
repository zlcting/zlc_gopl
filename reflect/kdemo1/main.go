package main

import "fmt"
import "reflect"

func main() {
	var s int = 42
	// 反射指针类型
	var v = reflect.ValueOf(&s)
	// 要拿出指针指向的元素进行修改
	v.Elem().SetInt(43)
	fmt.Println(s)
}
