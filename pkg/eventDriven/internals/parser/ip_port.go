package parser

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseIpAndPort(serverAddress string) (ip net.IP, port int, err error) {
	ip = net.ParseIP("127.0.0.1")
	port = 9040
	addr := strings.Split(serverAddress, ":")
	if len(addr) == 1 {
		ip = net.ParseIP(serverAddress)
	} else if len(addr) >= 2 {
		ip = net.ParseIP(serverAddress)
		tempPort, err := strconv.ParseInt(addr[len(addr)-1], 10, 32)
		if err != nil {
			return ip, port, fmt.Errorf("failed parse port %w", err)
		}

		port = int(tempPort)

	}
	return ip, port, nil
}
