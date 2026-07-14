package scene

import (
	"math/rand"
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// Station is a small railway station scene with a platform, roof, name board,
// clock, signals, and simplified passenger characters.
type Station struct {
	vp    Viewport
	style tcell.Style
}

// Name returns "station".
func (s *Station) Name() string { return "station" }

// Type returns SceneStation.
func (s *Station) Type() SceneType { return SceneStation }

// Build generates the station layout and populates the entity manager with
// signals and passengers.
func (s *Station) Build(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style) {
	s.vp = vp
	s.style = style
	// Remove old scenery entities.
	for _, e := range manager.ByType(entity.TypeSignal) {
		e.Dead = true
	}
	for _, e := range manager.ByType(entity.TypeCloud) {
		e.Dead = true
	}
	manager.Flush()

	// Spawn clouds (station is outdoors too).
	spawnClouds(manager, vp, rng)

	// Spawn entry and exit signals.
	trackY := vp.TrackY()
	signalY := trackY - 4
	if signalY < 1 {
		signalY = 1
	}
	// Entry signal (left side, red initially).
	manager.Add(NewSignal(8, signalY, SignalRed, 8*time.Second))
	// Exit signal (right side, green initially).
	manager.Add(NewSignal(vp.Width-9, signalY, SignalGreen, 10*time.Second))
}

// Update advances scene-specific state. Signals are updated by the entity
// manager; the station scene has no additional per-frame logic.
func (s *Station) Update(_ *entity.Manager, _ Viewport, _ float64) {}

// RenderBackground draws the sky, station building, platform, name board,
// clock, and passengers.
func (s *Station) RenderBackground(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	s.vp = vp
	s.style = style
	rng := rand.New(rand.NewSource(int64(vp.Width)*11 + int64(vp.Height)*17))

	// Draw distant hills.
	drawHills(canvas, vp, style, rng)

	// Draw the station building.
	s.drawStationBuilding(canvas, vp, style)

	// Draw the platform.
	s.drawPlatform(canvas, vp, style)

	// Draw passengers on the platform.
	s.drawPassengers(canvas, vp, style, rng)

	// Draw ground texture.
	drawGround(canvas, vp, style, ':')
}

// RenderBackgroundWithPalette draws the background with time-of-day palette.
func (s *Station) RenderBackgroundWithPalette(canvas *render.Canvas, vp Viewport, pal Palette, style tcell.Style) {
	s.vp = vp
	s.style = style
	rng := rand.New(rand.NewSource(int64(vp.Width)*11 + int64(vp.Height)*17))

	// Draw celestial elements (sun/moon/stars).
	drawCelestial(canvas, vp, pal, rng)

	// Draw distant hills with palette.
	hillStyle := style.Foreground(pal.HillColor)
	if pal.DimFactor {
		hillStyle = hillStyle.Dim(true)
	}
	drawHills(canvas, vp, hillStyle, rng)

	// Draw the station building.
	s.drawStationBuilding(canvas, vp, style)

	// Draw the platform.
	s.drawPlatform(canvas, vp, style)

	// Draw passengers on the platform.
	s.drawPassengers(canvas, vp, style, rng)

	// Draw ground texture with palette.
	groundStyle := style.Foreground(pal.GroundColor)
	if pal.DimFactor {
		groundStyle = groundStyle.Dim(true)
	}
	drawGroundStyled(canvas, vp, groundStyle, ':')
}

// RenderTrack draws the standard railway track.
func (s *Station) RenderTrack(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	drawTrack(canvas, vp, style)
}

// HandleResize regenerates the station layout for the new viewport.
func (s *Station) HandleResize(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style) {
	s.Build(manager, vp, rng, style)
}

// drawStationBuilding draws the station building with a roof, name board, and
// clock. The building is placed on the left side, above the platform.
func (s *Station) drawStationBuilding(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	buildX := 5
	buildY := vp.TrackY() - 7
	if buildY < 1 {
		buildY = 1
	}

	wallStyle := style.Dim(true)
	roofStyle := style.Foreground(tcell.ColorDarkRed).Dim(true)
	boardStyle := style.Bold(true)

	// Roof (triangular).
	for i := 0; i < 6; i++ {
		canvas.SetRune(buildX+i, buildY, '/', roofStyle)
	}
	canvas.SetRune(buildX+6, buildY, '^', roofStyle)
	for i := 0; i < 6; i++ {
		canvas.SetRune(buildX+7+i, buildY, '\\', roofStyle)
	}

	// Walls.
	for row := buildY + 1; row < buildY+5 && row < vp.Height; row++ {
		for col := buildX; col < buildX+13; col++ {
			if col == buildX || col == buildX+12 {
				canvas.SetRune(col, row, '|', wallStyle)
			} else {
				canvas.SetRune(col, row, '#', wallStyle)
			}
		}
	}

	// Name board above the roof.
	boardText := "STATION"
	boardX := buildX + 3
	boardY := buildY - 1
	if boardY >= 0 {
		for i, r := range boardText {
			canvas.SetRune(boardX+i, boardY, r, boardStyle)
		}
	}

	// Clock on the building wall.
	clockY := buildY + 2
	canvas.SetRune(buildX+2, clockY, 'O', boardStyle)
	canvas.SetRune(buildX+3, clockY, ':', boardStyle)
	canvas.SetRune(buildX+4, clockY, 'O', boardStyle)

	// Door.
	canvas.SetRune(buildX+9, buildY+3, '|', wallStyle)
	canvas.SetRune(buildX+10, buildY+3, '|', wallStyle)
	canvas.SetRune(buildX+9, buildY+4, '|', wallStyle)
	canvas.SetRune(buildX+10, buildY+4, '|', wallStyle)
}

// drawPlatform draws a rectangular platform along the track.
func (s *Station) drawPlatform(canvas *render.Canvas, vp Viewport, style tcell.Style) {
	platformStyle := style.Dim(true)
	platformY := vp.TrackY() + 2
	if platformY >= vp.Height {
		return
	}
	// Platform surface.
	for x := 0; x < vp.Width; x++ {
		canvas.SetRune(x, platformY, '_', platformStyle)
	}
	// Platform side.
	if platformY+1 < vp.Height {
		for x := 0; x < vp.Width; x++ {
			canvas.SetRune(x, platformY+1, '=', platformStyle)
		}
	}
}

// drawPassengers draws simple stick-figure passengers standing on the platform.
func (s *Station) drawPassengers(canvas *render.Canvas, vp Viewport, style tcell.Style, rng *rand.Rand) {
	passengerStyle := style
	platformY := vp.TrackY() + 1
	if platformY >= vp.Height {
		return
	}
	count := vp.Width / 25
	if count < 2 {
		count = 2
	}
	if count > 5 {
		count = 5
	}
	for i := 0; i < count; i++ {
		x := 15 + rng.Intn(vp.Width-30)
		if x < 0 {
			x = 0
		}
		// Head.
		canvas.SetRune(x, platformY-2, 'o', passengerStyle)
		// Body.
		canvas.SetRune(x, platformY-1, '|', passengerStyle)
		// Legs.
		canvas.SetRune(x-1, platformY, '/', passengerStyle)
		canvas.SetRune(x+1, platformY, '\\', passengerStyle)
	}
}
