package scene

import (
	"math/rand"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// Mountain is a mountain route scene with tall mountains, cliffs, a tunnel
// entrance, a bridge, a river, and pine trees.
type Mountain struct {
	vp    Viewport
	style tcell.Style
}

// Name returns "mountain".
func (m *Mountain) Name() string { return "mountain" }

// Type returns SceneMountain.
func (m *Mountain) Type() SceneType { return SceneMountain }

// Build generates the mountain layout and populates the entity manager with
// cloud entities.
func (m *Mountain) Build(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style) {
	m.vp = vp
	m.style = style
	// Remove old scenery entities.
	for _, e := range manager.ByType(entity.TypeCloud) {
		e.Dead = true
	}
	manager.Flush()
	// Spawn clouds (mountains are high, so clouds appear lower).
	spawnClouds(manager, vp, rng)
}

// Update advances scene-specific state. Clouds are updated by the entity
// manager; the mountain scene has no additional per-frame logic.
func (m *Mountain) Update(_ *entity.Manager, _ Viewport, _ float64) {}

// RenderBackground draws the mountains, cliffs, tunnel, bridge, river, and
// pine trees.
func (m *Mountain) RenderBackground(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	m.vp = vp
	m.style = style
	rng := rand.New(rand.NewSource(int64(vp.Width)*19 + int64(vp.Height)*23))

	// Draw mountains (large, in the background).
	m.drawMountains(canvas, vp, style, rng)

	// Draw pine trees on the slopes.
	treeCount := vp.Width / 15
	if treeCount < 4 {
		treeCount = 4
	}
	drawTrees(canvas, vp, style, rng, treeCount)

	// Draw tunnel entrance on the left side.
	m.drawTunnel(canvas, vp, style)

	// Draw bridge in the middle.
	m.drawBridge(canvas, vp, style)

	// Draw river below the bridge.
	m.drawRiver(canvas, vp, style)

	// Draw rocky ground.
	drawGround(canvas, vp, style, '.')
}

// RenderBackgroundWithPalette draws the background with time-of-day palette.
func (m *Mountain) RenderBackgroundWithPalette(canvas *render.Canvas, vp Viewport, pal Palette, style tcell.Style) {
	m.vp = vp
	m.style = style
	rng := rand.New(rand.NewSource(int64(vp.Width)*19 + int64(vp.Height)*23))

	// Draw celestial elements (sun/moon/stars).
	drawCelestial(canvas, vp, pal, rng)

	// Draw mountains with palette.
	mountainStyle := style.Foreground(pal.HillColor)
	if pal.DimFactor {
		mountainStyle = mountainStyle.Dim(true)
	}
	m.drawMountains(canvas, vp, mountainStyle, rng)

	// Draw pine trees with palette.
	treeCount := vp.Width / 15
	if treeCount < 4 {
		treeCount = 4
	}
	treeStyle := style.Foreground(pal.TreeColor)
	drawTrees(canvas, vp, treeStyle, rng, treeCount)

	// Draw tunnel entrance on the left side.
	m.drawTunnel(canvas, vp, style)

	// Draw bridge in the middle.
	m.drawBridge(canvas, vp, style)

	// Draw river below the bridge.
	m.drawRiver(canvas, vp, style)

	// Draw rocky ground with palette.
	groundStyle := style.Foreground(pal.GroundColor)
	if pal.DimFactor {
		groundStyle = groundStyle.Dim(true)
	}
	drawGroundStyled(canvas, vp, groundStyle, '.')
}

// RenderTrack draws the railway track. In the mountain scene, the track
// passes through the tunnel and over the bridge.
func (m *Mountain) RenderTrack(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	drawTrack(canvas, vp, style)
}

// HandleResize regenerates the mountain layout for the new viewport.
func (m *Mountain) HandleResize(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style) {
	m.Build(manager, vp, rng, style)
}

// drawMountains draws large triangular mountain shapes in the background.
func (m *Mountain) drawMountains(canvas *render.Canvas, vp Viewport, style tcell.Style, rng *rand.Rand) {
	mountainStyle := style.Dim(true)
	peakStyle := style.Foreground(tcell.ColorWhite).Dim(true)
	baseY := vp.TrackY() - 3
	if baseY < 2 {
		baseY = 2
	}

	// Draw 2-3 mountain peaks.
	peaks := []int{vp.Width / 4, vp.Width / 2, vp.Width * 3 / 4}
	for _, peakX := range peaks {
		height := 8 + rng.Intn(5)
		if height > baseY {
			height = baseY
		}
		// Draw triangular mountain.
		for row := 0; row < height; row++ {
			y := baseY - row
			if y < 0 {
				break
			}
			halfWidth := row + 2
			for col := -halfWidth; col <= halfWidth; col++ {
				x := peakX + col
				if x < 0 || x >= vp.Width {
					continue
				}
				if row >= height-2 {
					canvas.SetRune(x, y, '^', peakStyle)
				} else {
					canvas.SetRune(x, y, '/', mountainStyle)
				}
			}
		}
		// Snow cap.
		canvas.SetRune(peakX, baseY-height, '^', peakStyle)
	}
}

// drawTunnel draws a tunnel entrance on the left side of the track.
func (m *Mountain) drawTunnel(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	tunnelStyle := style.Dim(true)
	trackY := vp.TrackY()
	tunnelX := 0
	tunnelWidth := 8
	tunnelHeight := 6

	// Arch: draw an inverted-U shape.
	for col := 0; col < tunnelWidth; col++ {
		x := tunnelX + col
		if x >= vp.Width {
			break
		}
		// Top of the arch.
		canvas.SetRune(x, trackY-tunnelHeight, '=', tunnelStyle)
		// Sides.
		if col < 2 {
			for row := 1; row < tunnelHeight; row++ {
				y := trackY - tunnelHeight + row
				if y >= 0 && y < vp.Height {
					canvas.SetRune(x, y, '|', tunnelStyle)
				}
			}
		}
		if col >= tunnelWidth-2 {
			for row := 1; row < tunnelHeight; row++ {
				y := trackY - tunnelHeight + row
				if y >= 0 && y < vp.Height {
					canvas.SetRune(x, y, '|', tunnelStyle)
				}
			}
		}
	}
	// Dark interior.
	for col := 2; col < tunnelWidth-2; col++ {
		for row := 1; row < tunnelHeight; row++ {
			x := tunnelX + col
			y := trackY - tunnelHeight + row
			if x >= 0 && x < vp.Width && y >= 0 && y < vp.Height {
				canvas.SetRune(x, y, '#', tcell.StyleDefault.Reverse(true))
			}
		}
	}
}

// drawBridge draws a bridge structure in the center of the viewport.
func (m *Mountain) drawBridge(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	bridgeStyle := style.Dim(true)
	trackY := vp.TrackY()
	bridgeX := vp.Width/2 - 6
	bridgeWidth := 12

	// Bridge supports.
	for col := 0; col < bridgeWidth; col += 5 {
		x := bridgeX + col
		if x < 0 || x >= vp.Width {
			continue
		}
		for y := trackY + 2; y < vp.Height && y < trackY+6; y++ {
			canvas.SetRune(x, y, '|', bridgeStyle)
		}
	}
	// Bridge deck (just below the track).
	deckY := trackY + 2
	if deckY < vp.Height {
		for col := 0; col < bridgeWidth; col++ {
			x := bridgeX + col
			if x >= 0 && x < vp.Width {
				canvas.SetRune(x, deckY, '=', bridgeStyle)
			}
		}
	}
	// Trusses (X pattern above the track).
	trussY := trackY - 3
	if trussY >= 0 {
		for col := 0; col < bridgeWidth; col++ {
			x := bridgeX + col
			if x < 0 || x >= vp.Width {
				continue
			}
			if col%4 < 2 {
				canvas.SetRune(x, trussY, '\\', bridgeStyle)
			} else {
				canvas.SetRune(x, trussY, '/', bridgeStyle)
			}
		}
		// Top rail.
		canvas.SetRune(bridgeX, trussY-1, '=', bridgeStyle)
		for col := 1; col < bridgeWidth; col++ {
			x := bridgeX + col
			if x >= 0 && x < vp.Width {
				canvas.SetRune(x, trussY-1, '=', bridgeStyle)
			}
		}
	}
}

// drawRiver draws a river below the bridge area.
func (m *Mountain) drawRiver(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	riverStyle := style.Foreground(tcell.ColorBlue).Dim(true)
	riverY := vp.TrackY() + 5
	if riverY >= vp.Height {
		return
	}
	riverWidth := 12
	riverX := vp.Width/2 - 6

	for col := 0; col < riverWidth; col++ {
		x := riverX + col
		if x < 0 || x >= vp.Width {
			continue
		}
		if riverY < vp.Height {
			canvas.SetRune(x, riverY, '~', riverStyle)
		}
		if riverY+1 < vp.Height {
			canvas.SetRune(x, riverY+1, '~', riverStyle)
		}
	}
}
