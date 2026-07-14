package engine

import (
	"math/rand"
	"time"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/effects"
	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/example/conductor/internal/scene"
	"github.com/example/conductor/internal/train"
	"github.com/gdamore/tcell/v2"
)

// World holds the complete simulation state for the current scene. It owns an
// entity Manager, a Scene, a weather system, a time-of-day setting, a random
// event scheduler, and a train consist generator.
type World struct {
	width  int
	height int
	trackY int

	manager   *entity.Manager
	scene     scene.Scene
	weather   *effects.WeatherSystem
	scheduler *Scheduler
	style     tcell.Style
	seed      int64
	rng       *rand.Rand
	generator *train.Generator

	timePeriod    scene.TimePeriod
	trainSpeed    float64
	noColor       bool
	reducedMotion bool

	spawnDelay   time.Duration
	spawnElapsed time.Duration
	waiting      bool

	paused      bool
	helpVisible bool

	cfg config.Config
}

// NewWorld creates a world sized to the given viewport. The scene and weather
// are chosen randomly from the seed. A first train is spawned immediately.
func NewWorld(width, height int, seed int64, style tcell.Style) *World {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	w := &World{
		width:      width,
		height:     height,
		style:      style,
		seed:       seed,
		rng:        rand.New(rand.NewSource(seed)),
		trainSpeed: 25.0,
		spawnDelay: 2 * time.Second,
		manager:    entity.NewManager(),
		generator:  train.NewGenerator(seed, style),
		timePeriod: scene.TimeDay,
		weather:    effects.NewWeatherSystem(effects.WeatherClear, false),
		scheduler:  NewScheduler(seed),
	}
	w.layoutTrack()
	w.scene = scene.NewRandom(w.rng)
	w.scene.Build(w.manager, w.viewport(), w.rng, style)
	w.scheduler.Reset(seed)
	w.registerEventCallbacks()
	w.spawnTrain()
	w.manager.Flush()
	return w
}

// NewWorldWithConfig creates a world configured by the given Config. It
// applies train type, scene, weather, time, speed, no-color, and
// reduced-motion settings from the config.
func NewWorldWithConfig(width, height int, cfg config.Config, style tcell.Style) *World {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	noColor := cfg.NoColor || cfg.Color == "16" && false // color "16" still has color
	if cfg.NoColor {
		noColor = true
	}

	w := &World{
		width:         width,
		height:        height,
		style:         style,
		seed:          cfg.Seed,
		rng:           rand.New(rand.NewSource(cfg.Seed)),
		trainSpeed:    25.0 * cfg.Speed,
		spawnDelay:    2 * time.Second,
		manager:       entity.NewManager(),
		generator:     train.NewGenerator(cfg.Seed, style),
		timePeriod:    resolveTimePeriod(cfg.Time, rand.New(rand.NewSource(cfg.Seed+1))),
		weather:       effects.NewWeatherSystem(resolveWeatherType(cfg.Weather, rand.New(rand.NewSource(cfg.Seed+2))), cfg.ReducedMotion),
		scheduler:     NewScheduler(cfg.Seed),
		noColor:       noColor,
		reducedMotion: cfg.ReducedMotion,
		cfg:           cfg,
	}
	w.layoutTrack()

	// Resolve scene from config.
	w.scene = resolveScene(cfg.Scene, w.rng)
	w.scene.Build(w.manager, w.viewport(), w.rng, style)

	// Apply weather particles if weather is not clear.
	w.weather.SetWeather(w.weather.Current(), w.manager, w.width, w.height)

	w.scheduler.Reset(cfg.Seed)
	if !cfg.RandomEvents {
		w.scheduler.SetEnabled(false)
	}
	w.registerEventCallbacks()
	w.spawnTrain()
	w.manager.Flush()
	return w
}

// resolveTimePeriod converts a config string to a TimePeriod, choosing
// randomly for "auto".
func resolveTimePeriod(s string, rng *rand.Rand) scene.TimePeriod {
	tp, ok := scene.ParseTimePeriod(s)
	if ok && tp >= 0 {
		return tp
	}
	// "auto" or unknown: random.
	periods := scene.AllTimePeriods()
	return periods[rng.Intn(len(periods))]
}

