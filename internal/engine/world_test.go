package engine

import (
	"testing"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/example/conductor/internal/scene"
	"github.com/example/conductor/internal/train"
	"github.com/gdamore/tcell/v2"
)

func TestNewWorld(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	if w.Width() != 80 {
		t.Errorf("Width = %d, want 80", w.Width())
	}
	if w.Height() != 24 {
		t.Errorf("Height = %d, want 24", w.Height())
	}
	if w.TrackY() <= 0 {
		t.Errorf("TrackY = %d, want > 0", w.TrackY())
	}
	if w.manager.FirstByType(entity.TypeTrain) == nil {
		t.Error("new world should have a train entity")
	}
}

func TestWorldTrackYPlacement(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	if w.TrackY() != 16 {
		t.Errorf("TrackY = %d, want 16", w.TrackY())
	}
}

func TestWorldPauseStopsUpdate(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	train := w.manager.FirstByType(entity.TypeTrain)
	if train == nil {
		t.Fatal("expected a train entity")
	}
	initialX := train.X
	w.SetPaused(true)
	w.Update(1.0)
	if train.X != initialX {
		t.Errorf("paused train moved: X = %f, want %f", train.X, initialX)
	}
}

func TestWorldUpdateMovesTrain(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	train := w.manager.FirstByType(entity.TypeTrain)
	if train == nil {
		t.Fatal("expected a train entity")
	}
	initialX := train.X
	w.Update(1.0)
	if train.X == initialX {
		t.Error("train did not move after Update")
	}
}

func TestWorldRegenerate(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	if w.manager.FirstByType(entity.TypeTrain) == nil {
		t.Fatal("expected a train entity")
	}
	w.Regenerate()
	if w.manager.FirstByType(entity.TypeTrain) == nil {
		t.Fatal("no train after Regenerate")
	}
	// The regenerated train should start from an edge, not in the middle.
	x := w.manager.FirstByType(entity.TypeTrain).X
	if x > 0 && x < float64(w.Width()) {
		t.Errorf("regenerated train X = %f, should start from an edge", x)
	}
}

func TestWorldHandleResize(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	w.HandleResize(120, 40)
	if w.Width() != 120 {
		t.Errorf("Width = %d, want 120", w.Width())
	}
	if w.Height() != 40 {
		t.Errorf("Height = %d, want 40", w.Height())
	}
	if w.TrackY() != 28 {
		t.Errorf("TrackY = %d, want 28", w.TrackY())
	}
	train := w.manager.FirstByType(entity.TypeTrain)
	if train != nil {
		if train.Y != 28 {
			t.Errorf("train Y = %f, want 28", train.Y)
		}
	}
}

func TestWorldSetSpeed(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	w.SetSpeed(50.0)
	if w.Speed() != 50.0 {
		t.Errorf("Speed = %f, want 50.0", w.Speed())
	}
	train := w.manager.FirstByType(entity.TypeTrain)
	if train != nil {
		speed := train.VX
		if speed < 0 {
			speed = -speed
		}
		if speed != 50.0 {
			t.Errorf("train |VX| = %f, want 50.0", speed)
		}
	}
}

func TestWorldDeterministicSeed(t *testing.T) {
	w1 := NewWorld(80, 24, 42, tcell.StyleDefault)
	w2 := NewWorld(80, 24, 42, tcell.StyleDefault)
	t1 := w1.manager.FirstByType(entity.TypeTrain)
	t2 := w2.manager.FirstByType(entity.TypeTrain)
	if t1 == nil || t2 == nil {
		t.Fatal("both worlds should have trains")
	}
	d1 := t1.Data.(*train.TrainData)
	d2 := t2.Data.(*train.TrainData)
	if d1.Direction != d2.Direction {
		t.Error("same seed should produce same train direction")
	}
	if t1.X != t2.X {
		t.Error("same seed should produce same train start position")
	}
}

