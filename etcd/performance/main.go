package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	var tmp int
	//连接etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {

		fmt.Println("connect failed, err:", err)
		return
	}
	fmt.Println("connect succ")
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	go func() {
		for {
			resp, err := cli.Get(ctx, "/testkey/fenxiang")

			if err != nil {
				fmt.Println("get failed, err:", err)
				return
			}

			for _, ev := range resp.Kvs {
				fmt.Printf("%s : %s\n", ev.Key, ev.Value)
			}
			tmp = tmp + 1
		}

	}()

	time.Sleep(time.Second)
	cancel()

	fmt.Println(tmp)
}
