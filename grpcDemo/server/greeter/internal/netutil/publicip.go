package netutil

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// publicIPServices is a list of fallback HTTP endpoints that return the caller's public IP.
var publicIPServices = []string{
	"https://ifconfig.me",
	"https://api.ipify.org",
	"https://icanhazip.com",
	"https://checkip.amazonaws.com",
}

// DiscoverPublicIP returns the machine's public IP address.
//
// Priority:
//
//	1. ADVERTISED_IP env var (explicit override)
//	2. Query external HTTP services (multiple fallbacks)
//	3. Fall back to preferred outbound IP (local interface)
func DiscoverPublicIP() (string, error) {
	// 1. env override
	if ip := os.Getenv("ADVERTISED_IP"); ip != "" {
		return ip, nil
	}

	// 2. try external services
	client := &http.Client{Timeout: 3 * time.Second}
	for _, svc := range publicIPServices {
		if ip, err := queryPublicIP(client, svc); err == nil && net.ParseIP(ip) != nil {
			return ip, nil
		}
	}

	// 3. fallback: outbound IP from local interface
	return getOutboundIP()
}

func queryPublicIP(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s returned %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

// getOutboundIP gets the preferred outbound IP of this machine.
func getOutboundIP() (string, error) {
	conn, err := net.DialTimeout("udp", "8.8.8.8:80", 2*time.Second)
	if err != nil {
		return "", fmt.Errorf("detect outbound IP: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
