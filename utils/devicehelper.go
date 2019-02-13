package utils

import (
	"logging"
	"net"
)

// GetDeviceIpAddress returns current IP address of the device
func GetDeviceIpAddress() (string, error) {
	ipAddress := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logging.AddError("Retrieving host's IP address failed.", err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP.String()
			}
		}
	}

	return ipAddress, err
}