// resolveWeatherType converts a config string to a WeatherType, choosing
// randomly for "random".
func resolveWeatherType(s string, rng *rand.Rand) effects.WeatherType {
	wt, ok := effects.ParseWeatherType(s)
	if ok && wt >= 0 {
		return wt
	}
	// "random" or unknown: choose randomly.
	types := effects.AllWeatherTypes()
	return types[rng.Intn(len(types))]
}

// resolveScene converts a config string to a Scene, choosing randomly for
// "random".
func resolveScene(s string, rng *rand.Rand) scene.Scene {
	st, ok := scene.ParseSceneType(s)
	if ok {
		return scene.New(st)
	}
	return scene.NewRandom(rng)
}

// registerEventCallbacks wires scheduler events to world actions.
func (w *World) registerEventCallbacks() {
	w.scheduler.On(EventLongFreight, func(_ EventType, rng *rand.Rand) {
		w.spawnLongFreight(rng)
	})
	w.scheduler.On(EventBirdsCrossing, func(_ EventType, rng *rand.Rand) {
		w.spawnBirds(rng)
	})
	w.scheduler.On(EventLightning, func(_ EventType, rng *rand.Rand) {
		w.weather.MaybeLightning(rng)
	})
	w.scheduler.On(EventSignalChange, func(_ EventType, _ *rand.Rand) {
		w.toggleRandomSignal()
	})
	w.scheduler.On(EventSmokeBurst, func(_ EventType, _ *rand.Rand) {
		w.spawnSmokeBurst()
	})
}

// viewport returns the current viewport for scene calls.
func (w *World) viewport() scene.Viewport {
	return scene.Viewport{Width: w.width, Height: w.height}
}

// Width returns the viewport width.
func (w *World) Width() int { return w.width }

// Height returns the viewport height.
func (w *World) Height() int { return w.height }

// TrackY returns the Y coordinate of the top rail row.
func (w *World) TrackY() int { return w.trackY }

// Scene returns the current scene.
func (w *World) Scene() scene.Scene { return w.scene }

// TimePeriod returns the current time of day.
func (w *World) TimePeriod() scene.TimePeriod { return w.timePeriod }

// Weather returns the current weather type.
func (w *World) Weather() effects.WeatherType { return w.weather.Current() }

// Paused reports whether the world is paused.
func (w *World) Paused() bool { return w.paused }

// SetPaused pauses or resumes the world.
func (w *World) SetPaused(paused bool) {
	w.paused = paused
}

// SetSpeed sets the train speed in cells per second.
func (w *World) SetSpeed(cellsPerSecond float64) {
	w.trainSpeed = cellsPerSecond
	for _, e := range w.manager.ByType(entity.TypeTrain) {
		dir := train.LeftToRight
		if data, ok := e.Data.(*train.TrainData); ok && data != nil {
			dir = data.Direction
		}
		e.VX = cellsPerSecond * train.VelocityForDirection(dir)
	}
}

// Speed returns the current train speed in cells per second.
func (w *World) Speed() float64 { return w.trainSpeed }

// Regenerate recreates the world from scratch using the current seed.
func (w *World) Regenerate() {
	w.manager = entity.NewManager()
	w.rng = rand.New(rand.NewSource(w.seed))
	w.generator = train.NewGenerator(w.seed, w.style)
	w.waiting = false
	w.spawnElapsed = 0
	w.scene.Build(w.manager, w.viewport(), w.rng, w.style)
	w.scheduler.Reset(w.seed)
	w.weather.SetWeather(w.weather.Current(), w.manager, w.width, w.height)
	w.spawnTrain()
	w.manager.Flush()
}

// SpawnNext forces a new train to spawn immediately.
func (w *World) SpawnNext() {
	for _, e := range w.manager.ByType(entity.TypeTrain) {
		e.Dead = true
	}
	w.manager.Flush()
	w.spawnTrain()
	w.manager.Flush()
}

