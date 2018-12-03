package main

import "fmt"

//https://mp.weixin.qq.com/s/2naJiu3GeNeT8gVyhMo1fA
//用接口来模拟多态
type Fruitable interface {
	eat()
}

type Fruit struct {
	Name      string // 属性变量
	Fruitable        // 匿名内嵌接口变量
}

func (f Fruit) want() {
	fmt.Printf("I like ")
	f.eat() // 外结构体会自动继承匿名内嵌变量的方法
}

type Apple struct{}

func (a Apple) eat() {
	fmt.Println("eating apple")
}

type Banana struct{}

func (b Banana) eat() {
	fmt.Println("eating banana")
}

func main() {
	var f1 = Fruit{"Apple", Apple{}}
	var f2 = Fruit{"Banana", Banana{}}
	f1.want()
	f2.want()
}
