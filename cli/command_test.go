package cli

import (
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestCommand_Execute_Help(t *testing.T) {
	cmd := NewCommand("test", "Test command")

	helpFlags := []string{"--help", "-h", "help"}
	for _, flag := range helpFlags {
		exitCode := cmd.Execute([]string{flag})
		assert.Equal(t, 0, exitCode)
	}
}

func TestCommand_Execute_Version(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.SetVersion("v1.2.3")

	var exitCode int
	output := captureStdout(func() {
		exitCode = cmd.Execute([]string{"--version"})
	})

	assert.Equal(t, 0, exitCode)
	assert.Contains(t, output, "v1.2.3")
}

func TestCommand_Execute_VersionNotSet(t *testing.T) {
	cmd := NewCommand("test", "Test command")

	var exitCode int
	output := captureStdout(func() {
		exitCode = cmd.Execute([]string{"--version"})
	})

	assert.Equal(t, 0, exitCode)
	assert.Contains(t, output, "version not set")
}

func TestCommand_Execute_Subcommand(t *testing.T) {
	var receivedArgs []string

	handler := func(args []string) int {
		receivedArgs = args
		return 0
	}

	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("sub", "Subcommand description", handler)

	exitCode := cmd.Execute([]string{"sub", "--flag", "value"})

	assert.Equal(t, 0, exitCode)
	assert.Equal(t, []string{"--flag", "value"}, receivedArgs)
}

func TestCommand_Execute_SubcommandReturnsExitCode(t *testing.T) {
	handler := func(args []string) int {
		return 42
	}

	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("sub", "Subcommand", handler)

	exitCode := cmd.Execute([]string{"sub"})
	assert.Equal(t, 42, exitCode)
}

func TestCommand_Execute_UnknownSubcommand(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("valid", "Valid subcommand", func(args []string) int { return 0 })

	exitCode := cmd.Execute([]string{"unknown"})
	assert.Equal(t, 1, exitCode)
}

func TestCommand_Execute_NoArgs(t *testing.T) {
	cmd := NewCommand("test", "Test command")

	exitCode := cmd.Execute([]string{})
	assert.Equal(t, 1, exitCode)
}

func TestCommand_ShowUsage(t *testing.T) {
	cmd := NewCommand("myapp", "My awesome application")
	cmd.AddSubcommand("start", "Start the server", func(args []string) int { return 0 })
	cmd.AddSubcommand("stop", "Stop the server", func(args []string) int { return 0 })

	output := captureStdout(func() {
		cmd.ShowUsage()
	})

	assert.Contains(t, output, "myapp - My awesome application")
	assert.Contains(t, output, "Usage:")
	assert.Contains(t, output, "myapp <subcommand> [options]")
	assert.Contains(t, output, "Subcommands:")
	assert.Contains(t, output, "start")
	assert.Contains(t, output, "Start the server")
	assert.Contains(t, output, "stop")
	assert.Contains(t, output, "Stop the server")
	assert.Contains(t, output, "Global Flags:")
	assert.Contains(t, output, "--version")
	assert.Contains(t, output, "--help")
}

func TestCommand_SubcommandOrder(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("zebra", "Last alphabetically", func(args []string) int { return 0 })
	cmd.AddSubcommand("alpha", "First alphabetically", func(args []string) int { return 0 })
	cmd.AddSubcommand("middle", "Middle alphabetically", func(args []string) int { return 0 })

	output := captureStdout(func() {
		cmd.ShowUsage()
	})

	// Check that subcommands are sorted alphabetically in output.
	alphaIdx := strings.Index(output, "alpha")
	middleIdx := strings.Index(output, "middle")
	zebraIdx := strings.Index(output, "zebra")

	assert.True(t, alphaIdx < middleIdx, "alpha should come before middle")
	assert.True(t, middleIdx < zebraIdx, "middle should come before zebra")
}

func TestCommand_MultipleSubcommands(t *testing.T) {
	var executed []string

	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("first", "First command", func(args []string) int {
		executed = append(executed, "first")
		return 0
	})
	cmd.AddSubcommand("second", "Second command", func(args []string) int {
		executed = append(executed, "second")
		return 1
	})
	cmd.AddSubcommand("third", "Third command", func(args []string) int {
		executed = append(executed, "third")
		return 2
	})

	// Execute first
	exitCode := cmd.Execute([]string{"first"})
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, []string{"first"}, executed)

	// Execute second
	exitCode = cmd.Execute([]string{"second"})
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, []string{"first", "second"}, executed)

	// Execute third
	exitCode = cmd.Execute([]string{"third"})
	assert.Equal(t, 2, exitCode)
	assert.Equal(t, []string{"first", "second", "third"}, executed)
}

