package effects

import (
	"testing"
	"time"

	"github.com/example/conductor/internal/entity"
)

func TestNewSmoke(t *testing.T) {
	e := NewSmoke(10, 20)
	if e.Type != entity.TypeSmoke {
		t.Errorf("Type = %s, want smoke", e.Type)
	}
	if e.X != 10 {
		t.Errorf("X = %f, want 10", e.X)
	}
	if e.Y != 20 {
		t.Errorf("Y = %f, want 20", e.Y)
	}
	if e.VY >= 0 {
		t.Errorf("VY = %f, want negative (rising)", e.VY)
	}
	if e.Lifetime <= 0 {
		t.Error("smoke should have a finite lifetime")
	}
	if !e.RemoveOffscreen {
		t.Error("smoke should be removed when offscreen")
	}
	if e.Sprite == nil {
		t.Error("smoke should have a sprite")
	}
}

func TestSmokeRisesOverTime(t *testing.T) {
	e := NewSmoke(10, 20)
	initialY := e.Y
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	e.Behavior.Update(e, ctx, 100*time.Millisecond)
	if e.Y >= initialY {
		t.Errorf("smoke should rise: Y = %f, initial = %f", e.Y, initialY)
	}
}

func TestSmokeChangesFrameWithAge(t *testing.T) {
	e := NewSmoke(10, 20)
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	initialFrame := e.Frame
	// Simulate aging: the manager normally increments Age before the next
	// Update call. We set Age directly to test the frame selection logic.
	e.Age = e.Lifetime / 2
	e.Behavior.Update(e, ctx, 1*time.Millisecond)
	if e.Frame == initialFrame {
		t.Error("smoke frame should change as it ages")
	}
}

func TestSmokeDampensVerticalVelocity(t *testing.T) {
	e := NewSmoke(10, 20)
	initialVY := e.VY
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	e.Behavior.Update(e, ctx, 100*time.Millisecond)
	// Dampening makes VY closer to zero (less negative). Since both values
	// are negative, the dampened value should be greater than the initial.
	if e.VY <= initialVY {
		t.Errorf("VY should dampen toward zero: %f -> %f", initialVY, e.VY)
	}
}

func TestSmokeSpriteValidates(t *testing.T) {
	if err := smokeSprite.Validate(); err != nil {
		t.Errorf("smoke sprite failed validation: %v", err)
	}
}

func TestSmokeSpriteHasMultipleFrames(t *testing.T) {
	if smokeSprite.FrameCount() < 3 {
		t.Errorf("smoke sprite frames = %d, want >= 3", smokeSprite.FrameCount())
	}
}