// ChangeScene cycles to the next scene type.
func (w *World) ChangeScene() {
	types := scene.AllSceneTypes()
	current := w.scene.Type()
	nextIdx := 0
	for i, t := range types {
		if t == current {
			nextIdx = (i + 1) % len(types)
			break
		}
	}
	w.scene = scene.New(types[nextIdx])
	for _, e := range w.manager.ByType(entity.TypeCloud) {
		e.Dead = true
	}
	for _, e := range w.manager.ByType(entity.TypeSignal) {
		e.Dead = true
	}
	w.manager.Flush()
	w.scene.Build(w.manager, w.viewport(), w.rng, w.style)
	w.manager.Flush()
}

// ChangeWeather cycles to the next weather type (clear → rain → snow → clear).
func (w *World) ChangeWeather() {
	w.weather.CycleWeather(w.manager, w.width, w.height)
	w.manager.Flush()
}

// ChangeTime cycles to the next time period (morning → day → sunset → night).
func (w *World) ChangeTime() {
	periods := scene.AllTimePeriods()
	nextIdx := 0
	for i, p := range periods {
		if p == w.timePeriod {
			nextIdx = (i + 1) % len(periods)
			break
		}
	}
	w.timePeriod = periods[nextIdx]
}

// HandleResize updates the viewport dimensions and recalculates the track
// position. The scene and weather are rebuilt for the new dimensions.
func (w *World) HandleResize(width, height int) {
	w.width = width
	w.height = height
	w.layoutTrack()
	w.scene.HandleResize(w.manager, w.viewport(), w.rng, w.style)
	w.weather.HandleResize(w.manager, w.width, w.height)
	for _, e := range w.manager.ByType(entity.TypeTrain) {
		e.Y = float64(w.trackY)
	}
}

// layoutTrack places the track at roughly 70% of the terminal height.
func (w *World) layoutTrack() {
	w.trackY = w.height * 7 / 10
	if w.trackY < 3 {
		w.trackY = 3
	}
}

// spawnTrain creates a new train entity entering from the appropriate edge.
func (w *World) spawnTrain() {
	dir := train.LeftToRight
	if w.rng.Intn(2) == 0 {
		dir = train.RightToLeft
	}
	trainType, cars := w.generator.GenerateRandom(dir, w.width)
	trainW := train.TrainWidth(cars)
	var startX float64
	if dir == train.LeftToRight {
		startX = float64(-trainW)
	} else {
		startX = float64(w.width)
	}
	speed := train.TrainSpeed(trainType)
	if w.trainSpeed != 25.0 {
		speed = w.trainSpeed
	}
	var smokeFactory train.SmokeFactory
	if train.EmitsSmoke(trainType) {
		smokeFactory = train.SmokeFactory(effects.NewSmoke)
	}
	e := train.NewTrainEntity(startX, w.trackY, speed, dir, cars, smokeFactory)
	w.manager.Add(e)
	w.waiting = false
}

// spawnLongFreight creates a rare long freight train event.
func (w *World) spawnLongFreight(rng *rand.Rand) {
	dir := train.LeftToRight
	if rng.Intn(2) == 0 {
		dir = train.RightToLeft
	}
	trainType, cars := w.generator.GenerateLongFreight(dir, w.width)
	trainW := train.TrainWidth(cars)
	var startX float64
	if dir == train.LeftToRight {
		startX = float64(-trainW)
	} else {
		startX = float64(w.width)
	}
	speed := train.TrainSpeed(trainType)
	var smokeFactory train.SmokeFactory
	if train.EmitsSmoke(trainType) {
		smokeFactory = train.SmokeFactory(effects.NewSmoke)
	}
	e := train.NewTrainEntity(startX, w.trackY, speed, dir, cars, smokeFactory)
	w.manager.Add(e)
}

