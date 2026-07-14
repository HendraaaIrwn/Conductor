// Package input handles keyboard events and maps them to world and clock
// actions. The handler returns true from Handle when the application should
// quit.
//
// To avoid an import cycle with the engine package, the handler depends on
// small interfaces (WorldController, ClockController) rather than on the
// concrete engine types. The engine types satisfy these interfaces via Go's
// structural typing.
package input

import "github.com/gdamore/tcell/v2"

// WorldController is the subset of world operations that keyboard input can
// trigger.
type WorldController interface {
	Paused() bool
	SetPaused(bool)
	Regenerate()
	SpawnNext()
	ChangeScene()
	ChangeWeather()
	ChangeTime()
	ChangeColorPalette()
	ToggleHelp()
	SetHelpVisible(bool)
	HelpVisible() bool
	Speed() float64
	SetSpeed(float64)
}

// ClockController is the subset of clock operations that keyboard input can
// trigger.
type ClockController interface {
	SetPaused(bool)
}

// Handler maps key presses to simulation actions.
type Handler struct {
	world WorldController
	clock ClockController
}

// NewHandler creates a Handler bound to the given world and clock controllers.
func NewHandler(world WorldController, clock ClockController) *Handler {
	return &Handler{world: world, clock: clock}
}

// Handle processes a single key event. It returns true if the application
// should quit.
func (h *Handler) Handle(ev tcell.Event) bool {
	keyEvent, ok := ev.(*tcell.EventKey)
	if !ok {
		return false
	}
	switch keyEvent.Key() {
	case tcell.KeyCtrlC:
		return true
	case tcell.KeyEscape:
		// Close help if visible, otherwise no-op.
		if h.world.HelpVisible() {
			h.world.SetHelpVisible(false)
		}
		return false
	}

	r := keyEvent.Rune()
	switch r {
	case 'q', 'Q':
		return true
	case 'p', 'P':
		h.togglePause()
	case ' ':
		h.togglePause()
	case 'r', 'R':
		h.world.Regenerate()
	case 'n', 'N':
		h.world.SpawnNext()
	case 's', 'S':
		h.world.ChangeScene()
	case 'w', 'W':
		h.world.ChangeWeather()
	case 'd', 'D':
		h.world.ChangeTime()
	case 'c', 'C':
		h.world.ChangeColorPalette()
	case 'h', 'H', '?':
		h.world.ToggleHelp()
	case '+':
		h.adjustSpeed(5.0)
	case '-':
		h.adjustSpeed(-5.0)
	case '0':
		h.resetSpeed()
	}
	return false
}

// togglePause pauses or resumes both the clock and the world.
func (h *Handler) togglePause() {
	paused := !h.world.Paused()
	h.world.SetPaused(paused)
	h.clock.SetPaused(paused)
}

// adjustSpeed increases or decreases the train speed by the given delta
// (cells per second). The speed is clamped to a reasonable range.
func (h *Handler) adjustSpeed(delta float64) {
	newSpeed := h.world.Speed() + delta
	if newSpeed < 5.0 {
		newSpeed = 5.0
	}
	if newSpeed > 80.0 {
		newSpeed = 80.0
	}
	h.world.SetSpeed(newSpeed)
}

// resetSpeed restores the default train speed.
func (h *Handler) resetSpeed() {
	h.world.SetSpeed(25.0)
}
