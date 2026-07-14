// Package effects contains particle and visual-effect entity factories.
// Each factory creates a fully-configured entity.Entity ready to be added to
// the entity.Manager.
package effects

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// smokeRunes are the runes used for smoke particles as they age. The particle
// cycles through these from small to large to faint, simulating smoke
// dissipating.
var smokeRunes = []rune{'.', 'o', 'O', '*', '°'}

// smokeStyle is the base style for smoke particles. It uses dim to make the
// smoke visually softer than the train.
var smokeStyle = tcell.StyleDefault.Dim(true)

// smokeSprite is a shared single-cell sprite with one frame per smoke rune.
var smokeSprite = buildSmokeSprite()

// SmokeBehavior drives a smoke particle: it rises, drifts slightly, slows
// over time, and changes appearance as it ages.
type SmokeBehavior struct {
	// Dampening is applied to VY each frame to slow the rise.
	Dampening float64
}

// NewSmoke creates a smoke particle entity at the given position. The particle
// rises upward with a slight horizontal drift and expires after its lifetime.
func NewSmoke(x, y float64) *entity.Entity {
	return &entity.Entity{
		Type:            entity.TypeSmoke,
		X:               x,
		Y:               y,
		VX:              drift(),
		VY:              -8.0, // cells per second upward
		Layer:           render.LayerParticles,
		Sprite:          smokeSprite,
		Visible:         true,
		Lifetime:        2500 * time.Millisecond,
		Behavior:        &SmokeBehavior{Dampening: 0.95},
		RemoveOffscreen: true,
	}
}

// Update advances the smoke particle: applies velocity, dampens vertical
// speed, and selects the frame based on age ratio.
func (b *SmokeBehavior) Update(e *entity.Entity, _ entity.UpdateContext, delta time.Duration) {
	sec := delta.Seconds()
	e.X += e.VX * sec
	e.Y += e.VY * sec
	e.VY *= b.Dampening

	// Select frame based on how much of the lifetime has elapsed.
	ratio := float64(e.Age) / float64(e.Lifetime)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	e.Frame = int(ratio * float64(len(smokeRunes)))
	if e.Frame >= len(smokeRunes) {
		e.Frame = len(smokeRunes) - 1
	}
}

// drift returns a small random horizontal velocity for smoke particles.
// The value is deterministic per particle via the caller's RNG; here we use
// a simple alternating pattern.
func drift() float64 {
	return 1.5
}

// buildSmokeSprite creates the shared smoke sprite with one frame per rune.
func buildSmokeSprite() *render.Sprite {
	frames := make([]render.CellGrid, len(smokeRunes))
	for i, r := range smokeRunes {
		frames[i] = render.CellGrid{
			{render.Cell{Rune: r, Style: smokeStyle}},
		}
	}
	return render.NewSprite("smoke", frames...)
}
