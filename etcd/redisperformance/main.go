package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func main() {
	var tmp int

	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()

	go func() {
		for {
			username, err := redis.String(c.Do("GET", "mykey"))
			if err != nil {
				fmt.Println("redis get failed:", err)
			} else {
				tmp = tmp + 1
				fmt.Printf("Get mykey: %v \n", username)
			}
		}

	}()
	time.Sleep(time.Second)
	fmt.Println(tmp)

}
