package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sipeed/picoclaw/pkg/config"
)

const Logo = "🦞"

var (
	version   = "2.0.5"
	gitCommit string
	buildTime string
	goVersion string
)

func GetConfigPath() string {
	if configPath := os.Getenv("PICOCLAW_CONFIG"); configPath != "" {
		return configPath
	}
	// v0.2.1: Use PICOCLAW_HOME if set
	var homePath string
	if picoclawHome := os.Getenv("PICOCLAW_HOME"); picoclawHome != "" {
		homePath = picoclawHome
	} else {
		home, _ := os.UserHomeDir()
		homePath = filepath.Join(home, ".picoclaw")
	}
	return filepath.Join(homePath, "config.json")
}

func LoadConfig() (*config.Config, error) {
	return config.LoadConfig(GetConfigPath())
}

// FormatVersion returns the version string with optional git commit
func FormatVersion() string {
	v := version
	if gitCommit != "" {
		v += fmt.Sprintf(" (git: %s)", gitCommit)
	}
	return v
}

// FormatBuildInfo returns build time and go version info
func FormatBuildInfo() (string, string) {
	build := buildTime
	goVer := goVersion
	if goVer == "" {
		goVer = runtime.Version()
	}
	return build, goVer
}

// GetVersion returns the version string
func GetVersion() string {
	return version
}
