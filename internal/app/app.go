// Package app wires together the terminal, engine loop, and lifecycle
// management. It is the top-level orchestrator invoked by main.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/engine"
	"github.com/example/conductor/internal/terminal"
)

// App holds the runtime configuration and owns the terminal screen.
type App struct {
	cfg config.Config
}

// New creates an App from the given configuration.
func New(cfg config.Config) *App {
	return &App{cfg: cfg}
}

// Run initialises the terminal, starts the animation loop, and blocks until
// the user quits. It guarantees that the terminal is restored on exit,
// including when a panic occurs or an OS signal is received.
func (a *App) Run() (err error) {
	screen, err := terminal.New()
	if err != nil {
		return fmt.Errorf("terminal init: %w", err)
	}

	// Ensure the terminal is restored no matter how we exit. This catches
	// normal returns, panics in the loop, and signal-driven cancellation.
	defer func() {
		if r := recover(); r != nil {
			screen.Close()
			err = fmt.Errorf("panic: %v", r)
			return
		}
		screen.Close()
	}()

	loop, err := engine.NewLoop(screen, a.cfg)
	if err != nil {
		return fmt.Errorf("create loop: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS signals and cancel the context so the loop exits
	// gracefully. SIGINT covers Ctrl+C at the process level; SIGTERM covers
	// process termination.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		<-sigCh
		cancel()
	}()

	if err := loop.Run(ctx); err != nil {
		return fmt.Errorf("animation loop: %w", err)
	}
	return nil
}
