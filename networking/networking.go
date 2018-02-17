package networking

import (
	"fmt"
	"net"
	"strings"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("Outbound IP couldn't be fetched; returning 127.0.0.1:0")
		return "127.0.0.1:0"
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	index := strings.LastIndex(localAddr, ":")
	return localAddr[0:index]
}
