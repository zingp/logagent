package logagent

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

type AppConfig struct {
	EtcdAddr     string
	EtcdTimeOut  int
	EtcdWatchKey string

	KafkaAddr string

	ThreadNum int
	LogFile   string
	LogLevel  string
}

var appConf = &AppConfig{}

func initConfig(file string) (err error) {
	conf, err := config.NewConfig("ini", file)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	appConf.EtcdAddr = conf.String("etcd_addr")
	appConf.EtcdTimeOut = conf.DefaultInt("etcd_timeout", 5)
	appConf.EtcdWatchKey = conf.String("etcd_watch_key")

	appConf.KafkaAddr = conf.String("kafka_addr")

	appConf.ThreadNum = conf.DefaultInt("thread_num", 4)
	appConf.LogFile = conf.String("log")
	appConf.LogLevel = conf.String("level")
	return
}
