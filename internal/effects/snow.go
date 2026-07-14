package effects

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// snowStyle is the base style for snow particles.
var snowStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite)

// snowSprite is a shared single-cell sprite for snow.
var snowSprite = render.NewSprite("snow",
	render.CellGrid{{render.Cell{Rune: '*', Style: snowStyle}}})

// SnowBehavior drives a snow particle: it falls slowly with slight
// horizontal drift, and recycles to the top when it exits the bottom.
type SnowBehavior struct {
	// DriftPhase controls the sine-wave horizontal drift.
	DriftPhase float64
}

// snowData holds per-particle drift state.
type snowData struct {
	Phase    float64
	RecycleX float64
	BaseX    float64
}

// Update moves the snow particle down with a gentle sine-wave drift, recycling
// when it exits the viewport.
func (SnowBehavior) Update(e *entity.Entity, ctx entity.UpdateContext, delta time.Duration) {
	sec := delta.Seconds()
	data, ok := e.Data.(*snowData)
	if !ok || data == nil {
		e.Y += e.VY * sec
		return
	}

	// Fall slowly.
	e.Y += e.VY * sec

	// Sine-wave drift: X = baseX + sin(phase + age) * amplitude.
	ageSec := e.Age.Seconds()
	e.X = data.BaseX + sinApprox(data.Phase+ageSec*2.0)*2.0

	// Recycle when the particle exits the bottom.
	if e.Y > float64(ctx.Height) {
		e.Y = -1
		e.X = data.RecycleX
		data.BaseX = data.RecycleX
	}
}

// NewSnow creates a snow particle entity at the given position. The particle
// falls slowly with a gentle horizontal drift and recycles at the bottom.
func NewSnow(x, y, recycleX float64, phase float64) *entity.Entity {
	return &entity.Entity{
		Type:            entity.TypeSnow,
		X:               x,
		Y:               y,
		VY:              6.0, // slow fall
		Layer:           render.LayerParticles,
		Sprite:          snowSprite,
		Visible:         true,
		Behavior:        SnowBehavior{},
		Data:            &snowData{Phase: phase, RecycleX: recycleX, BaseX: x},
		RemoveOffscreen: false, // recycled, not removed
	}
}

// SpawnSnow creates a set of snow particles distributed across the viewport.
func SpawnSnow(manager *entity.Manager, width, height int, count int) {
	for i := 0; i < count; i++ {
		x := float64(i * width / count)
		y := float64(i % height)
		phase := float64(i) * 0.7
		manager.Add(NewSnow(x, y, x, phase))
	}
}

// SnowParticleCount returns the recommended number of snow particles. Snow
// uses fewer particles than rain for visual variety and performance.
func SnowParticleCount(width, height int) int {
	count := width * height / 150
	if count < 10 {
		count = 10
	}
	if count > 80 {
		count = 80
	}
	return count
}

// sinApprox is a simple sine approximation using a Taylor-series-like
// polynomial. It avoids importing math.Sin in the hot path and is accurate
// enough for visual drift (which doesn't need precision).
func sinApprox(x float64) float64 {
	// Normalize to [0, 2π).
	const pi2 = 6.283185307179586
	for x < 0 {
		x += pi2
	}
	for x >= pi2 {
		x -= pi2
	}
	// Shift to [-π, π].
	if x > 3.141592653589793 {
		x -= pi2
	}
	// Taylor series: x - x³/6 + x⁵/120.
	x2 := x * x
	return x * (1 - x2*(1-x2*(1-x2/42)/20)/6)
}
