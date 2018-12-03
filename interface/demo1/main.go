package main

import "fmt"

// 可以闻
type Smellable interface {
	smell()
}

// 可以吃
type Eatable interface {
	eat()
}

// 苹果既可能闻又能吃
type Apple struct{}

func (a Apple) smell() {
	fmt.Println("apple can smell")
}

func (a Apple) eat() {
	fmt.Println("apple can eat")
}

// 花只可以闻
type Flower struct{}

func (f Flower) smell() {
	fmt.Println("flower can smell")
}

func main() {
	var s1 Smellable //interface
	var s2 Eatable
	var apple = Apple{}
	var flower = Flower{}

	s1 = apple
	s1.smell()

	s1 = flower
	s1.smell()

	s2 = apple
	s2.eat()

	//s2 = apple
	apple.eat()
	//???我们还要用接口呢？
}
