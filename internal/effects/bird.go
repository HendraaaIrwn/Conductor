package effects

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
)

// birdFrames are the two animation frames for a flying bird. The wings
// alternate between up and down positions.
var birdFrame1 = render.CellGrid{
	{render.Cell{Rune: ' '}, render.Cell{Rune: 'v'}, render.Cell{Rune: ' '}},
}

var birdFrame2 = render.CellGrid{
	{render.Cell{Rune: '^'}, render.Cell{Rune: ' '}, render.Cell{Rune: '^'}},
}

var birdSprite = render.NewSprite("bird", birdFrame1, birdFrame2)

// BirdBehavior drives a bird entity: it flies horizontally across the screen
// and flaps its wings.
type BirdBehavior struct {
	FrameTime time.Duration
	Elapsed   time.Duration
}

// Update moves the bird horizontally and cycles the wing animation.
func (b *BirdBehavior) Update(e *entity.Entity, _ entity.UpdateContext, delta time.Duration) {
	sec := delta.Seconds()
	e.X += e.VX * sec

	b.Elapsed += delta
	if b.Elapsed >= b.FrameTime {
		b.Elapsed = 0
		e.Frame++
		if e.Frame >= birdSprite.FrameCount() {
			e.Frame = 0
		}
	}
}

// NewBird creates a bird entity at the given position, flying in the given
// direction. Birds fly across the screen and are removed when offscreen.
func NewBird(x, y int, vx float64) *entity.Entity {
	return &entity.Entity{
		Type:            entity.TypeBird,
		X:               float64(x),
		Y:               float64(y),
		VX:              vx,
		Layer:           render.LayerCelestial,
		Sprite:          birdSprite,
		Visible:         true,
		Behavior:        &BirdBehavior{FrameTime: 200 * time.Millisecond},
		RemoveOffscreen: true,
	}
}
