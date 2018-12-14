package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}

	_, err = c.Do("SET", "password", "123456", "EX", "10")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	time.Sleep(11 * time.Second)
	password, err := redis.String(c.Do("GET", "password"))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Got password %v \n", password)
	}
	defer c.Close()
}
