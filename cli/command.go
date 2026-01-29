// Package cli provides utilities for building command-line interface applications.
// It supports subcommand routing, flag parsing with struct tags, positional arguments,
// environment variable integration, and organized usage output.
package cli

import (
	"fmt"
	"os"
	"sort"
)

// Command represents a command with subcommands.
type Command struct {
	name        string
	description string
	subcommands map[string]*Subcommand
	version     string
}

// Subcommand represents a registered subcommand.
type Subcommand struct {
	name        string
	description string
	handler     SubcommandHandler
}

// SubcommandHandler is called when a subcommand is executed.
// It receives the remaining arguments after the subcommand name
// and returns an exit code.
type SubcommandHandler func(args []string) int

// NewCommand creates a new Command with subcommands.
func NewCommand(name, description string) *Command {
	return &Command{
		name:        name,
		description: description,
		subcommands: make(map[string]*Subcommand),
	}
}

// AddSubcommand registers a subcommand with its handler.
func (cmd *Command) AddSubcommand(name, description string, handler SubcommandHandler) {
	cmd.subcommands[name] = &Subcommand{
		name:        name,
		description: description,
		handler:     handler,
	}
}

// SetVersion sets the version string for --version flag.
func (cmd *Command) SetVersion(version string) {
	cmd.version = version
}

// Execute routes to the appropriate subcommand and returns an exit code.
func (cmd *Command) Execute(args []string) int {
	if len(args) == 0 {
		cmd.ShowUsage()
		return 1
	}

	// Handle global flags.
	switch args[0] {
	case "--version", "-version":
		if cmd.version != "" {
			fmt.Println(cmd.version)
		} else {
			fmt.Println("version not set")
		}
		return 0
	case "--help", "-h", "help":
		cmd.ShowUsage()
		return 0
	}

	// Route to subcommand.
	cmdName := args[0]
	subcmd, exists := cmd.subcommands[cmdName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n\n", cmdName)
		cmd.ShowUsage()
		return 1
	}

	// Execute subcommand handler with remaining args.
	return subcmd.handler(args[1:])
}

// ShowUsage prints the main usage message with all subcommands.
func (cmd *Command) ShowUsage() {
	fmt.Printf("%s - %s\n\n", cmd.name, cmd.description)
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s <subcommand> [options]\n\n", cmd.name)

	if len(cmd.subcommands) > 0 {
		fmt.Println("Subcommands:")

		// Sort subcommands for consistent output.
		names := make([]string, 0, len(cmd.subcommands))
		for name := range cmd.subcommands {
			names = append(names, name)
		}
		sort.Strings(names)

		// Find max width for alignment.
		maxWidth := 0
		for _, name := range names {
			maxWidth = max(maxWidth, len(name))
		}

		// Print subcommands.
		for _, name := range names {
			subcmd := cmd.subcommands[name]
			fmt.Printf("  %-*s  %s\n", maxWidth, name, subcmd.description)
		}
		fmt.Println()
	}

	fmt.Println("Global Flags:")
	fmt.Println("  --version     Show version information")
	fmt.Println("  --help        Show this help message")
	fmt.Println()
	fmt.Printf("Run '%s <subcommand> --help' for subcommand-specific options.\n", cmd.name)
}
