package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	client "github.com/coreos/etcd/clientv3"
)

var (
	confChan  = make(chan string, 10)
	cli       *client.Client
	waitGroup sync.WaitGroup
)

func initEtcd(addr []string, keyFormat string, timeout time.Duration) (err error) {
	// init a global var cli and can not close
	cli, err = client.New(client.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("connect etcd error:", err)
		return
	}
	logs.Debug("init etcd success")
	// defer cli.Close()   //can not close

	var etcdKeys []string
	ips, err := getLocalIP()
	if err != nil {
		fmt.Println("get local ip error:", err)
		return
	}
	for _, ip := range ips {
		key := fmt.Sprintf(keyFormat, ip)
		etcdKeys = append(etcdKeys, key)
	}

	// first, pull conf from etcd
	for _, key := range etcdKeys {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, key)
		cancel()
		if err != nil {
			fmt.Println("get etcd key failed, error:", err)
			continue
		}

		for _, ev := range resp.Kvs {
			// return result is not string
			confChan <- string(ev.Value)
			fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
		}
	}

	waitGroup.Add(1)
	// second, start a goroutine to watch etcd
	go etcdWatch(etcdKeys)
	return
}

// watch etcd
func etcdWatch(keys []string) {
	defer waitGroup.Done()

	var watchChans []client.WatchChan
	for _, key := range keys {
		rch := cli.Watch(context.Background(), key)
		watchChans = append(watchChans, rch)
	}

	for {
		for _, watchC := range watchChans {
			select {
			case wresp := <-watchC:
				for _, ev := range wresp.Events {
					confChan <- string(ev.Kv.Value)
					logs.Debug("etcd key = %s , etcd value = %s", ev.Kv.Key, ev.Kv.Value)
				}
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

//GetEtcdConfChan is func get etcd conf add to chan
func GetEtcdConfChan() chan string {
	return confChan
}