func TestWorldHasScene(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	if w.Scene() == nil {
		t.Fatal("world should have a scene")
	}
}

func TestWorldChangeScene(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	originalScene := w.Scene()
	w.ChangeScene()
	newScene := w.Scene()
	if newScene == nil {
		t.Fatal("no scene after ChangeScene")
	}
	if newScene.Type() == originalScene.Type() {
		t.Error("ChangeScene should switch to a different scene type")
	}
}

func TestWorldChangeSceneCyclesThroughAll(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	seen := map[scene.SceneType]bool{}
	for i := 0; i < 3; i++ {
		seen[w.Scene().Type()] = true
		w.ChangeScene()
	}
	// After 3 changes we should have seen at least 2 different scene types
	// (the initial scene plus at least one from cycling).
	if len(seen) < 2 {
		t.Errorf("only saw %d scene types after cycling, want >= 2", len(seen))
	}
}

func TestWorldRenderProducesOutput(t *testing.T) {
	w := NewWorld(120, 40, 42, tcell.StyleDefault)
	canvas := render.NewCanvas(120, 40)
	w.Render(canvas)
	// The rendered frame should have non-blank cells (track, scenery, train).
	nonBlank := 0
	for y := 0; y < 40; y++ {
		for x := 0; x < 120; x++ {
			if !canvas.CellAt(x, y).IsBlank() {
				nonBlank++
			}
		}
	}
	if nonBlank == 0 {
		t.Error("Render should produce non-blank cells")
	}
}

func TestWorldHasTimePeriod(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	// Default should be day.
	if w.TimePeriod() != scene.TimeDay {
		t.Errorf("default time = %v, want day", w.TimePeriod())
	}
}

func TestWorldChangeTime(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	original := w.TimePeriod()
	w.ChangeTime()
	if w.TimePeriod() == original {
		t.Error("ChangeTime should switch to a different time period")
	}
}

func TestWorldChangeTimeCyclesThroughAll(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	seen := map[scene.TimePeriod]bool{}
	for i := 0; i < 4; i++ {
		seen[w.TimePeriod()] = true
		w.ChangeTime()
	}
	if len(seen) != 4 {
		t.Errorf("should have seen 4 time periods, saw %d", len(seen))
	}
}

func TestWorldHasWeather(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	// Default should be clear.
	if w.Weather() != 0 {
		t.Errorf("default weather = %v, want clear", w.Weather())
	}
}

func TestWorldChangeWeather(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	w.ChangeWeather()
	// After one change, should not be clear (cycles clear → rain → snow).
	if w.Weather() == 0 {
		t.Error("ChangeWeather should switch away from clear")
	}
}

func TestWorldChangeWeatherCyclesThroughAll(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	seen := map[interface{}]bool{}
	for i := 0; i < 3; i++ {
		seen[w.Weather()] = true
		w.ChangeWeather()
	}
	if len(seen) != 3 {
		t.Errorf("should have seen 3 weather types, saw %d", len(seen))
	}
}

func TestWorldRenderNightProducesStars(t *testing.T) {
	w := NewWorld(120, 40, 42, tcell.StyleDefault)
	// Cycle time to night.
	for w.TimePeriod() != scene.TimeNight {
		w.ChangeTime()
	}
	canvas := render.NewCanvas(120, 40)
	w.Render(canvas)
	// At night, stars should be present in the sky.
	hasStar := false
	for y := 0; y < w.TrackY() && !hasStar; y++ {
		for x := 0; x < 120; x++ {
			if canvas.CellAt(x, y).Rune == '*' {
				hasStar = true
				break
			}
		}
	}
	if !hasStar {
		t.Error("night render should have stars in the sky")
	}
}

