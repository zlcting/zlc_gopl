package main

import "fmt"
import "time"
import "math/rand"

func send(ch chan int) {
	for {
		var value = rand.Intn(100)
		ch <- value
		fmt.Printf("send %d\n", value)
	}
}

func recv(ch chan int) {
	for {
		value := <-ch
		fmt.Printf("recv %d\n", value)
		time.Sleep(time.Second)
	}
}

func main() {
	var ch = make(chan int, 1)
	// 子协程循环读
	go recv(ch)
	// 主协程循环写
	send(ch)
}
