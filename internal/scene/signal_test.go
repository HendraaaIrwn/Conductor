package scene

import (
	"testing"
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
)

func TestSignalStateString(t *testing.T) {
	if SignalRed.String() != "R" {
		t.Errorf("red string = %q, want 'R'", SignalRed.String())
	}
	if SignalGreen.String() != "G" {
		t.Errorf("green string = %q, want 'G'", SignalGreen.String())
	}
}

func TestNewSignal(t *testing.T) {
	sig := NewSignal(10, 20, SignalRed, 5*time.Second)
	if sig.Type != entity.TypeSignal {
		t.Errorf("Type = %s, want signal", sig.Type)
	}
	if sig.X != 10 {
		t.Errorf("X = %f, want 10", sig.X)
	}
	if sig.Y != 20 {
		t.Errorf("Y = %f, want 20", sig.Y)
	}
	data, ok := sig.Data.(*SignalData)
	if !ok {
		t.Fatal("Data is not *SignalData")
	}
	if data.State != SignalRed {
		t.Error("initial state should be red")
	}
	if data.ChangeAfter != 5*time.Second {
		t.Errorf("ChangeAfter = %v, want 5s", data.ChangeAfter)
	}
}

func TestSignalBehaviorTogglesState(t *testing.T) {
	sig := NewSignal(10, 20, SignalRed, 100*time.Millisecond)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	// Before timeout: state should stay red.
	sig.Behavior.Update(sig, ctx, 50*time.Millisecond)
	data := sig.Data.(*SignalData)
	if data.State != SignalRed {
		t.Error("state should still be red after 50ms")
	}
	// After timeout: state should toggle to green.
	sig.Behavior.Update(sig, ctx, 60*time.Millisecond)
	if data.State != SignalGreen {
		t.Error("state should be green after 110ms total")
	}
	// After another timeout: back to red.
	sig.Behavior.Update(sig, ctx, 100*time.Millisecond)
	if data.State != SignalRed {
		t.Error("state should be red after another 100ms")
	}
}

func TestSignalBehaviorNoToggleWithZeroChangeAfter(t *testing.T) {
	sig := NewSignal(10, 20, SignalRed, 0)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	sig.Behavior.Update(sig, ctx, 10*time.Second)
	data := sig.Data.(*SignalData)
	if data.State != SignalRed {
		t.Error("state should not change with ChangeAfter=0")
	}
}

func TestSignalRenderFuncRedState(t *testing.T) {
	sig := NewSignal(10, 20, SignalRed, 0)
	canvas := render.NewCanvas(30, 30)
	sig.RenderFunc(canvas, sig)
	// Should have [R] at the signal head position.
	headY := int(sig.Y)
	if canvas.CellAt(9, headY).Rune != '[' {
		t.Errorf("expected '[' at (9,%d), got %q", headY, canvas.CellAt(9, headY).Rune)
	}
	if canvas.CellAt(10, headY).Rune != 'R' {
		t.Errorf("expected 'R' at (10,%d), got %q", headY, canvas.CellAt(10, headY).Rune)
	}
	if canvas.CellAt(11, headY).Rune != ']' {
		t.Errorf("expected ']' at (11,%d), got %q", headY, canvas.CellAt(11, headY).Rune)
	}
}

func TestSignalRenderFuncGreenState(t *testing.T) {
	sig := NewSignal(10, 20, SignalGreen, 0)
	canvas := render.NewCanvas(30, 30)
	sig.RenderFunc(canvas, sig)
	// Should have (G) at the signal head position.
	headY := int(sig.Y)
	if canvas.CellAt(9, headY).Rune != '(' {
		t.Errorf("expected '(' at (9,%d), got %q", headY, canvas.CellAt(9, headY).Rune)
	}
	if canvas.CellAt(10, headY).Rune != 'G' {
		t.Errorf("expected 'G' at (10,%d), got %q", headY, canvas.CellAt(10, headY).Rune)
	}
	if canvas.CellAt(11, headY).Rune != ')' {
		t.Errorf("expected ')' at (11,%d), got %q", headY, canvas.CellAt(11, headY).Rune)
	}
}

func TestSignalRenderFuncDrawsPost(t *testing.T) {
	sig := NewSignal(10, 20, SignalRed, 0)
	canvas := render.NewCanvas(30, 30)
	sig.RenderFunc(canvas, sig)
	// Post should be drawn below the head.
	if canvas.CellAt(10, 22).Rune != '|' {
		t.Errorf("expected '|' post at (10,22), got %q", canvas.CellAt(10, 22).Rune)
	}
	if canvas.CellAt(10, 23).Rune != '|' {
		t.Errorf("expected '|' post at (10,23), got %q", canvas.CellAt(10, 23).Rune)
	}
}

func TestSignalRenderFuncDifferentShapesForAccessibility(t *testing.T) {
	// Red uses [R] (square brackets), green uses (G) (parentheses).
	// This ensures signals are distinguishable in no-color mode.
	redSig := NewSignal(10, 20, SignalRed, 0)
	greenSig := NewSignal(10, 20, SignalGreen, 0)
	canvas := render.NewCanvas(30, 30)
	redSig.RenderFunc(canvas, redSig)
	redBracket := canvas.CellAt(9, 20).Rune
	greenSig.RenderFunc(canvas, greenSig)
	greenBracket := canvas.CellAt(9, 20).Rune
	if redBracket == greenBracket {
		t.Error("red and green signals should use different bracket shapes")
	}
}
