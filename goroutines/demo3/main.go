package main

import (
	"fmt"
	"sync"
)

func main() {
	var memoryAccess sync.Mutex //1 这里我们添加一个变量，它允许我们的代码同步对数据变量内存的访问。
	var value int
	go func() {
		memoryAccess.Lock() //2 在这里我们声明，除非解锁，否则我们的goroutine应该独占访问此内存。
		value++
		memoryAccess.Unlock() //3 在这里，我们声明这个对该内存的访问已经完成了。
	}()

	memoryAccess.Lock() //4 在这里，我们再次声明接下来的条件语句应该独占访问数据变量的内存。
	if value == 0 {
		fmt.Printf("the value is %v.\n", value)
	} else {
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock() //5 在这里，我们声明对内存的访问已经完成。
}
