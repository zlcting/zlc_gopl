package discoveryv3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type Master struct {
	members map[string]*Member
	Client  *clientv3.Client
}

type Member struct {
	InGroup bool
	IP      string
	Name    string
	CPU     int
}

func NewMaster(endpoints []string) *Master {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}
	master := &Master{
		members: make(map[string]*Member),
		Client:  etcdClient,
	}

	return master
}

func (m *Master) WatchWorkers() {
	Client := m.Client
	wc := Client.Watch(context.Background(), "workers/", clientv3.WithPrefix())
	for v := range wc {
		for _, e := range v.Events {
			fmt.Printf("type:%v kv:%v  value:%v\n", e.Type, string(e.Kv.Key), string(e.Kv.Value))

			info := &WorkerInfo{}
			json.Unmarshal([]byte(e.Kv.Value), info)

			switch e.Type {
			case mvccpb.PUT:
				log.Println("Update worker name：", info.Name)
			case mvccpb.DELETE:
				log.Println("del worker key:", string(e.Kv.Key))
			}

		}
	}
}
