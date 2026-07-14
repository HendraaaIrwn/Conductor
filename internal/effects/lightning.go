package effects

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// LightningData holds state for a lightning flash effect.
type LightningData struct {
	// Duration is how long the flash lasts.
	Duration time.Duration
	// BoltRune is the character used to draw the bolt.
	BoltRune rune
}

// lightningSprite is a shared sprite for the lightning bolt.
var lightningSprite = render.NewSprite("lightning",
	render.CellGrid{
		{render.Cell{Rune: '/', Style: tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)}},
	})

// NewLightning creates a lightning flash entity at the given column. The
// flash is a brief vertical bolt that expires after a short lifetime.
func NewLightning(x float64, height int) *entity.Entity {
	boltHeight := height / 3
	if boltHeight < 5 {
		boltHeight = 5
	}
	return &entity.Entity{
		Type:            entity.TypeOverlay,
		X:               x,
		Y:               0,
		Layer:           render.LayerParticles + 5, // above rain
		Sprite:          lightningSprite,
		Visible:         true,
		Lifetime:        300 * time.Millisecond,
		Behavior:        &LightningBehavior{BoltHeight: boltHeight},
		RenderFunc:      LightningRenderFunc,
		Data:            &LightningData{Duration: 300 * time.Millisecond, BoltRune: '/'},
		RemoveOffscreen: false,
	}
}

// LightningBehavior updates the lightning flash. It has no movement; it just
// exists for its lifetime and then expires.
type LightningBehavior struct {
	BoltHeight int
}

// Update is a no-op; the lightning flash simply expires after its lifetime.
func (b *LightningBehavior) Update(e *entity.Entity, _ entity.UpdateContext, _ time.Duration) {
	// No movement needed. The manager handles lifetime expiration.
}

// LightningRenderFunc draws a jagged lightning bolt from the top of the screen
// down to a fraction of the viewport height.
func LightningRenderFunc(canvas *render.Canvas, e *entity.Entity) {
	data, ok := e.Data.(*LightningData)
	if !ok || data == nil {
		return
	}
	x := int(e.X)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	bolt := data.BoltRune

	// Draw a jagged vertical bolt.
	for y := 0; y < 20; y++ {
		// Alternate the X position slightly for a jagged effect.
		offsetX := 0
		if y%3 == 1 {
			offsetX = 1
		} else if y%3 == 2 {
			offsetX = -1
		}
		canvas.SetRune(x+offsetX, y, bolt, style)
	}
}

// LightningFlash creates a brief full-screen brightening effect entity. This
// is separate from the bolt itself and makes the whole screen flash white for
// a single frame.
type LightningFlash struct {
	Active   bool
	Elapsed  time.Duration
	Duration time.Duration
}

// NewLightningFlash creates a LightningFlash state object (not an entity).
func NewLightningFlash() *LightningFlash {
	return &LightningFlash{
		Duration: 150 * time.Millisecond,
	}
}

// Trigger activates the flash.
func (f *LightningFlash) Trigger() {
	f.Active = true
	f.Elapsed = 0
}

// Update advances the flash timer. Returns true while the flash is still
// active.
func (f *LightningFlash) Update(delta time.Duration) bool {
	if !f.Active {
		return false
	}
	f.Elapsed += delta
	if f.Elapsed >= f.Duration {
		f.Active = false
		return false
	}
	return true
}

// Intensity returns the current flash intensity (0.0 to 1.0) based on how
// much of the flash duration has elapsed. The flash is brightest at the
// start and fades out.
func (f *LightningFlash) Intensity() float64 {
	if !f.Active {
		return 0
	}
	remaining := float64(f.Duration-f.Elapsed) / float64(f.Duration)
	if remaining < 0 {
		return 0
	}
	if remaining > 1 {
		return 1
	}
	return remaining
}
