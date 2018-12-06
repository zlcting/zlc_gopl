package main

import "fmt"

func main() {
	//多维map的声明与实现方法
	//方法1 初始化一个空的多维映射
	mainMapA := map[string]map[string]string{}
	subMapA := map[string]string{"A_Key_1": "A_SubValue_1", "A_Key_2": "A_SubValue_2"}
	mainMapA["MapA"] = subMapA
	fmt.Println("MultityMapA")
	for keyA, valA := range mainMapA {
		for subKeyA, subValA := range valA {
			fmt.Printf("mapName=%s	Key=%s	Value=%s\n", keyA, subKeyA, subValA)
		}
	}

}
