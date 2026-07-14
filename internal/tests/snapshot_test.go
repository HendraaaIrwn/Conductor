// Package tests contains snapshot and integration tests for Conductor.
// Snapshot tests generate deterministic rendered frames without a real
// terminal and compare them against golden files.
package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/conductor/internal/config"
	"github.com/example/conductor/internal/engine"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// snapshotDir is the directory where golden files are stored, relative to
// the project root. Tests run from the package directory, so we go up two
// levels from internal/tests/ to reach the project root.
var snapshotDir = filepath.Join("..", "..", "testdata", "snapshots")

// renderFrame renders a single frame from the world into a canvas string.
// This is deterministic given the config, viewport, and delta.
func renderFrame(cfg config.Config, width, height int, deltaSeconds float64) string {
	style := tcell.StyleDefault
	w := engine.NewWorldWithConfig(width, height, cfg, style)
	canvas := render.NewCanvas(width, height)
	w.Update(deltaSeconds)
	w.Render(canvas)
	return canvasToString(canvas)
}

// canvasToString converts a canvas to a portable string representation.
// Trailing blanks on each line are trimmed for stable snapshots.
func canvasToString(c *render.Canvas) string {
	var b strings.Builder
	for y := 0; y < c.Height(); y++ {
		line := ""
		for x := 0; x < c.Width(); x++ {
			cell := c.CellAt(x, y)
			if cell.IsBlank() {
				line += " "
			} else {
				line += string(cell.Rune)
			}
		}
		line = strings.TrimRight(line, " ")
		if line != "" || y < c.Height()-1 {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// snapshotPath returns the golden file path for the given test name.
func snapshotPath(name string) string {
	return filepath.Join(snapshotDir, name+".txt")
}

// updateFlag is set with -update to regenerate golden files.
var updateFlag = os.Getenv("UPDATE_SNAPSHOTS") == "1"

// verifySnapshot compares the rendered output against the golden file, or
// writes the golden file if UPDATE_SNAPSHOTS=1 is set.
func verifySnapshot(t *testing.T, name string, got string) {
	t.Helper()
	shouldUpdate := os.Getenv("UPDATE_SNAPSHOTS") == "1"
	path := snapshotPath(name)
	if shouldUpdate {
		os.MkdirAll(filepath.Dir(path), 0755)
		if err := os.WriteFile(path, []byte(got), 0644); err != nil {
			t.Fatalf("writing snapshot %s: %v", name, err)
		}
		t.Logf("snapshot %s written", name)
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading snapshot %s: %v. To generate, run: UPDATE_SNAPSHOTS=1 go test ./internal/tests/...", name, err)
	}
	if got != string(want) {
		t.Errorf("snapshot %s mismatch", name)
		// Write the actual output for diffing.
		actualPath := filepath.Join(snapshotDir, name+"_actual.txt")
		os.WriteFile(actualPath, []byte(got), 0644)
		t.Logf("actual output written to %s", actualPath)
	}
}

func TestSnapshotCountrysideDay(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "countryside"
	cfg.Time = "day"
	cfg.Weather = "clear"
	cfg.Seed = 42
	cfg.FPS = 20
	cfg.Speed = 1.0

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "countryside_day", got)
}

func TestSnapshotCountrysideNight(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "countryside"
	cfg.Time = "night"
	cfg.Weather = "clear"
	cfg.Seed = 42

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "countryside_night", got)
}

func TestSnapshotStationDay(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "station"
	cfg.Time = "day"
	cfg.Weather = "clear"
	cfg.Seed = 42

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "station_day", got)
}

func TestSnapshotStationNight(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "station"
	cfg.Time = "night"
	cfg.Weather = "clear"
	cfg.Seed = 42

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "station_night", got)
}

func TestSnapshotMountainDay(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "mountain"
	cfg.Time = "day"
	cfg.Weather = "clear"
	cfg.Seed = 42

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "mountain_day", got)
}

func TestSnapshotMountainRain(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "mountain"
	cfg.Time = "day"
	cfg.Weather = "rain"
	cfg.Seed = 42

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "mountain_rain", got)
}

func TestSnapshotNoColor(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "countryside"
	cfg.Time = "day"
	cfg.Weather = "clear"
	cfg.Seed = 42
	cfg.NoColor = true

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "no_color", got)
}

func TestSnapshotReducedMotion(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "station"
	cfg.Time = "night"
	cfg.Weather = "rain"
	cfg.Seed = 42
	cfg.ReducedMotion = true

	got := renderFrame(cfg, 100, 30, 2.0)
	verifySnapshot(t, "reduced_motion", got)
}

func TestSnapshotDeterministic(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Scene = "countryside"
	cfg.Time = "day"
	cfg.Weather = "clear"
	cfg.Seed = 42

	frame1 := renderFrame(cfg, 100, 30, 1.0)
	frame2 := renderFrame(cfg, 100, 30, 1.0)
	if frame1 != frame2 {
		t.Error("same seed should produce identical frames")
	}
}

func TestSnapshotDifferentSeedsDiffer(t *testing.T) {
	cfg1 := config.DefaultConfig()
	cfg1.Scene = "countryside"
	cfg1.Time = "day"
	cfg1.Weather = "clear"
	cfg1.Seed = 1

	cfg2 := config.DefaultConfig()
	cfg2.Scene = "countryside"
	cfg2.Time = "day"
	cfg2.Weather = "clear"
	cfg2.Seed = 2

	frame1 := renderFrame(cfg1, 100, 30, 1.0)
	frame2 := renderFrame(cfg2, 100, 30, 1.0)
	if frame1 == frame2 {
		t.Error("different seeds should produce different frames")
	}
}
