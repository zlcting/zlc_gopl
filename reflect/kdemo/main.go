package main

import "reflect"
import "fmt"

func main() {
	type SomeInt int
	var s SomeInt = 42
	var t = reflect.TypeOf(s)
	var v = reflect.ValueOf(s)
	// reflect.ValueOf(s).Type() 等价于 reflect.TypeOf(s)
	fmt.Println(t == v.Type())
	fmt.Println(v.Kind() == reflect.Int) // 元类型
	// 将 Value 还原成原来的变量
	var is = v.Interface()
	fmt.Println(is.(SomeInt))
}