func TestWorldRenderRainProducesParticles(t *testing.T) {
	w := NewWorld(120, 40, 42, tcell.StyleDefault)
	// Cycle weather to rain.
	for w.Weather() == 0 { // WeatherClear == 0
		w.ChangeWeather()
	}
	canvas := render.NewCanvas(120, 40)
	w.Render(canvas)
	// Rain particles ('|') should be present.
	hasRain := false
	for y := 0; y < 40 && !hasRain; y++ {
		for x := 0; x < 120; x++ {
			if canvas.CellAt(x, y).Rune == '|' {
				hasRain = true
				break
			}
		}
	}
	if !hasRain {
		t.Error("rain render should have '|' particles")
	}
}

func TestWorldHelpToggle(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	if w.HelpVisible() {
		t.Error("help should not be visible initially")
	}
	w.ToggleHelp()
	if !w.HelpVisible() {
		t.Error("help should be visible after ToggleHelp")
	}
	w.ToggleHelp()
	if w.HelpVisible() {
		t.Error("help should be hidden after second ToggleHelp")
	}
}

func TestWorldSetHelpVisible(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	w.SetHelpVisible(true)
	if !w.HelpVisible() {
		t.Error("help should be visible after SetHelpVisible(true)")
	}
	w.SetHelpVisible(false)
	if w.HelpVisible() {
		t.Error("help should be hidden after SetHelpVisible(false)")
	}
}

func TestWorldRenderHelpOverlay(t *testing.T) {
	w := NewWorld(120, 40, 42, tcell.StyleDefault)
	w.SetHelpVisible(true)
	canvas := render.NewCanvas(120, 40)
	w.Render(canvas)
	// The help overlay should contain "Conductor" text.
	found := false
	for y := 0; y < 40 && !found; y++ {
		for x := 0; x < 120-8; x++ {
			text := ""
			for i := 0; i < 9; i++ {
				cell := canvas.CellAt(x+i, y)
				if cell.IsBlank() {
					text += " "
				} else {
					text += string(cell.Rune)
				}
			}
			if text == "Conductor" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("help overlay should contain 'Conductor' text")
	}
}

func TestWorldNoColor(t *testing.T) {
	w := NewWorld(80, 24, 42, tcell.StyleDefault)
	if w.NoColor() {
		t.Error("NoColor should be false by default")
	}
	w.ChangeColorPalette()
	if !w.NoColor() {
		t.Error("NoColor should be true after ChangeColorPalette")
	}
	w.ChangeColorPalette()
	if w.NoColor() {
		t.Error("NoColor should be false after second ChangeColorPalette")
	}
}

func TestWorldNewWorldWithConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 99
	cfg.FPS = 15
	cfg.Speed = 1.5
	cfg.Train = "steam"
	cfg.Scene = "station"
	cfg.Weather = "rain"
	cfg.Time = "night"
	cfg.NoColor = true
	cfg.ReducedMotion = true

	w := NewWorldWithConfig(120, 40, cfg, tcell.StyleDefault)
	if w.NoColor() != true {
		t.Error("NoColor should be true from config")
	}
	if w.ReducedMotion() != true {
		t.Error("ReducedMotion should be true from config")
	}
	if w.TimePeriod() != scene.TimeNight {
		t.Errorf("TimePeriod = %v, want night", w.TimePeriod())
	}
	// Station scene.
	if w.Scene().Type() != scene.SceneStation {
		t.Errorf("Scene type = %v, want station", w.Scene().Type())
	}
	// Rain weather.
	if w.Weather() == 0 { // WeatherClear == 0
		t.Error("Weather should not be clear")
	}
}

func TestWorldNewWorldWithConfigRandom(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	cfg.Scene = "random"
	cfg.Weather = "random"
	cfg.Time = "auto"

	w := NewWorldWithConfig(120, 40, cfg, tcell.StyleDefault)
	// Should not crash and should have a scene, weather, and time.
	if w.Scene() == nil {
		t.Fatal("random scene should produce a scene")
	}
	_ = w.Weather()
	_ = w.TimePeriod()
}
