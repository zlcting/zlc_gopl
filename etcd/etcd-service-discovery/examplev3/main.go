package main

import (
	"flag"
	"fmt"
	"zlc_gopl/etcd/etcd-service-discovery/discoveryv3"
)

func main() {
	var role = flag.String("role", "", "master | worker")
	flag.Parse()
	endpoints := []string{"http://127.0.0.1:2379"}
	if *role == "master" {
		master := discoveryv3.NewMaster(endpoints)
		master.WatchWorkers()
	} else if *role == "worker" {
		worker := discoveryv3.NewWorker("localhost", "127.0.0.1", endpoints)
		worker.HeartBeat()
	} else {
		fmt.Println("example -h for usage")
	}
}