func TestCommand_SubcommandWithComplexArgs(t *testing.T) {
	var receivedArgs []string

	handler := func(args []string) int {
		receivedArgs = args
		return 0
	}

	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("build", "Build command", handler)

	cmd.Execute([]string{"build", "--verbose", "-o", "output.txt", "file1.go", "file2.go"})

	expected := []string{"--verbose", "-o", "output.txt", "file1.go", "file2.go"}
	assert.Equal(t, expected, receivedArgs)
}

func TestCommand_SubcommandWithFlagSet(t *testing.T) {
	// This tests that subcommand handlers can use FlagSet
	type Config struct {
		Input  string `flag:"input,i" usage:"input file" required:"true"`
		Output string `flag:"output,o" usage:"output file" default:"out.txt"`
		Count  int    `flag:"count,n" usage:"repeat count" default:"1"`
	}

	var config *Config

	handler := func(args []string) int {
		fs := NewFlagSet("test sub")
		config = &Config{}
		fs.AddSection("Options", config)

		_, err := fs.Parse(args)
		if err != nil {
			return 1
		}

		return 0
	}

	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("sub", "Subcommand with flags", handler)

	exitCode := cmd.Execute([]string{"sub", "-i", "in.txt", "-o", "custom.txt", "-n", "5"})

	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "in.txt", config.Input)
	assert.Equal(t, "custom.txt", config.Output)
	assert.Equal(t, 5, config.Count)
}

func TestCommand_NoSubcommands(t *testing.T) {
	cmd := NewCommand("test", "Test command with no subcommands")

	output := captureStdout(func() {
		cmd.ShowUsage()
	})

	// Should still show usage even with no subcommands
	assert.Contains(t, output, "test - Test command with no subcommands")
	assert.Contains(t, output, "Global Flags:")
}

func TestCommand_VersionAlternativeFlag(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.SetVersion("v2.0.0")

	// Test -version (single dash) as well as --version
	output := captureStdout(func() {
		exitCode := cmd.Execute([]string{"-version"})
		assert.Equal(t, 0, exitCode)
	})

	assert.Contains(t, output, "v2.0.0")
}

func TestCommand_EmptySubcommandName(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("", "Empty name", func(args []string) int { return 0 })
	cmd.AddSubcommand("valid", "Valid subcommand", func(args []string) int { return 0 })

	// Empty string subcommand name should still work (edge case)
	exitCode := cmd.Execute([]string{""})
	assert.Equal(t, 0, exitCode)
}

func TestCommand_DuplicateSubcommandRegistration(t *testing.T) {
	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("sub", "First registration", func(args []string) int { return 1 })
	cmd.AddSubcommand("sub", "Second registration", func(args []string) int { return 2 })

	// Last registration should win
	exitCode := cmd.Execute([]string{"sub"})
	assert.Equal(t, 2, exitCode)
}

func TestCommand_SubcommandWithOnlyFlags(t *testing.T) {
	var receivedArgs []string

	handler := func(args []string) int {
		receivedArgs = args
		return 0
	}

	cmd := NewCommand("test", "Test command")
	cmd.AddSubcommand("sub", "Subcommand", handler)

	exitCode := cmd.Execute([]string{"sub"})

	assert.Equal(t, 0, exitCode)
	assert.Empty(t, receivedArgs)
}
