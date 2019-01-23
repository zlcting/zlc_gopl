package discoveryv3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type Worker struct {
	Name   string
	IP     string
	Client *clientv3.Client
	kv     clientv3.KV
}

type WorkerInfo struct {
	Name string
	IP   string
	CPU  int
}

func NewWorker(name, IP string, endpoints []string) *Worker {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	w := &Worker{
		Name:   name,
		IP:     IP,
		Client: etcdClient,
	}
	return w
}

func (w *Worker) HeartBeat() {
	Client := w.Client
	kv := clientv3.NewKV(Client)

	info := &WorkerInfo{
		Name: w.Name,
		IP:   w.IP,
		CPU:  runtime.NumCPU(),
	}

	key := "workers/" + w.Name
	value, _ := json.Marshal(info)
	//设置租约
	lease := clientv3.NewLease(Client)
	leaseResp, err := lease.Grant(context.TODO(), 3)

	if err != nil {
		fmt.Printf("设置租约时间失败:%s\n", err.Error())
	}

	//put操作
	leaseID := leaseResp.ID
	_, err = kv.Put(context.TODO(), key, string(value), clientv3.WithLease(leaseID))

	//设置续租
	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseID)

	for {
		select {
		case leaseKeepResp := <-leaseRespChan:
			if leaseKeepResp == nil {
				fmt.Printf("已经关闭续租功能\n")
				break
			} else {
				fmt.Printf("续租成功\n")
				goto END
			}
		}
	END:
		time.Sleep(500 * time.Millisecond)
	}

	cancelFunc()

}
