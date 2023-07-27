package mask

import (
	"net"
	"strings"
)

// list of private subnets
var privateMasks, _ = ParseFrom([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"fc00::/7",
})

// ParseFrom  converts a list of subnets' string to a list of net.IPNet.
func ParseFrom(ips []string) (masks []net.IPNet, err error) {
	for _, cidr := range ips {
		var network *net.IPNet
		_, network, err = net.ParseCIDR(cidr)
		if err != nil {
			return
		}
		masks = append(masks, *network)
	}
	return
}

// CheckIpInMasks checks if a net.IP is in a list of net.IPNet
func CheckIpInMasks(ip net.IP, masks []net.IPNet) bool {
	for _, mask := range masks {
		if mask.Contains(ip) {
			return true
		}
	}
	return false
}

// isPublicIP returns true if the given IP can be routed on the Internet.
func isPublicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() {
		return false
	}
	return !CheckIpInMasks(ip, privateMasks)
}

// GetPublicIpFrom Parse parses the value of the X-Forwarded-For Header
// and returns the IP address.
func GetPublicIpFrom(ipList string) string {
	for _, ip := range strings.Split(ipList, ",") {
		ip = strings.TrimSpace(ip)
		if IP := net.ParseIP(ip); IP != nil && isPublicIP(IP) {
			return ip
		}
	}
	return ""
}
