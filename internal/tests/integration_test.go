package tests

import (
	"testing"
	"time"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/engine"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// TestIntegrationStartupShutdown creates a world, verifies it initializes
// correctly, and runs a few frames without crashing.
func TestIntegrationStartupShutdown(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)
	if w == nil {
		t.Fatal("world should not be nil")
	}
	if w.Width() != 120 || w.Height() != 40 {
		t.Errorf("viewport = %dx%d, want 120x40", w.Width(), w.Height())
	}
	// Run a few update cycles.
	for i := 0; i < 10; i++ {
		w.Update(0.05)
	}
}

// TestIntegrationPauseResume verifies that pause stops movement and resume
// continues it.
func TestIntegrationPauseResume(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 1
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	// Get the initial position of the train.
	trains := w.Trains()
	if len(trains) == 0 {
		t.Skip("no train in world")
	}
	initialX := trains[0].X

	// Pause and update.
	w.SetPaused(true)
	w.Update(1.0)
	if trains[0].X != initialX {
		t.Error("train should not move while paused")
	}

	// Resume and update.
	w.SetPaused(false)
	w.Update(0.5)
	if trains[0].X == initialX {
		t.Error("train should move after resume")
	}
}

// TestIntegrationResize verifies that the world handles resize events
// without crashing.
func TestIntegrationResize(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	// Resize to a larger viewport.
	w.HandleResize(160, 50)
	if w.Width() != 160 || w.Height() != 50 {
		t.Errorf("after resize = %dx%d, want 160x50", w.Width(), w.Height())
	}
	// Update after resize.
	w.Update(0.1)

	// Resize back to a smaller viewport.
	w.HandleResize(80, 24)
	if w.Width() != 80 || w.Height() != 24 {
		t.Errorf("after resize = %dx%d, want 80x24", w.Width(), w.Height())
	}
	w.Update(0.1)
}

// TestIntegrationSceneChange verifies that scene switching works correctly.
func TestIntegrationSceneChange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	initialType := w.Scene().Type()
	w.ChangeScene()
	if w.Scene().Type() == initialType {
		t.Error("scene type should change")
	}
	// Update after scene change.
	w.Update(0.1)
}

// TestIntegrationWeatherChange verifies that weather switching works.
func TestIntegrationWeatherChange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	initialWeather := w.Weather()
	w.ChangeWeather()
	if w.Weather() == initialWeather {
		t.Error("weather should change")
	}
	// Update after weather change.
	w.Update(0.1)
}

// TestIntegrationTimeChange verifies that time switching works.
func TestIntegrationTimeChange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	initialTime := w.TimePeriod()
	w.ChangeTime()
	if w.TimePeriod() == initialTime {
		t.Error("time should change")
	}
	// Update after time change.
	w.Update(0.1)
}

// TestIntegrationRegenerate verifies that regenerating the world works.
func TestIntegrationRegenerate(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)
	w.Regenerate()
	// Update after regenerate.
	w.Update(0.1)
}

// TestIntegrationSpeedChange verifies that speed changes are applied.
func TestIntegrationSpeedChange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	w.SetSpeed(50.0)
	if w.Speed() != 50.0 {
		t.Errorf("speed = %f, want 50.0", w.Speed())
	}
	// Update after speed change.
	w.Update(0.1)
}

// TestIntegrationRenderProducesOutput verifies that rendering produces
// non-blank cells.
func TestIntegrationRenderProducesOutput(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	canvas := render.NewCanvas(120, 40)
	w.Update(0.1)
	w.Render(canvas)

	nonBlank := 0
	for y := 0; y < 40; y++ {
		for x := 0; x < 120; x++ {
			if !canvas.CellAt(x, y).IsBlank() {
				nonBlank++
			}
		}
	}
	if nonBlank == 0 {
		t.Error("render should produce non-blank cells")
	}
}

// TestIntegrationLongRunning verifies that the world can run for an extended
// number of frames without crashing.
func TestIntegrationLongRunning(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	// Simulate 60 seconds of runtime at 20 FPS = 1200 frames.
	for i := 0; i < 1200; i++ {
		w.Update(0.05)
	}
}

// TestIntegrationNoColor verifies no-color mode rendering.
func TestIntegrationNoColor(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	cfg.NoColor = true
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	if !w.NoColor() {
		t.Error("NoColor should be true")
	}
	canvas := render.NewCanvas(120, 40)
	w.Update(0.1)
	w.Render(canvas)
}

// TestIntegrationReducedMotion verifies reduced-motion mode.
func TestIntegrationReducedMotion(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	cfg.ReducedMotion = true
	cfg.Weather = "rain"
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	if !w.ReducedMotion() {
		t.Error("ReducedMotion should be true")
	}
	canvas := render.NewCanvas(120, 40)
	w.Update(0.1)
	w.Render(canvas)
}

// TestIntegrationMinTerminalSize checks that the minimum terminal size
// calculation works.
func TestIntegrationMinTerminalSize(t *testing.T) {
	// Small viewport should still work (the size check is in the loop,
	// not the world constructor).
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(50, 15, cfg, style)
	if w == nil {
		t.Fatal("world should not be nil for small viewport")
	}
	canvas := render.NewCanvas(50, 15)
	w.Update(0.1)
	w.Render(canvas)
}

// TestIntegrationNoCrashOnZeroSize checks that zero-sized viewport doesn't
// cause a panic.
func TestIntegrationNoCrashOnZeroSize(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(0, 0, cfg, style)
	if w == nil {
		t.Fatal("world should not be nil for zero viewport")
	}
	canvas := render.NewCanvas(0, 0)
	w.Update(0.1)
	w.Render(canvas)
}

// TestIntegrationSpawnNext verifies that spawning a new train works.
func TestIntegrationSpawnNext(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)
	w.SpawnNext()
	w.Update(0.1)
	w.SpawnNext()
	w.Update(0.1)
}

// TestIntegrationColorToggle verifies color palette toggling.
func TestIntegrationColorToggle(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Seed = 42
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(120, 40, cfg, style)

	w.ChangeColorPalette()
	if !w.NoColor() {
		t.Error("NoColor should be true after toggle")
	}
	w.ChangeColorPalette()
	if w.NoColor() {
		t.Error("NoColor should be false after second toggle")
	}
}

// _ is used to prevent unused import warnings.
var _ = time.Second
