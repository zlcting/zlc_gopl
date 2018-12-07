package main

import (
	"fmt"
	"time"
)

func send(ch chan int) {
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	close(ch)

}

func recv(ch chan int) {
	value := <-ch
	fmt.Println(value)
	value = <-ch
	fmt.Println(value)
}

func main() {
	var ch = make(chan int, 4)
	go recv(ch)
	send(ch)
	time.Sleep(1 * time.Second)
}
