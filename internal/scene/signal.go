package scene

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// SignalState represents the state of a railway signal. Signals use both
// color and shape so they are understandable in no-color mode.
type SignalState int

const (
	SignalRed SignalState = iota
	SignalGreen
)

// String returns the text label for the signal state, used in no-color mode.
func (s SignalState) String() string {
	switch s {
	case SignalRed:
		return "R"
	case SignalGreen:
		return "G"
	default:
		return "?"
	}
}

// SignalData holds the state for a signal entity. It is stored in the
// entity's Data field and updated by SignalBehavior.
type SignalData struct {
	State       SignalState
	Elapsed     time.Duration
	ChangeAfter time.Duration // duration after which to toggle; 0 = manual
}

// SignalBehavior updates a signal entity: it toggles the state after a
// configured duration if ChangeAfter is non-zero.
type SignalBehavior struct{}

// Update advances the signal timer and toggles state when the duration
// elapses.
func (SignalBehavior) Update(e *entity.Entity, _ entity.UpdateContext, delta time.Duration) {
	data, ok := e.Data.(*SignalData)
	if !ok || data == nil || data.ChangeAfter == 0 {
		return
	}
	data.Elapsed += delta
	if data.Elapsed >= data.ChangeAfter {
		data.Elapsed = 0
		if data.State == SignalRed {
			data.State = SignalGreen
		} else {
			data.State = SignalRed
		}
	}
}

// SignalRenderFunc draws a signal as a post with a colored/lettered head.
// The shape differs by state so it is readable in monochrome:
//
//	[R]   red   — square brackets with 'R'
//	(G)   green — parentheses with 'G'
func SignalRenderFunc(canvas *render.Canvas, e *entity.Entity) {
	data, ok := e.Data.(*SignalData)
	if !ok || data == nil {
		return
	}
	x, y := int(e.X), int(e.Y)

	// Draw the signal post.
	postStyle := tcell.StyleDefault
	canvas.SetRune(x, y+2, '|', postStyle)
	canvas.SetRune(x, y+3, '|', postStyle)

	// Draw the signal head with state-specific shape and color.
	var headStyle tcell.Style
	var openBracket, closeBracket rune
	var letter rune

	if data.State == SignalRed {
		headStyle = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
		openBracket = '['
		closeBracket = ']'
		letter = 'R'
	} else {
		headStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
		openBracket = '('
		closeBracket = ')'
		letter = 'G'
	}

	canvas.SetRune(x-1, y, openBracket, headStyle)
	canvas.SetRune(x, y, letter, headStyle)
	canvas.SetRune(x+1, y, closeBracket, headStyle)
}

// NewSignal creates a signal entity at the given position. If changeAfter is
// non-zero, the signal automatically toggles state after that duration.
func NewSignal(x, y int, initialState SignalState, changeAfter time.Duration) *entity.Entity {
	return &entity.Entity{
		Type:       entity.TypeSignal,
		X:          float64(x),
		Y:          float64(y),
		Layer:      render.LayerPlatform,
		Visible:    true,
		Behavior:   SignalBehavior{},
		RenderFunc: SignalRenderFunc,
		Data: &SignalData{
			State:       initialState,
			ChangeAfter: changeAfter,
		},
	}
}
