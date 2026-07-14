package engine

import (
	"context"
	"time"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/input"
	"github.com/example/conductor/internal/render"
	"github.com/example/conductor/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

// Loop is the main animation driver. It owns the terminal screen, canvas,
// clock, and world, and runs the frame loop until the context is cancelled or
// the user quits.
type Loop struct {
	screen *terminal.Screen
	canvas *render.Canvas
	clock  *Clock
	world  *World
	input  *input.Handler
	quit   bool
}

// NewLoop creates a Loop bound to the given screen, configured by the given
// Config. The config values are applied to the clock and world.
func NewLoop(screen *terminal.Screen, cfg config.Config) (*Loop, error) {
	width := screen.Width()
	height := screen.Height()
	canvas := render.NewCanvas(width, height)

	// Apply reduced-motion to FPS.
	fps := cfg.FPS
	if cfg.ReducedMotion && fps > 10 {
		fps = 10
	}
	clock := NewClock(fps)

	// Build the world with the config's seed and style.
	style := tcell.StyleDefault
	world := NewWorldWithConfig(width, height, cfg, style)

	handler := input.NewHandler(world, clock)
	return &Loop{
		screen: screen,
		canvas: canvas,
		clock:  clock,
		world:  world,
		input:  handler,
	}, nil
}

// Run starts the animation loop. It blocks until the user quits or the context
// is cancelled. The caller is responsible for closing the terminal screen.
func (l *Loop) Run(ctx context.Context) error {
	for {
		if l.quit {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// 1. Handle all pending terminal events (non-blocking).
		l.processEvents()
		if l.quit {
			return nil
		}

		// 2. If the terminal is too small, show the message and wait.
		if l.screen.IsTooSmall() {
			l.screen.DrawTooSmallMessage()
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// 3. Calculate delta time (zero when paused).
		delta := l.clock.Tick()

		// 4. Update simulation state.
		l.world.Update(delta)

		// 5. Render the world to the back buffer.
		l.canvas.Clear()
		l.world.Render(l.canvas)

		// 6. Flush only changed cells to the terminal.
		l.canvas.Flush(l.screen)

		// 7. Sleep to maintain the target frame rate.
		elapsed := time.Since(l.clock.lastTick)
		remaining := l.clock.FrameDuration() - elapsed
		if remaining > 0 {
			time.Sleep(remaining)
		}
	}
}

// processEvents polls all pending terminal events and dispatches them. Resize
// events update the canvas and world; key events go to the input handler.
func (l *Loop) processEvents() {
	for {
		ev := l.screen.PollEvent()
		if ev == nil {
			return
		}
		switch e := ev.(type) {
		case *tcell.EventResize:
			l.handleResize()
		case *tcell.EventKey:
			if l.input.Handle(e) {
				l.quit = true
				return
			}
		}
	}
}

// handleResize updates the canvas and world after a terminal resize.
func (l *Loop) handleResize() {
	width := l.screen.Width()
	height := l.screen.Height()
	l.canvas.Resize(width, height)
	l.world.HandleResize(width, height)
	l.screen.Sync()
}
