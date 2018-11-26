package main

import "fmt"

//函数作为一等公民
func main() {
	s := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}
	fmt.Printf("原来素组             ：%#4v\n", s)

	//对slice的元素进行平方 把函数作为参数，传递给其它函数
	result, e := operateSliceByFunc(s, func(n int) (int, error) { return n * n, nil })

	if e == nil {
		fmt.Printf("对每个元素进行平方运算：%#4v\n", result)
	}

	s = make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	//对slice的元素进行加10操作
	result, e = operateSliceByFunc(s, func(n int) (int, error) { return n + 10, nil })

	if e == nil {
		fmt.Printf("对每个元素进行加10运算：%#4v\n", result)
	}

	//在其它函数中返回函数
	a := simple()
	fmt.Println(a(60, 7))

}

func operateSliceByFunc(s []int, f func(n int) (int, error)) ([]int, error) {
	for i := 0; i < len(s); i++ {
		result, e := f(s[i])
		if e == nil {
			s[i] = result
		} else {
			return make([]int, 0, 0), e
		}
	}
	return s, nil
}

//赋值返回 f()函数
func simple() func(a, b int) int {
	f := func(a, b int) int {
		return a + b
	}
	return f
}
