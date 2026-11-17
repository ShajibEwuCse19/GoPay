package main

import (
	"fmt"
	"os"

	"github.com/ShajibEwuCse19/gopay-core/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		cmd.Serve()
	case "migrate":
		cmd.Migrate(os.Args[2:]) // Passes sub-arguments (e.g. "up", "down") to cmd.Migrate
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "GoPay Core CLI")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  go run . <command> [arguments]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Available Commands:")
	fmt.Fprintln(os.Stderr, "  serve                       Start the HTTP REST API server")
	fmt.Fprintln(os.Stderr, "  migrate up [path]           Apply all pending database migrations (default: config or ./migrations)")
	fmt.Fprintln(os.Stderr, "  migrate down [path]         Roll back the most recent migration")
}
