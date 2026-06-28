package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/Gazi2050/cleanup/internal/ui"
)

// version is overridden via -ldflags="-X main.version=$TAG" in CI.
// Defaults to "dev" for local builds.
var version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "--version", "-v", "version":
			fmt.Printf("cleanup %s\n", version)
			return
		}
	}

	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}
