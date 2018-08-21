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

	cli, err := client.New(client.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("connect etcd error:", err)
		return
	}
	logs.Debug("init etcd success")
	// defer cli.Close()   //这里千万不能关闭
	// 生成etcd key
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

	// 第一次运行主动从etcd拉取配置
	for _, key := range etcdKeys {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, key)
		cancel()
		if err != nil {
			fmt.Println("get etcd key failed, error:", err)
			continue
		}

		for _, ev := range resp.Kvs {
			// 返回的类型不是string,需要强制转换
			// fmt.Printf("etcd conf is %s", string(ev.Value))
			confChan <- string(ev.Value)
			fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
		}
	}
	
	waitGroup.Add(1)
	go etcdWatch(etcdKeys)
	return
}

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

//GetEtcdConfChan is func get etcd conf
func GetEtcdConfChan() chan string {
	return confChan
}
