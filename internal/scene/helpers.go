package scene

import (
	"math/rand"
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// cloudArt defines a few cloud shapes. Clouds are drawn as entities so they
// can drift across the sky.
var cloudArts = [][]string{
	{
		`   .--.   `,
		`.(    ). `,
		` '-..-'  `,
	},
	{
		`  .-.    `,
		` (   ).  `,
		`'(___)'  `,
	},
	{
		`   .---.   `,
		` .'     '. `,
		`'---------'`,
	},
}

// CloudData holds state for a cloud entity.
type CloudData struct {
	ArtIndex int
}

// CloudBehavior slowly drifts clouds horizontally. Clouds wrap around when
// they exit the viewport.
type CloudBehavior struct{}

// Update moves the cloud and wraps it around the viewport edges.
func (CloudBehavior) Update(e *entity.Entity, ctx entity.UpdateContext, delta time.Duration) {
	sec := delta.Seconds()
	e.X += e.VX * sec
	// Wrap around: if a cloud exits the right edge, reappear on the left.
	spriteW := 10
	if e.Sprite != nil {
		spriteW = e.Sprite.Width
	}
	if e.VX > 0 && e.X > float64(ctx.Width) {
		e.X = float64(-spriteW)
	}
	if e.VX < 0 && e.X < float64(-spriteW) {
		e.X = float64(ctx.Width)
	}
}

// CloudRenderFunc draws a cloud using its art index stored in CloudData.
func CloudRenderFunc(canvas *render.Canvas, e *entity.Entity) {
	data, ok := e.Data.(*CloudData)
	if !ok || data == nil {
		return
	}
	if data.ArtIndex < 0 || data.ArtIndex >= len(cloudArts) {
		return
	}
	style := tcell.StyleDefault.Dim(true)
	art := cloudArts[data.ArtIndex]
	x, y := int(e.X), int(e.Y)
	for row, line := range art {
		for col, r := range line {
			if r != ' ' {
				canvas.SetRune(x+col, y+row, r, style)
			}
		}
	}
}

// NewCloud creates a cloud entity at the given position with a random art
// variant.
func NewCloud(x, y int, vx float64, rng *rand.Rand) *entity.Entity {
	artIdx := rng.Intn(len(cloudArts))
	return &entity.Entity{
		Type:       entity.TypeCloud,
		X:          float64(x),
		Y:          float64(y),
		VX:         vx,
		Layer:      render.LayerCelestial,
		Visible:    true,
		Behavior:   CloudBehavior{},
		RenderFunc: CloudRenderFunc,
		Data:       &CloudData{ArtIndex: artIdx},
	}
}

// spawnClouds adds a number of clouds to the entity manager, distributed
// across the sky. The number scales with viewport width.
func spawnClouds(manager *entity.Manager, vp Viewport, rng *rand.Rand) {
	count := vp.Width / 30
	if count < 2 {
		count = 2
	}
	if count > 6 {
		count = 6
	}
	for i := 0; i < count; i++ {
		x := rng.Intn(vp.Width)
		y := rng.Intn(vp.TrackY()/3) + 1
		if y < 1 {
			y = 1
		}
		vx := 2.0 + rng.Float64()*3.0
		if rng.Intn(2) == 0 {
			vx = -vx
		}
		manager.Add(NewCloud(x, y, vx, rng))
	}
}

// drawHills draws rolling hills at the base of the sky area. The hills are
// drawn directly on the canvas (not as entities) because they are static.
func drawHills(canvas *render.Canvas, vp Viewport, style tcell.Style, rng *rand.Rand) {
	hillY := vp.TrackY() - 4
	if hillY < 1 {
		hillY = 1
	}
	hillStyle := style.Dim(true)
	for x := 0; x < vp.Width; x++ {
		// Create a rolling hill silhouette using a simple pseudo-random wave.
		h := int(float64(rng.Intn(3)) * 0.5)
		baseY := hillY + h
		for y := baseY; y < vp.TrackY(); y++ {
			canvas.SetRune(x, y, '^', hillStyle)
		}
	}
}

// drawStars fills the sky with star characters. Used at night (Milestone 5
// will call this conditionally). Kept here as a shared helper.
func drawStars(canvas *render.Canvas, vp Viewport, rng *rand.Rand, density int) {
	starStyle := tcell.StyleDefault
	starY := vp.TrackY() - 2
	if starY < 1 {
		starY = 1
	}
	for i := 0; i < density; i++ {
		x := rng.Intn(vp.Width)
		y := rng.Intn(starY)
		canvas.SetRune(x, y, '*', starStyle)
	}
}

// drawGround fills the area below the track with a ground texture.
func drawGround(canvas *render.Canvas, vp Viewport, style tcell.Style, char rune) {
	groundStyle := style.Dim(true)
	groundY := vp.TrackY() + 2
	for y := groundY; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			canvas.SetRune(x, y, char, groundStyle)
		}
	}
}

// drawGroundStyled fills the area below the track with a ground texture using
// a caller-provided style (for palette-based rendering).
func drawGroundStyled(canvas *render.Canvas, vp Viewport, style tcell.Style, char rune) {
	groundY := vp.TrackY() + 2
	for y := groundY; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			canvas.SetRune(x, y, char, style)
		}
	}
}

// drawTrack draws the standard railway track across the full viewport width.
// All scenes use this for consistent track appearance.
func drawTrack(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	trackY := vp.TrackY()
	if trackY < 0 || trackY >= vp.Height {
		return
	}
	railStyle := style
	sleeperStyle := tcell.StyleDefault.Dim(true)
	for x := 0; x < vp.Width; x++ {
		canvas.SetRune(x, trackY, '=', railStyle)
	}
	sleeperY := trackY + 1
	if sleeperY < vp.Height {
		for x := 2; x < vp.Width; x += 4 {
			canvas.SetRune(x, sleeperY, 'o', sleeperStyle)
		}
	}
}

// treeArt is a simple pine tree shape used by countryside and mountain scenes.
var treeArt = []string{
	`   ^   `,
	`  ^^^  `,
	` ^^^^^ `,
	`^^^^^^^`,
	`   |   `,
}

// drawTree draws a single tree at the given position.
func drawTree(canvas *render.Canvas, x, y int, style tcell.Style) {
	treeStyle := style.Foreground(tcell.ColorGreen).Dim(true)
	for row, line := range treeArt {
		for col, r := range line {
			if r != ' ' {
				canvas.SetRune(x+col, y+row, r, treeStyle)
			}
		}
	}
}

// drawTrees scatters trees along the ground. The number and positions are
// deterministic from the RNG.
func drawTrees(canvas *render.Canvas, vp Viewport, style tcell.Style, rng *rand.Rand, count int) {
	groundY := vp.TrackY() - 5
	if groundY < 1 {
		groundY = 1
	}
	maxX := vp.Width - 7
	if maxX < 1 {
		maxX = 1
	}
	for i := 0; i < count; i++ {
		x := rng.Intn(maxX)
		// Trees sit at varying heights on the hills.
		y := groundY - rng.Intn(3)
		if y < 0 {
			y = 0
		}
		drawTree(canvas, x, y, style)
	}
}

// _ prevents unused import warnings during incremental development.
var _ = time.Millisecond
