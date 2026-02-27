package utils

import (
	"errors"
	"fmt"
	"net"

	"github.com/peterouob/seckill_service/utils/logs"
)

func getLocalIP() (ipv4 string) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isNetIp bool
		err     error
	)

	addrs, err = net.InterfaceAddrs()
	logs.HandelError("GetLocalIP error ", err)
	for _, addr = range addrs {
		if ipNet, isNetIp = addr.(*net.IPNet); isNetIp && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	err = errors.New("not found ip from get local ip")
	logs.HandelError("GetLocalIP error ", err)
	return
}

func FormatIP(port string) string {
	localP := getLocalIP()
	return fmt.Sprintf("%s:%s", localP, port)
}
