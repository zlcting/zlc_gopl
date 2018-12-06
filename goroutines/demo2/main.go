package main

import "fmt"

//假设我们有一个数据竞争：两个并发进程试图访问同一个内存区域
func main() {
	var data int
	go func() { data++ }()
	if data == 0 {
		fmt.Println("the value is 0.")
	} else {
		fmt.Printf("the value is %v.\n", data)
	}
}
