package main

//send的时候会直接阻塞main goroutine
//https://mp.weixin.qq.com/s/1-MwFEBEr0uu_P4SCSgtkQ
import (
	"fmt"
)

func main() {
	goo(32)
}
func goo(s int) {

	counter := make(chan int)

	go func() {
		counter <- s
	}()

	fmt.Println(<-counter)
}
