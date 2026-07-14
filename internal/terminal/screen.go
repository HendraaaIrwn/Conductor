// Package terminal wraps tcell to provide a thin abstraction over the
// underlying terminal screen. It is intentionally small so that the rest of
// the application can depend on a concrete type rather than on tcell directly,
// which keeps the simulation and rendering code testable.
package terminal

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

// MinWidth and MinHeight are the smallest terminal dimensions that Conductor
// supports. When the terminal is smaller than this the application shows a
// readable message instead of attempting to render.
const (
	MinWidth  = 80
	MinHeight = 24
)

// Screen is a thin wrapper around tcell.Screen that exposes only the
// operations Conductor needs. It also guarantees that the terminal is restored
// to its original state when Close is called.
type Screen struct {
	ts tcell.Screen
}

// New creates and initialises a Screen. The terminal is put into alternate
// screen mode, the cursor is hidden, and raw input is enabled. Close must be
// called to restore the terminal.
func New() (*Screen, error) {
	ts, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("create tcell screen: %w", err)
	}
	if err := ts.Init(); err != nil {
		return nil, fmt.Errorf("init tcell screen: %w", err)
	}
	ts.Clear()
	hideCursor(ts)
	return &Screen{ts: ts}, nil
}

// Close restores the terminal to its original state. It is safe to call
// multiple times; subsequent calls are no-ops.
func (s *Screen) Close() {
	if s == nil || s.ts == nil {
		return
	}
	showCursor(s.ts)
	s.ts.Fini()
	s.ts = nil
}

// Width returns the current terminal width in cells.
func (s *Screen) Width() int {
	w, _ := s.ts.Size()
	return w
}

// Height returns the current terminal height in cells.
func (s *Screen) Height() int {
	_, h := s.ts.Size()
	return h
}

// SetCell writes a single cell at the given coordinate with the supplied
// style. Coordinates outside the visible area are ignored.
func (s *Screen) SetCell(x, y int, r rune, style tcell.Style) {
	w, h := s.ts.Size()
	if x < 0 || x >= w || y < 0 || y >= h {
		return
	}
	s.ts.SetContent(x, y, r, nil, style)
}

// Show flushes pending cell changes to the physical terminal.
func (s *Screen) Show() {
	s.ts.Show()
}

// Sync forces a full redraw of the terminal. This is used after a resize event
// to discard any stale state the terminal might be holding.
func (s *Screen) Sync() {
	s.ts.Sync()
}

// PollEvent returns the next pending terminal event, or nil if no event is
// available. Polling is non-blocking.
func (s *Screen) PollEvent() tcell.Event {
	return s.ts.PollEvent()
}

// PostEvent tries to enqueue a custom event. It is used to inject quit signals
// from OS signal handlers.
func (s *Screen) PostEvent(ev tcell.Event) error {
	return s.ts.PostEvent(ev)
}

// IsTooSmall reports whether the terminal is smaller than the supported
// minimum.
func (s *Screen) IsTooSmall() bool {
	return s.Width() < MinWidth || s.Height() < MinHeight
}

// DrawTooSmallMessage writes the "needs more track" message onto the raw
// terminal. It is used when the terminal is below the minimum supported size.
// The caller is responsible for calling Show afterwards.
func (s *Screen) DrawTooSmallMessage() {
	s.ts.Clear()
	style := tcell.StyleDefault
	lines := []string{
		"Conductor needs a little more track.",
		"",
		fmt.Sprintf("Minimum terminal size: %dx%d", MinWidth, MinHeight),
		fmt.Sprintf("Current terminal size: %dx%d", s.Width(), s.Height()),
	}
	y := s.Height()/2 - len(lines)/2
	if y < 0 {
		y = 0
	}
	for _, line := range lines {
		x := s.Width()/2 - len(line)/2
		if x < 0 {
			x = 0
		}
		for _, r := range line {
			s.SetCell(x, y, r, style)
			x++
		}
		y++
	}
	s.ts.Show()
}

// hideCursor disables the cursor via the terminal escape sequence. tcell does
// not expose a direct API for this on every backend, so we write the DECSCUSR
// sequence to make the cursor invisible.
func hideCursor(ts tcell.Screen) {
	fmt.Fprint(os.Stderr, "\x1b[?25l")
}

// showCursor re-enables the cursor.
func showCursor(ts tcell.Screen) {
	fmt.Fprint(os.Stderr, "\x1b[?25h")
}
