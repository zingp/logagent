### logagent

logagent 是一个实时收集日志的并发送到kafka集群的客户端。

    1 支持多个日志实时收集
    2 支持限流功能
    3 依赖etcd, 支持动态配置收集日志

### 部署
    1 编译logagent下所有.go文件得到一个二进制，如logagent
    2 将二进制放在任何一个目录下，如agent/,并在该目录下建立conf、log目录：
        ├── conf
        │   └── app.cfg
        ├── log
        │   └── logagent.log
        └── logagent
    3 在conf下编辑配置文件app.cfg, 配置文件详解：

        etcd_addr = 10.134.123.183:2379         # etcd 地址
        etcd_timeout = 5                        # 连接etcd超时时间
        etcd_watch_key = /logagent/%s/logconfig    # etcd key 格式

        kafka_addr = 10.134.123.183:9092           # 卡夫卡地址

        thread_num = 4                             # 线程数
        log = ./log/logagent.log                   # agent的日志文件
        level = debug                              # 日志级别


### etcd value说明
`[
	{
	"service":"test_service",        
	"log_path":"/search/nginx/logs/ping-android.shouji.sogoucom_access_log",
    "topic": "nginx_log",
	"send_rate": 1000
	},
	{
	"service":"srv.android.shouji.sogou.com",
	"log_path":"/search/nginx/logs/srv.android.shouji.sogou.com_access_log",
    "topic": "nginx_log",
	"send_rate": 2000
	}
]`

    "service":"服务名称",        
    "log_path": "应该监听的日志文件",   
    "topic": "kfk topic",
    "send_rate": "日志条数限制"
