package main

import "fmt"

func send(ch chan int) {
	// ch <- 1
	// ch <- 2
	// ch <- 3
	// ch <- 4
	// close(ch)

	i := 0
	for {
		i++
		ch <- i
	}
}

func recv(ch chan int) {
	value := <-ch
	fmt.Println(value)
	value = <-ch
	fmt.Println(value)
	close(ch)
}

func main() {
	var ch = make(chan int, 4)
	go recv(ch)
	send(ch)
}
