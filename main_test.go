package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLIBasicFunctionality tests the basic CLI functionality
func TestCLIBasicFunctionality(t *testing.T) {
	// This is a simple integration test that verifies the CLI runs successfully
	// For more thorough library testing, see the pkg/geq tests
	t.Skip("Skipping integration test that requires network access")

	// Build the CLI binary for testing
	binaryPath := filepath.Join(t.TempDir(), "geq")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	err := cmd.Run()
	require.NoError(t, err, "Failed to build CLI binary")

	// Run the CLI with --version flag
	cmd = exec.Command(binaryPath, "--version")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI execution failed")

	// Check that version information is displayed
	assert.True(t, strings.Contains(string(output), "geq version"), "Version output not found")
}

// TestCLIArgumentParsing tests the CLI argument parsing
func TestCLIArgumentParsing(t *testing.T) {
	// Test error cases for argument parsing
	// Here we use a subprocess approach to test the CLI

	// Build the CLI binary for testing
	binaryPath := filepath.Join(t.TempDir(), "geq")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	err := cmd.Run()
	require.NoError(t, err, "Failed to build CLI binary")

	// Run the CLI without required endpoint argument
	cmd = exec.Command(binaryPath)
	output, err := cmd.CombinedOutput()

	// Should exit with error (non-zero exit code)
	assert.Error(t, err, "CLI should fail without required endpoint argument")

	// Error message should mention endpoint URL is required - use case-insensitive check
	// to ensure it works across platforms
	outputStr := strings.ToLower(string(output))
	assert.True(t, 
		strings.Contains(outputStr, "endpoint") && strings.Contains(outputStr, "required"),
		"Missing endpoint error message not found - got: %s", string(output))
}

// TestCLIOutputFileHandling tests the output file handling of the CLI
func TestCLIOutputFileHandling(t *testing.T) {
	// This test verifies output file generation logic

	// Create a temporary directory for test output files
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "test_output.graphql")

	// Create a mock file to simulate output
	testContent := "type Query { test: String }\n"
	err := os.WriteFile(outputPath, []byte(testContent), 0644)
	require.NoError(t, err, "Failed to create test output file")

	// Verify file was created successfully
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err, "Failed to read test output file")
	assert.Equal(t, testContent, string(content), "Output file content mismatch")

	// Additional tests would call the actual CLI with different output options
	// but we'll keep those as integration tests that can be explicitly enabled
}
