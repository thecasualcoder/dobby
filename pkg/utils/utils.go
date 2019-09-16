package utils

import (
	"fmt"
	"net"
)

// GetOutboundIP will return a single IP of the host through which outbound traffic to internet is sent
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", fmt.Errorf("error when resolving outbound IP: %s", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
