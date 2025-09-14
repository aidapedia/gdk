package util

import "net"

func CheckSubnet(ip, subnet string) bool {
	ipAddr := net.ParseIP(ip)
	_, subnetNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}
	return subnetNet.Contains(ipAddr)
}
