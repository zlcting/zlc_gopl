package main

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		println(err)
		return
	}

	for index := 0; index < 100000; index++ {
		c.Do("pfadd", "codehole", "user"+strconv.Itoa(index))
	}

	if err != nil {
		println(err)
		return
	}

	println("finish")
}
