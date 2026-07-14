package scene

import (
	"math/rand"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// Countryside is a rural scene with hills, trees, small houses, utility poles,
// and drifting clouds.
type Countryside struct {
	vp    Viewport
	style tcell.Style
}

// Name returns "countryside".
func (c *Countryside) Name() string { return "countryside" }

// Type returns SceneCountryside.
func (c *Countryside) Type() SceneType { return SceneCountryside }

// Build generates the countryside layout and populates the entity manager
// with cloud entities. Static scenery (hills, trees, houses, poles) is drawn
// directly in RenderBackground for performance.
func (c *Countryside) Build(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style) {
	c.vp = vp
	c.style = style
	// Remove old scenery entities.
	for _, e := range manager.ByType(entity.TypeCloud) {
		e.Dead = true
	}
	manager.Flush()
	// Spawn clouds.
	spawnClouds(manager, vp, rng)
}

// Update advances scene-specific state. Clouds are updated by the entity
// manager; the countryside scene has no additional per-frame logic.
func (c *Countryside) Update(_ *entity.Manager, _ Viewport, _ float64) {}

// RenderBackground draws the sky, hills, trees, houses, and utility poles.
func (c *Countryside) RenderBackground(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	c.vp = vp
	c.style = style
	rng := rand.New(rand.NewSource(int64(vp.Width)*7 + int64(vp.Height)*13))

	// Draw hills as the base of the scenery.
	drawHills(canvas, vp, style, rng)

	// Draw scattered trees on the hills.
	treeCount := vp.Width / 20
	if treeCount < 3 {
		treeCount = 3
	}
	drawTrees(canvas, vp, style, rng, treeCount)

	// Draw a few small houses.
	c.drawHouses(canvas, vp, style, rng)

	// Draw utility poles along the ground.
	c.drawPoles(canvas, vp, style)

	// Draw ground texture below the track.
	drawGround(canvas, vp, style, '"')
}

// RenderBackgroundWithPalette draws the background with time-of-day palette.
func (c *Countryside) RenderBackgroundWithPalette(canvas *render.Canvas, vp Viewport, pal Palette, style tcell.Style) {
	c.vp = vp
	c.style = style
	rng := rand.New(rand.NewSource(int64(vp.Width)*7 + int64(vp.Height)*13))

	// Draw celestial elements (sun/moon/stars).
	drawCelestial(canvas, vp, pal, rng)

	// Draw hills with palette colors.
	hillStyle := style.Foreground(pal.HillColor)
	if pal.DimFactor {
		hillStyle = hillStyle.Dim(true)
	}
	drawHills(canvas, vp, hillStyle, rng)

	// Draw scattered trees on the hills.
	treeCount := vp.Width / 20
	if treeCount < 3 {
		treeCount = 3
	}
	treeStyle := style.Foreground(pal.TreeColor)
	drawTrees(canvas, vp, treeStyle, rng, treeCount)

	// Draw a few small houses.
	c.drawHouses(canvas, vp, style, rng)

	// Draw utility poles along the ground.
	c.drawPoles(canvas, vp, style)

	// Draw ground texture below the track.
	groundStyle := style.Foreground(pal.GroundColor)
	if pal.DimFactor {
		groundStyle = groundStyle.Dim(true)
	}
	drawGroundStyled(canvas, vp, groundStyle, '"')
}

// RenderTrack draws the standard railway track.
func (c *Countryside) RenderTrack(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	drawTrack(canvas, vp, style)
}

// HandleResize updates the viewport and regenerates cloud positions.
func (c *Countryside) HandleResize(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style) {
	c.Build(manager, vp, rng, style)
}

// drawHouses draws 1-3 small houses on the hills.
func (c *Countryside) drawHouses(canvas *render.Canvas, vp Viewport, style tcell.Style, rng *rand.Rand) {
	houseCount := 1 + rng.Intn(3)
	houseY := vp.TrackY() - 5
	if houseY < 2 {
		houseY = 2
	}
	houseStyle := style.Dim(true)
	roofStyle := style.Foreground(tcell.ColorDarkRed).Dim(true)

	for i := 0; i < houseCount; i++ {
		x := 10 + rng.Intn(vp.Width-20)
		if x < 0 {
			x = 0
		}
		y := houseY - rng.Intn(2)
		if y < 0 {
			y = 0
		}
		// Roof.
		canvas.SetRune(x+1, y, '/', roofStyle)
		canvas.SetRune(x+2, y+1, '\\', roofStyle)
		canvas.SetRune(x, y+1, '/', roofStyle)
		canvas.SetRune(x+3, y+1, '\\', roofStyle)
		// Walls.
		for row := y + 2; row < y+5 && row < vp.Height; row++ {
			canvas.SetRune(x, row, '|', houseStyle)
			canvas.SetRune(x+3, row, '|', houseStyle)
		}
		// Window.
		canvas.SetRune(x+1, y+3, 'o', houseStyle)
		canvas.SetRune(x+2, y+3, 'o', houseStyle)
	}
}

// drawPoles draws utility poles at regular intervals along the ground.
func (c *Countryside) drawPoles(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	poleStyle := style.Dim(true)
	poleY := vp.TrackY() - 3
	if poleY < 2 {
		poleY = 2
	}
	for x := 15; x < vp.Width; x += 25 {
		// Pole.
		for y := poleY; y < vp.TrackY() && y < vp.Height; y++ {
			canvas.SetRune(x, y, '|', poleStyle)
		}
		// Crossbar.
		canvas.SetRune(x-1, poleY, '-', poleStyle)
		canvas.SetRune(x+1, poleY, '-', poleStyle)
	}
}
