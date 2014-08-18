/*
    iputils - Package with utility functions related to working with IPs
    (c) 2014 Zencube, Geekonaut
    Author: Martin Naumann
    License: MIT
*/
package iputils

import (
	"net"
	"net/http"
	"strings"
)

// Takes an IP address as a string and checks
// if it is a valid IPv4 address and if it belongs in a private IP range
func IsIpPrivate(ipAddr string) bool {
	localIp := net.ParseIP(ipAddr)
	privateNetMasks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	if localIp.To4() == nil {
		return false
	}

	ipIsPrivate := false
	for _, netmask := range privateNetMasks {
		_, theNet, _ := net.ParseCIDR(netmask)
		if theNet.Contains(localIp) {
			ipIsPrivate = true
			break
		}
	}

	return ipIsPrivate
}

// Extracts the remote address either from an X-Forwarded-For header
// or if such header is not present from the RemoteAddr in the request itself
// It returns only the IP, not a port that is potentially included in the original request
func GetRemoteIpFromRequest(request *http.Request) string {
	var remoteIp string
	remoteIp = request.Header.Get("X-Forwarded-For")
	if remoteIp == "" {
		remoteIp = request.RemoteAddr
	}
	remoteIp = strings.Split(remoteIp, ":")[0]

    return remoteIp
}