// spawnBirds creates a flock of birds flying across the sky.
func (w *World) spawnBirds(rng *rand.Rand) {
	count := 3 + rng.Intn(4)
	startX := -5
	vx := 15.0 + rng.Float64()*10.0
	if rng.Intn(2) == 0 {
		startX = w.width + 5
		vx = -vx
	}
	skyY := w.trackY / 3
	if skyY < 2 {
		skyY = 2
	}
	for i := 0; i < count; i++ {
		y := skyY + rng.Intn(5)
		if y >= w.trackY {
			y = w.trackY - 1
		}
		w.manager.Add(effects.NewBird(startX+i*4, y, vx))
	}
}

// toggleRandomSignal flips the state of a random signal entity.
func (w *World) toggleRandomSignal() {
	signals := w.manager.ByType(entity.TypeSignal)
	if len(signals) == 0 {
		return
	}
	sig := signals[w.rng.Intn(len(signals))]
	if data, ok := sig.Data.(*scene.SignalData); ok && data != nil {
		if data.State == scene.SignalRed {
			data.State = scene.SignalGreen
		} else {
			data.State = scene.SignalRed
		}
		data.Elapsed = 0
	}
}

// spawnSmokeBurst emits a cluster of smoke particles from the active train's
// chimney.
func (w *World) spawnSmokeBurst() {
	trains := w.manager.ByType(entity.TypeTrain)
	if len(trains) == 0 {
		return
	}
	e := trains[0]
	data, ok := e.Data.(*train.TrainData)
	if !ok || data == nil || len(data.Cars) == 0 {
		return
	}
	loco := data.Cars[0]
	spriteW := loco.Sprite.Width
	chimneyOffset := spriteW - 3
	if data.Direction == train.RightToLeft {
		chimneyOffset = 2
	}
	smokeX := e.X + float64(chimneyOffset)
	smokeY := e.Y - 1
	for i := 0; i < 5; i++ {
		w.manager.Add(effects.NewSmoke(smokeX+float64(i-2), smokeY))
	}
}

// Update advances the simulation by delta seconds.
func (w *World) Update(delta float64) {
	if w.paused {
		return
	}
	dt := time.Duration(delta * float64(time.Second))
	w.manager.Update(w.width, w.height, dt)
	w.scene.Update(w.manager, w.viewport(), delta)
	w.scheduler.Update(dt)

	// Update lightning flash.
	w.weather.Flash().Update(dt)

	// Check if the train has been culled (offscreen or dead).
	trains := w.manager.ByType(entity.TypeTrain)
	if len(trains) == 0 {
		if !w.waiting {
			w.waiting = true
			w.spawnElapsed = 0
		}
	} else {
		w.waiting = false
	}

	if w.waiting {
		w.spawnElapsed += dt
		if w.spawnElapsed >= w.spawnDelay {
			w.spawnTrain()
		}
	}
}

// Render draws the scene background, track, and all entities onto the canvas.
func (w *World) Render(canvas *render.Canvas) {
	pal := scene.PaletteFor(w.timePeriod)
	if w.noColor {
		pal = monochromePalette()
	}
	w.scene.RenderBackgroundWithPalette(canvas, w.viewport(), pal, w.style)
	w.scene.RenderTrack(canvas, w.viewport(), w.style)
	w.manager.Render(canvas)
	if w.weather.Flash().Active {
		w.renderLightningFlash(canvas, w.weather.Flash().Intensity())
	}
	if w.paused {
		w.renderPausedIndicator(canvas)
	}
	if w.helpVisible {
		w.renderHelpOverlay(canvas)
	}
}

// monochromePalette returns a palette with no color attributes, suitable for
// no-color mode. All elements use default colors with dim/reverse for
// contrast.
func monochromePalette() scene.Palette {
	return scene.Palette{
		DimFactor: false,
		ShowStars: true,
		ShowSun:   true,
		ShowMoon:  true,
	}
}

