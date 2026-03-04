package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
)

type SafetyLevel int

const (
	SafetyStrict     SafetyLevel = 0 // Maximum protection, blocks most dangerous commands
	SafetyModerate   SafetyLevel = 1 // Balanced protection, allows common dev operations
	SafetyPermissive SafetyLevel = 2 // Minimal protection, only blocks catastrophic commands
	SafetyOff        SafetyLevel = 3 // No safety checks (use with extreme caution)
)

type ExecTool struct {
	workingDir          string
	timeout             time.Duration
	denyPatterns        []*regexp.Regexp
	warningPatterns     []*regexp.Regexp
	allowPatterns       []*regexp.Regexp
	customAllowPatterns []*regexp.Regexp
	restrictToWorkspace bool
	safetyLevel         SafetyLevel
}

var (
	// catastrophicPatterns - Always blocked regardless of safety level
	// These can cause immediate, irreversible system damage
	catastrophicPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\brm\s+-[rf]{1,2}\s+/\s*$`),                 // rm -rf / (root deletion)
		regexp.MustCompile(`\bdd\s+if=/dev/(zero|random)\s+of=/dev/sd`), // Disk wipe
		regexp.MustCompile(`\b(format|mkfs)\s+[A-Z]:\s*$`),              // Format entire drive
		regexp.MustCompile(`:\(\)\s*\{.*:\|:.*\};\s*:`),                 // Fork bomb
		regexp.MustCompile(`>\s*/dev/(sd[a-z]|nvme\d+n\d+)\s*$`),        // Write to raw disk
		regexp.MustCompile(`\bchmod\s+-R\s+000\s+/`),                    // Recursive permission removal on root
		regexp.MustCompile(`\b(shutdown|reboot|poweroff)\s+-[fh]\b`),    // Force shutdown/reboot
	}

	// strictPatterns - Blocked in Strict mode only
	// Common dangerous operations that might be needed in dev environments
	strictPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\brm\s+-[rf]{1,2}\b`),
		regexp.MustCompile(`\bdel\s+/[fq]\b`),
		regexp.MustCompile(`\brmdir\s+/s\b`),
		regexp.MustCompile(`\b(format|mkfs|diskpart)\b\s`),
		regexp.MustCompile(`\bdd\s+if=`),
		regexp.MustCompile(`>\s*/dev/(sd[a-z]|hd[a-z]|vd[a-z]|xvd[a-z]|nvme\d|mmcblk\d|loop\d|dm-\d|md\d|sr\d|nbd\d)`),
		regexp.MustCompile(`\b(shutdown|reboot|poweroff)\b`),
		regexp.MustCompile(`\$\([^)]+\)`),
		regexp.MustCompile(`\$\{[^}]+\}`),
		regexp.MustCompile("`[^`]+`"),
		regexp.MustCompile(`\|\s*sh\b`),
		regexp.MustCompile(`\|\s*bash\b`),
		regexp.MustCompile(`\bsudo\b`),
		regexp.MustCompile(`\bchmod\s+[0-7]{3,4}\b`),
		regexp.MustCompile(`\bchown\b`),
		regexp.MustCompile(`\bpkill\b`),
		regexp.MustCompile(`\bkillall\b`),
		regexp.MustCompile(`\bkill\s+-[9]\b`),
		regexp.MustCompile(`\bcurl\b.*\|\s*(sh|bash)`),
		regexp.MustCompile(`\bwget\b.*\|\s*(sh|bash)`),
		regexp.MustCompile(`\beval\b`),
		regexp.MustCompile(`\bsource\s+.*\.sh\b`),
	}

	// moderatePatterns - Blocked in Strict and Moderate modes
	// Potentially dangerous but commonly needed for system administration
	moderatePatterns = []*regexp.Regexp{
		regexp.MustCompile(`\brm\s+-[rf]{1,2}\s+/`),         // rm -rf on root paths
		regexp.MustCompile(`\bdd\s+if=/dev/(zero|random)`),  // dd with dangerous sources
		regexp.MustCompile(`>\s*/dev/(sd[a-z]|nvme\d)`),     // Write to block devices
		regexp.MustCompile(`\bcurl\b.*\|\s*(sh|bash)`),      // Pipe curl to shell
		regexp.MustCompile(`\bwget\b.*\|\s*(sh|bash)`),      // Pipe wget to shell
		regexp.MustCompile(`\bchmod\s+-R\s+[0-7]{3,4}\s+/`), // Recursive chmod on root
	}

	// warningPatterns - Generate warnings but don't block (in Moderate/Permissive modes)
	// Operations that should be done carefully but are often legitimate
	warningPatterns = []*regexp.Regexp{
		regexp.MustCompile(`\brm\s+-[rf]`),
		regexp.MustCompile(`\bgit\s+push\s+(-f|--force)`),
		regexp.MustCompile(`\bgit\s+reset\s+--hard`),
		regexp.MustCompile(`\bdocker\s+system\s+prune\s+-a`),
		regexp.MustCompile(`\bnpm\s+install\s+-g\b`),
		regexp.MustCompile(`\bpip\s+install\s+--user\b`),
		regexp.MustCompile(`\bapt\s+(install|remove|purge)\b`),
		regexp.MustCompile(`\byum\s+(install|remove)\b`),
		regexp.MustCompile(`\bdnf\s+(install|remove)\b`),
	}

	// absolutePathPattern matches absolute file paths in commands (Unix and Windows).
	absolutePathPattern = regexp.MustCompile(`[A-Za-z]:\\[^\\\"']+|/[^\s\"']+`)

	// safePaths are kernel pseudo-devices that are always safe to reference in
	// commands, regardless of workspace restriction. They contain no user data
	// and cannot cause destructive writes.
	safePaths = map[string]bool{
		"/dev/null":    true,
		"/dev/zero":    true,
		"/dev/random":  true,
		"/dev/urandom": true,
		"/dev/stdin":   true,
		"/dev/stdout":  true,
		"/dev/stderr":  true,
	}
)

