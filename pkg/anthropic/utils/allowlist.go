package utils

import (
	"fmt"
	"net"
)

// AllowedIPRanges are the IP ranges that Anthropics services live on.
var AllowedIPRanges = []string{
	"160.79.104.0/23", // IPv4 range
	"2607:6bc0::/48",  // IPv6 range
}

// IsIPAllowed checks if the provided IP address is within the allowed IP ranges.
func IsIPAllowed(ipStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	for _, cidrStr := range AllowedIPRanges {
		_, ipNet, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return false, err
		}

		if ipNet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
}
