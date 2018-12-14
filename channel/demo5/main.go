package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

//
// 模拟流水线作业：
// 1. 一个字符串全部改成小写
// 2. 在该字符串前面添加HELLO: 字样
// 3. 统计字数，并打印字符串
//

func starter(ch1 chan string) {
	for i := 0; i < 10; i++ {
		v := fmt.Sprintf("HoVer: %d", i)
		fmt.Fprintf(os.Stderr, "starter send: %s\n", v)
		ch1 <- v
		time.Sleep(1e9)
	}
}

func lowString(ch1 chan string, ch2 chan string) {
	for {
		str := <-ch1
		res := strings.ToLower(str)
		ch2 <- res
	}
}

func addString(ch2 chan string, ch3 chan string) {
	for {
		str := <-ch2
		res := fmt.Sprintf("HELLO: %s", str)
		ch3 <- res
	}
}

func getResult(ch3 chan string) {
	for {
		v := <-ch3
		len := len(v)
		fmt.Fprintf(os.Stderr, "len=%d, result=%s\n", len, v)
	}
}

func main() {
	ch := make(chan string)
	go starter(ch)

	ch1 := make(chan string)
	go lowString(ch, ch1)

	ch2 := make(chan string)
	go addString(ch1, ch2)
	go getResult(ch2)

	time.Sleep(10e9)
}