// renderLightningFlash overlays a brief brightening on the canvas.
func (w *World) renderLightningFlash(canvas *render.Canvas, intensity float64) {
	if intensity <= 0 {
		return
	}
	// Render a subtle flash by adding a faint white character overlay
	// at a few random positions. This is intentionally lightweight to
	// avoid performance impact.
	flashStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Dim(true)
	for i := 0; i < w.width*w.height/50; i++ {
		x := w.rng.Intn(w.width)
		y := w.rng.Intn(w.height)
		canvas.SetRune(x, y, '.', flashStyle)
	}
}

// renderPausedIndicator draws a small [PAUSED] label in the top-right corner.
func (w *World) renderPausedIndicator(canvas *render.Canvas) {
	label := "[PAUSED]"
	x := w.width - len(label) - 2
	if x < 0 {
		x = 0
	}
	style := tcell.StyleDefault.Reverse(true)
	for i, r := range label {
		canvas.SetRune(x+i, 0, r, style)
	}
}

// HelpVisible reports whether the help overlay is currently shown.
func (w *World) HelpVisible() bool { return w.helpVisible }

// Trains returns all train entities in the world.
func (w *World) Trains() []*entity.Entity {
	return w.manager.ByType(entity.TypeTrain)
}

// EntityCount returns the number of active entities.
func (w *World) EntityCount() int {
	return w.manager.Count()
}

// PendingEvents returns the number of pending scheduler events.
func (w *World) PendingEvents() int {
	return w.scheduler.PendingCount()
}

// ToggleHelp shows or hides the help overlay.
func (w *World) ToggleHelp() {
	w.helpVisible = !w.helpVisible
}

// SetHelpVisible explicitly sets the help overlay visibility.
func (w *World) SetHelpVisible(visible bool) {
	w.helpVisible = visible
}

// NoColor reports whether color is disabled.
func (w *World) NoColor() bool { return w.noColor }

// ReducedMotion reports whether reduced-motion mode is active.
func (w *World) ReducedMotion() bool { return w.reducedMotion }

// ChangeColorPalette cycles through color palette modes. This is triggered by
// the 'c' key. In Milestone 6 this toggles between no-color and color modes.
func (w *World) ChangeColorPalette() {
	w.noColor = !w.noColor
}

// renderHelpOverlay draws the help text overlay in the center of the screen.
func (w *World) renderHelpOverlay(canvas *render.Canvas) {
	lines := []string{
		"Conductor — Keyboard Controls",
		"",
		"  q / Ctrl+C   Quit",
		"  p / Space    Pause or resume",
		"  r            Regenerate scene",
		"  n            Spawn next train",
		"  s            Change scene",
		"  w            Change weather",
		"  d            Change time of day",
		"  c            Toggle color / no-color",
		"  + / -        Increase / decrease speed",
		"  0            Reset speed",
		"  h / ?        Toggle this help",
		"  Esc          Close this help",
		"",
		"Press Esc or h to close.",
	}

	overlayW := 0
	for _, line := range lines {
		if len(line) > overlayW {
			overlayW = len(line)
		}
	}
	overlayH := len(lines)
	startX := w.width/2 - overlayW/2
	startY := w.height/2 - overlayH/2
	if startX < 1 {
		startX = 1
	}
	if startY < 1 {
		startY = 1
	}

	boxStyle := tcell.StyleDefault.Reverse(true)
	textStyle := tcell.StyleDefault

	// Draw border.
	for y := startY - 1; y < startY+overlayH+1 && y < w.height; y++ {
		for x := startX - 1; x < startX+overlayW+1 && x < w.width; x++ {
			if x == startX-1 || x == startX+overlayW {
				canvas.SetRune(x, y, '|', boxStyle)
			} else if y == startY-1 || y == startY+overlayH {
				canvas.SetRune(x, y, '-', boxStyle)
			} else {
				canvas.SetRune(x, y, ' ', textStyle)
			}
		}
	}

	// Draw text.
	for row, line := range lines {
		y := startY + row
		if y >= w.height {
			break
		}
		for col, r := range line {
			x := startX + col
			if x >= w.width {
				break
			}
			canvas.SetRune(x, y, r, textStyle)
		}
	}
}