func NewExecTool(workingDir string, restrict bool) (*ExecTool, error) {
	return NewExecToolWithConfig(workingDir, restrict, nil)
}

func NewExecToolWithConfig(workingDir string, restrict bool, config *config.Config) (*ExecTool, error) {
	denyPatterns := make([]*regexp.Regexp, 0)
	customAllowPatterns := make([]*regexp.Regexp, 0)
	safetyLevel := SafetyModerate // Default to moderate

	if config != nil {
		execConfig := config.Tools.Exec

		// Determine safety level from config
		if execConfig.SafetyLevel != "" {
			switch strings.ToLower(execConfig.SafetyLevel) {
			case "strict":
				safetyLevel = SafetyStrict
			case "moderate":
				safetyLevel = SafetyModerate
			case "permissive":
				safetyLevel = SafetyPermissive
			case "off":
				safetyLevel = SafetyOff
				fmt.Println("⚠️  WARNING: Safety checks are DISABLED. All commands will be allowed!")
			default:
				fmt.Printf("Unknown safety level %q, using 'moderate'\n", execConfig.SafetyLevel)
			}
		}

		// Build deny patterns based on safety level
		if safetyLevel != SafetyOff {
			// Always include catastrophic patterns
			denyPatterns = append(denyPatterns, catastrophicPatterns...)

			switch safetyLevel {
			case SafetyStrict:
				denyPatterns = append(denyPatterns, strictPatterns...)
				denyPatterns = append(denyPatterns, moderatePatterns...)
			case SafetyModerate:
				denyPatterns = append(denyPatterns, moderatePatterns...)
			case SafetyPermissive:
				// Only catastrophic patterns
			}
		}

		// Legacy support for EnableDenyPatterns
		if !execConfig.EnableDenyPatterns && safetyLevel == SafetyModerate {
			fmt.Println("⚠️  Warning: deny patterns are disabled via config. Consider using safetyLevel='off' instead.")
			denyPatterns = make([]*regexp.Regexp, 0)
			safetyLevel = SafetyOff
		}

		// Add custom deny patterns
		if len(execConfig.CustomDenyPatterns) > 0 {
			fmt.Printf("Adding custom deny patterns: %v\n", execConfig.CustomDenyPatterns)
			for _, pattern := range execConfig.CustomDenyPatterns {
				re, err := regexp.Compile(pattern)
				if err != nil {
					return nil, fmt.Errorf("invalid custom deny pattern %q: %w", pattern, err)
				}
				denyPatterns = append(denyPatterns, re)
			}
		}

		// Add custom allow patterns (bypass all checks)
		for _, pattern := range execConfig.CustomAllowPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid custom allow pattern %q: %w", pattern, err)
			}
			customAllowPatterns = append(customAllowPatterns, re)
		}
	} else {
		// No config provided, use moderate defaults
		denyPatterns = append(denyPatterns, catastrophicPatterns...)
		denyPatterns = append(denyPatterns, moderatePatterns...)
	}

	return &ExecTool{
		workingDir:          workingDir,
		timeout:             60 * time.Second,
		denyPatterns:        denyPatterns,
		warningPatterns:     warningPatterns,
		allowPatterns:       nil,
		customAllowPatterns: customAllowPatterns,
		restrictToWorkspace: restrict,
		safetyLevel:         safetyLevel,
	}, nil
}

