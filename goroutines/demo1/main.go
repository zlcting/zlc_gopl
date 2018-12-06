package main

import "fmt"
import "time"

func main() {
	fmt.Println("run in main goroutine")
	go func() {
		fmt.Println("run in child goroutine")
		go func() {
			fmt.Println("run in grand child goroutine")
			go func() {
				fmt.Println("run in grand grand child goroutine")
			}()
		}()
	}()
	time.Sleep(time.Second)
	fmt.Println("main goroutine will quit")
}
