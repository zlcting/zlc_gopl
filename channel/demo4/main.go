package main

import "fmt"
import "time"

func send(ch1 chan int, ch2 chan int) {
	i := 0
	for {
		i++
		select {
		case ch1 <- i:
			fmt.Printf("send ch1 %d\n", i)
		case ch2 <- i:
			fmt.Printf("send ch2 %d\n", i)
		default:
		}
	}
}

func recv(ch chan int, gap time.Duration, name string) {
	for v := range ch {
		fmt.Printf("receive %s %d\n", name, v)
		time.Sleep(gap)
	}
}

func main() {
	// 无缓冲通道
	var ch1 = make(chan int)
	var ch2 = make(chan int)
	// 两个消费者的休眠时间不一样，名称不一样
	go recv(ch1, time.Second, "ch1")
	go recv(ch2, 2*time.Second, "ch2")
	send(ch1, ch2)
}

//------------
// send ch1 27
// send ch2 28
// receive ch1 27
// receive ch2 28
// send ch1 6708984
// receive ch1 6708984
// send ch2 13347544
// send ch1 13347775
// receive ch2 13347544
// receive ch1 13347775
// send ch1 20101642
// receive ch1 20101642
// send ch2 26775795
// receive ch2 26775795