func (t *ExecTool) Name() string {
	return "exec"
}

func (t *ExecTool) Description() string {
	safetyDesc := ""
	switch t.safetyLevel {
	case SafetyStrict:
		safetyDesc = " [Safety: STRICT - Most dangerous commands blocked]"
	case SafetyModerate:
		safetyDesc = " [Safety: MODERATE - Catastrophic commands blocked]"
	case SafetyPermissive:
		safetyDesc = " [Safety: PERMISSIVE - Only extreme dangers blocked]"
	case SafetyOff:
		safetyDesc = " [Safety: OFF - No restrictions]"
	}
	return "Execute a shell command and return its output. Use with caution." + safetyDesc
}

func (t *ExecTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "The shell command to execute",
			},
			"working_dir": map[string]any{
				"type":        "string",
				"description": "Optional working directory for the command",
			},
		},
		"required": []string{"command"},
	}
}

func (t *ExecTool) Execute(ctx context.Context, args map[string]any) *ToolResult {
	// Add panic recovery for shell execution
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorCF("shell", "Panic in shell execution", map[string]any{
				"panic": r,
			})
		}
	}()

	command, ok := args["command"].(string)
	if !ok {
		return ErrorResult("command is required")
	}

	cwd := t.workingDir
	if wd, ok := args["working_dir"].(string); ok && wd != "" {
		if t.restrictToWorkspace && t.workingDir != "" {
			resolvedWD, err := validatePath(wd, t.workingDir, true)
			if err != nil {
				return ErrorResult("Command blocked by safety guard (" + err.Error() + ")")
			}
			cwd = resolvedWD
		} else {
			cwd = wd
		}
	}

	if cwd == "" {
		wd, err := os.Getwd()
		if err == nil {
			cwd = wd
		}
	}

	if guardError := t.guardCommand(command, cwd); guardError != "" {
		return ErrorResult(guardError)
	}

	// timeout == 0 means no timeout
	var cmdCtx context.Context
	var cancel context.CancelFunc
	if t.timeout > 0 {
		cmdCtx, cancel = context.WithTimeout(ctx, t.timeout)
	} else {
		cmdCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(cmdCtx, "powershell", "-NoProfile", "-NonInteractive", "-Command", command)
	} else {
		cmd = exec.CommandContext(cmdCtx, "sh", "-c", command)
	}
	if cwd != "" {
		cmd.Dir = cwd
	}

	prepareCommandForTermination(cmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return ErrorResult(fmt.Sprintf("failed to start command: %v", err))
	}

	done := make(chan error, 1)
	waitCtx, waitCancel := context.WithCancel(context.Background())
	defer waitCancel() // Ensure goroutine cleanup

	go func() {
		defer func() {
			// Ensure we don't panic if channel is somehow closed
			recover()
		}()
		select {
		case done <- cmd.Wait():
		case <-waitCtx.Done():
			// Parent cancelled, exit goroutine
			return
		}
	}()

	var err error
	select {
	case err = <-done:
		// Command completed normally
	case <-cmdCtx.Done():
		// Context cancelled or timeout
		_ = terminateProcessTree(cmd)
		select {
		case err = <-done:
			// Process terminated gracefully
		case <-time.After(2 * time.Second):
			// Force kill after 2 seconds
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			// Wait for goroutine to finish with timeout
			select {
			case err = <-done:
			case <-time.After(1 * time.Second):
				// Goroutine still blocked, cancel it to prevent leak
				waitCancel()
				err = fmt.Errorf("process kill timeout")
			}
		}
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		if errors.Is(cmdCtx.Err(), context.DeadlineExceeded) {
			msg := fmt.Sprintf("Command timed out after %v", t.timeout)
			return &ToolResult{
				ForLLM:  msg,
				ForUser: msg,
				IsError: true,
			}
		}
		output += fmt.Sprintf("\nExit code: %v", err)
	}

	if output == "" {
		output = "(no output)"
	}

	maxLen := 10000
	if len(output) > maxLen {
		output = output[:maxLen] + fmt.Sprintf("\n... (truncated, %d more chars)", len(output)-maxLen)
	}

	if err != nil {
		return &ToolResult{
			ForLLM:  output,
			ForUser: output,
			IsError: true,
		}
	}

	return &ToolResult{
		ForLLM:  output,
		ForUser: output,
		IsError: false,
	}
}

