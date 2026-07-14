package effects

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// rainStyle is the base style for rain particles.
var rainStyle = tcell.StyleDefault.Foreground(tcell.ColorBlue)

// rainSprite is a shared single-cell sprite for rain.
var rainSprite = render.NewSprite("rain",
	render.CellGrid{{render.Cell{Rune: '|', Style: rainStyle}}})

// RainBehavior drives a rain particle: it falls diagonally and recycles to
// the top when it exits the bottom of the viewport. Each particle has a
// fixed horizontal offset for its recycle position, determined at spawn time
// via the entity's Data field.
type RainBehavior struct{}

// rainData holds the recycle X position for a rain particle.
type rainData struct {
	RecycleX float64
}

// Update moves the rain particle down and slightly sideways, recycling when
// it exits the viewport.
func (RainBehavior) Update(e *entity.Entity, ctx entity.UpdateContext, delta time.Duration) {
	sec := delta.Seconds()
	e.X += e.VX * sec
	e.Y += e.VY * sec

	// Recycle when the particle exits the bottom.
	if e.Y > float64(ctx.Height) {
		e.Y = -1
		if data, ok := e.Data.(*rainData); ok && data != nil {
			e.X = data.RecycleX
		}
	}
	// Wrap horizontally.
	if e.X < 0 {
		e.X = float64(ctx.Width)
	}
	if e.X > float64(ctx.Width) {
		e.X = 0
	}
}

// NewRain creates a rain particle entity at the given position. The particle
// falls diagonally with a downward and slight horizontal velocity. When it
// exits the bottom, it recycles to recycleX at the top.
func NewRain(x, y, recycleX float64) *entity.Entity {
	return &entity.Entity{
		Type:            entity.TypeRain,
		X:               x,
		Y:               y,
		VX:              3.0,  // slight diagonal
		VY:              30.0, // fast fall
		Layer:           render.LayerParticles,
		Sprite:          rainSprite,
		Visible:         true,
		Behavior:        RainBehavior{},
		Data:            &rainData{RecycleX: recycleX},
		RemoveOffscreen: false, // recycled, not removed
	}
}

// SpawnRain creates a set of rain particles distributed across the viewport.
// The count scales with viewport area.
func SpawnRain(manager *entity.Manager, width, height int, count int) {
	for i := 0; i < count; i++ {
		x := float64(i * width / count)
		y := float64(i % height)
		recycleX := x
		manager.Add(NewRain(x, y, recycleX))
	}
}

// RainParticleCount returns the recommended number of rain particles for the
// given viewport dimensions. Density scales with width but is capped.
func RainParticleCount(width, height int) int {
	count := width * height / 100
	if count < 20 {
		count = 20
	}
	if count > 150 {
		count = 150
	}
	return count
}
