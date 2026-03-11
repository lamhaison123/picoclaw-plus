package server

import (
	"net"
	"os"
	"path/filepath"
)

func DefaultConfigPath() string {
	// v0.2.1: Use PICOCLAW_HOME if set
	var homePath string
	if picoclawHome := os.Getenv("PICOCLAW_HOME"); picoclawHome != "" {
		homePath = picoclawHome
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "config.json"
		}
		homePath = filepath.Join(home, ".picoclaw")
	}
	return filepath.Join(homePath, "config.json")
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}
