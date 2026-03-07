package tools

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecTool_GrepExitCode(t *testing.T) {
	// Skip on Windows as grep is not available by default
	if runtime.GOOS == "windows" {
		t.Skip("Skipping grep tests on Windows")
	}

	// Check if grep is available
	if _, err := exec.LookPath("grep"); err != nil {
		t.Skip("grep command not found, skipping test")
	}
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Create a test file with some content
	err := os.WriteFile(testFile, []byte("hello world\nfoo bar\nbaz qux\n"), 0644)
	require.NoError(t, err)

	tool, err := NewExecTool(tmpDir, false)
	require.NoError(t, err)

	tests := []struct {
		name        string
		command     string
		expectError bool
		expectMsg   string
	}{
		{
			name:        "grep finds matches - exit code 0",
			command:     "grep hello " + testFile,
			expectError: false,
			expectMsg:   "hello world",
		},
		{
			name:        "grep no matches - exit code 1 (should NOT be error)",
			command:     "grep nonexistent " + testFile,
			expectError: false,
			expectMsg:   "(no matches found)",
		},
		{
			name:        "grep with -r flag no matches",
			command:     "grep -r \"subagent-16\" " + tmpDir,
			expectError: false,
			expectMsg:   "(no matches found)",
		},
		{
			name:        "non-grep command with exit code 1 (should be error)",
			command:     "sh -c 'exit 1'",
			expectError: true,
			expectMsg:   "Exit code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tool.Execute(context.Background(), map[string]any{
				"command": tt.command,
			})

			assert.Equal(t, tt.expectError, result.IsError,
				"Expected IsError=%v but got %v. Output: %s",
				tt.expectError, result.IsError, result.ForLLM)

			if tt.expectMsg != "" {
				assert.Contains(t, result.ForLLM, tt.expectMsg,
					"Expected output to contain '%s' but got: %s",
					tt.expectMsg, result.ForLLM)
			}
		})
	}
}

func TestExecTool_GrepVariations(t *testing.T) {
	// Skip on Windows as grep is not available by default
	if runtime.GOOS == "windows" {
		t.Skip("Skipping grep tests on Windows")
	}

	// Check if grep is available
	if _, err := exec.LookPath("grep"); err != nil {
		t.Skip("grep command not found, skipping test")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test content\n"), 0644)
	require.NoError(t, err)

	tool, err := NewExecTool(tmpDir, false)
	require.NoError(t, err)

	grepCommands := []string{
		"grep nonexistent " + testFile,
		"grep -r nonexistent " + tmpDir,
		"grep -i NONEXISTENT " + testFile,
		"grep -n nonexistent " + testFile,
		"grep -v test " + testFile + " | grep nonexistent",
	}

	for _, cmd := range grepCommands {
		t.Run(cmd, func(t *testing.T) {
			result := tool.Execute(context.Background(), map[string]any{
				"command": cmd,
			})

			// All grep commands with no matches should NOT be errors
			assert.False(t, result.IsError,
				"grep command with no matches should not be error: %s\nOutput: %s",
				cmd, result.ForLLM)

			// Should contain "no matches found" message
			assert.True(t,
				strings.Contains(result.ForLLM, "(no matches found)") ||
					strings.Contains(result.ForLLM, "test content"),
				"Expected '(no matches found)' or actual output, got: %s",
				result.ForLLM)
		})
	}
}

func TestExecTool_GrepActualError(t *testing.T) {
	// Skip on Windows as grep is not available by default
	if runtime.GOOS == "windows" {
		t.Skip("Skipping grep tests on Windows")
	}

	// Check if grep is available
	if _, err := exec.LookPath("grep"); err != nil {
		t.Skip("grep command not found, skipping test")
	}

	tmpDir := t.TempDir()
	tool, err := NewExecTool(tmpDir, false)
	require.NoError(t, err)

	// grep with invalid file should be exit code 2 (actual error)
	result := tool.Execute(context.Background(), map[string]any{
		"command": "grep test /nonexistent/file/path",
	})

	// This should be an error (exit code 2)
	assert.True(t, result.IsError,
		"grep with invalid file should be error. Output: %s", result.ForLLM)
}
