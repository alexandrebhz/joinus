package ssrf

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var blockedCIDRs []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	} {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil {
			blockedCIDRs = append(blockedCIDRs, network)
		}
	}
}

// ValidatePublicHTTPURL ensures raw is an http(s) URL that does not resolve to private,
// link-local, loopback, or metadata IP addresses.
func ValidatePublicHTTPURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("URL must have a host")
	}

	lowerHost := strings.ToLower(strings.TrimSuffix(host, "."))
	if lowerHost == "localhost" || strings.HasSuffix(lowerHost, ".localhost") {
		return fmt.Errorf("localhost URLs are not allowed")
	}

	if ip := net.ParseIP(host); ip != nil {
		if isBlockedIP(ip) {
			return fmt.Errorf("URL resolves to blocked IP address: %s", ip.String())
		}
		return nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("failed to resolve host: %w", err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("no IP addresses resolved for host")
	}

	for _, ip := range ips {
		if isBlockedIP(ip) {
			return fmt.Errorf("URL resolves to blocked IP address: %s", ip.String())
		}
	}

	return nil
}

func isBlockedIP(ip net.IP) bool {
	ip = ip.To16()
	if ip == nil {
		return true
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}

	for _, network := range blockedCIDRs {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