func (t *ExecTool) guardCommand(command, cwd string) string {
	// Safety OFF - no checks
	if t.safetyLevel == SafetyOff {
		return ""
	}

	cmd := strings.TrimSpace(command)
	lower := strings.ToLower(cmd)

	// Custom allow patterns exempt a command from all checks
	explicitlyAllowed := false
	for _, pattern := range t.customAllowPatterns {
		if pattern.MatchString(lower) {
			explicitlyAllowed = true
			break
		}
	}

	if !explicitlyAllowed {
		// Check deny patterns
		for _, pattern := range t.denyPatterns {
			if pattern.MatchString(lower) {
				return fmt.Sprintf("Command blocked by safety guard (dangerous pattern detected, safety level: %s)", t.getSafetyLevelName())
			}
		}

		// Check warning patterns (informational only, don't block)
		if t.safetyLevel == SafetyModerate || t.safetyLevel == SafetyPermissive {
			for _, pattern := range t.warningPatterns {
				if pattern.MatchString(lower) {
					// Log warning but don't block
					fmt.Printf("⚠️  Warning: Potentially dangerous command detected: %s\n", cmd)
					break
				}
			}
		}

		// Enhanced path traversal detection
		// Check for various path traversal patterns
		pathTraversalPatterns := []string{
			"../", "..\\",
			"%2e%2e/", "%2e%2e\\",
			"..%2f", "..%5c",
			"%252e%252e/", "%252e%252e\\",
		}
		for _, pattern := range pathTraversalPatterns {
			if strings.Contains(strings.ToLower(cmd), pattern) {
				return "Command blocked by safety guard (path traversal pattern detected)"
			}
		}

		// Check for command injection attempts
		injectionPatterns := []string{
			";", "|", "&", "$(",
			"`", "$(", "${",
		}
		if t.safetyLevel == SafetyStrict {
			for _, pattern := range injectionPatterns {
				if strings.Contains(cmd, pattern) {
					return fmt.Sprintf("Command blocked by safety guard (potential command injection: %s)", pattern)
				}
			}
		}
	}

	// Check allowlist if configured
	if len(t.allowPatterns) > 0 {
		allowed := false
		for _, pattern := range t.allowPatterns {
			if pattern.MatchString(lower) {
				allowed = true
				break
			}
		}
		if !allowed {
			return "Command blocked by safety guard (not in allowlist)"
		}
	}

	// Workspace restriction checks
	if t.restrictToWorkspace {
		if strings.Contains(cmd, "..\\") || strings.Contains(cmd, "../") {
			return "Command blocked by safety guard (path traversal detected)"
		}

		cwdPath, err := filepath.Abs(cwd)
		if err != nil {
			return ""
		}

		matches := absolutePathPattern.FindAllString(cmd, -1)

		for _, raw := range matches {
			p, err := filepath.Abs(raw)
			if err != nil {
				continue
			}

			if safePaths[p] {
				continue
			}

			rel, err := filepath.Rel(cwdPath, p)
			if err != nil {
				continue
			}

			if strings.HasPrefix(rel, "..") {
				return "Command blocked by safety guard (path outside working dir)"
			}
		}
	}

	return ""
}

func (t *ExecTool) getSafetyLevelName() string {
	switch t.safetyLevel {
	case SafetyStrict:
		return "strict"
	case SafetyModerate:
		return "moderate"
	case SafetyPermissive:
		return "permissive"
	case SafetyOff:
		return "off"
	default:
		return "unknown"
	}
}

func (t *ExecTool) SetTimeout(timeout time.Duration) {
	t.timeout = timeout
}

func (t *ExecTool) SetRestrictToWorkspace(restrict bool) {
	t.restrictToWorkspace = restrict
}

func (t *ExecTool) SetAllowPatterns(patterns []string) error {
	t.allowPatterns = make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return fmt.Errorf("invalid allow pattern %q: %w", p, err)
		}
		t.allowPatterns = append(t.allowPatterns, re)
	}
	return nil
}
