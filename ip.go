package main

import (
	"fmt"
	"net"
)

// 获取本机IP地址列表
var ipArray []string

func getLocalIP() (ips []string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("get ip interfaces error:", err)
		return
	}

	for _, i := range ifaces {
		addrs, errRet := i.Addrs()
		if errRet != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if ip.IsGlobalUnicast() {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return
}
