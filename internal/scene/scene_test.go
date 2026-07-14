package scene

import (
	"math/rand"
	"testing"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

func testStyle() tcell.Style {
	return tcell.StyleDefault
}

func testViewport() Viewport {
	return Viewport{Width: 120, Height: 40}
}

func TestParseSceneType(t *testing.T) {
	cases := []struct {
		input string
		want  SceneType
		ok    bool
	}{
		{"countryside", SceneCountryside, true},
		{"station", SceneStation, true},
		{"mountain", SceneMountain, true},
		{"unknown", 0, false},
		{"", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseSceneType(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseSceneType(%q) ok = %v, want %v", tc.input, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Errorf("ParseSceneType(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestSceneTypeString(t *testing.T) {
	if SceneCountryside.String() != "countryside" {
		t.Errorf("countryside string = %q", SceneCountryside.String())
	}
	if SceneStation.String() != "station" {
		t.Errorf("station string = %q", SceneStation.String())
	}
	if SceneMountain.String() != "mountain" {
		t.Errorf("mountain string = %q", SceneMountain.String())
	}
}

func TestAllSceneTypes(t *testing.T) {
	types := AllSceneTypes()
	if len(types) != 3 {
		t.Errorf("AllSceneTypes length = %d, want 3", len(types))
	}
}

func TestNewScene(t *testing.T) {
	for _, st := range AllSceneTypes() {
		s := New(st)
		if s == nil {
			t.Errorf("New(%v) returned nil", st)
			continue
		}
		if s.Type() != st {
			t.Errorf("New(%v).Type() = %v", st, s.Type())
		}
	}
}

func TestNewRandomScene(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	s := NewRandom(rng)
	if s == nil {
		t.Fatal("NewRandom returned nil")
	}
}

func TestViewportTrackY(t *testing.T) {
	vp := Viewport{Width: 120, Height: 40}
	// 40 * 7 / 10 = 28
	if vp.TrackY() != 28 {
		t.Errorf("TrackY = %d, want 28", vp.TrackY())
	}
}

func TestViewportTrackYMinClamp(t *testing.T) {
	vp := Viewport{Width: 80, Height: 3}
	// 3 * 7 / 10 = 2, but min is 3
	if vp.TrackY() != 3 {
		t.Errorf("TrackY = %d, want 3 (clamped)", vp.TrackY())
	}
}

func TestCountrysideBuild(t *testing.T) {
	s := &Countryside{}
	mgr := entity.NewManager()
	rng := rand.New(rand.NewSource(42))
	vp := testViewport()
	s.Build(mgr, vp, rng, testStyle())
	mgr.Flush()
	// Clouds should have been spawned.
	clouds := mgr.ByType(entity.TypeCloud)
	if len(clouds) == 0 {
		t.Error("countryside should spawn clouds")
	}
}

func TestCountrysideName(t *testing.T) {
	s := &Countryside{}
	if s.Name() != "countryside" {
		t.Errorf("Name = %q, want 'countryside'", s.Name())
	}
	if s.Type() != SceneCountryside {
		t.Errorf("Type = %v, want SceneCountryside", s.Type())
	}
}

func TestCountrysideRenderBackgroundProducesOutput(t *testing.T) {
	s := &Countryside{}
	mgr := entity.NewManager()
	rng := rand.New(rand.NewSource(42))
	vp := testViewport()
	s.Build(mgr, vp, rng, testStyle())
	mgr.Flush()
	canvas := render.NewCanvas(vp.Width, vp.Height)
	s.RenderBackground(canvas, vp, testStyle())
	// At least some cells should be non-blank.
	nonBlank := 0
	for y := 0; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			if !canvas.CellAt(x, y).IsBlank() {
				nonBlank++
			}
		}
	}
	if nonBlank == 0 {
		t.Error("countryside background should produce non-blank cells")
	}
}

func TestCountrysideRenderTrackProducesOutput(t *testing.T) {
	s := &Countryside{}
	vp := testViewport()
	canvas := render.NewCanvas(vp.Width, vp.Height)
	s.RenderTrack(canvas, vp, testStyle())
	trackY := vp.TrackY()
	// The track row should have '=' characters.
	hasRail := false
	for x := 0; x < vp.Width; x++ {
		if canvas.CellAt(x, trackY).Rune == '=' {
			hasRail = true
			break
		}
	}
	if !hasRail {
		t.Error("countryside track should have rail characters")
	}
}

func TestStationBuild(t *testing.T) {
	s := &Station{}
	mgr := entity.NewManager()
	rng := rand.New(rand.NewSource(42))
	vp := testViewport()
	s.Build(mgr, vp, rng, testStyle())
	mgr.Flush()
	// Signals should have been spawned.
	signals := mgr.ByType(entity.TypeSignal)
	if len(signals) < 2 {
		t.Errorf("station should spawn at least 2 signals, got %d", len(signals))
	}
	// Clouds should have been spawned.
	clouds := mgr.ByType(entity.TypeCloud)
	if len(clouds) == 0 {
		t.Error("station should spawn clouds")
	}
}

func TestStationName(t *testing.T) {
	s := &Station{}
	if s.Name() != "station" {
		t.Errorf("Name = %q, want 'station'", s.Name())
	}
	if s.Type() != SceneStation {
		t.Errorf("Type = %v, want SceneStation", s.Type())
	}
}

func TestStationRenderBackgroundProducesOutput(t *testing.T) {
	s := &Station{}
	mgr := entity.NewManager()
	rng := rand.New(rand.NewSource(42))
	vp := testViewport()
	s.Build(mgr, vp, rng, testStyle())
	mgr.Flush()
	canvas := render.NewCanvas(vp.Width, vp.Height)
	s.RenderBackground(canvas, vp, testStyle())
	nonBlank := 0
	for y := 0; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			if !canvas.CellAt(x, y).IsBlank() {
				nonBlank++
			}
		}
	}
	if nonBlank == 0 {
		t.Error("station background should produce non-blank cells")
	}
}

func TestMountainBuild(t *testing.T) {
	s := &Mountain{}
	mgr := entity.NewManager()
	rng := rand.New(rand.NewSource(42))
	vp := testViewport()
	s.Build(mgr, vp, rng, testStyle())
	mgr.Flush()
	// Clouds should have been spawned.
	clouds := mgr.ByType(entity.TypeCloud)
	if len(clouds) == 0 {
		t.Error("mountain should spawn clouds")
	}
}

func TestMountainName(t *testing.T) {
	s := &Mountain{}
	if s.Name() != "mountain" {
		t.Errorf("Name = %q, want 'mountain'", s.Name())
	}
	if s.Type() != SceneMountain {
		t.Errorf("Type = %v, want SceneMountain", s.Type())
	}
}

func TestMountainRenderBackgroundProducesOutput(t *testing.T) {
	s := &Mountain{}
	mgr := entity.NewManager()
	rng := rand.New(rand.NewSource(42))
	vp := testViewport()
	s.Build(mgr, vp, rng, testStyle())
	mgr.Flush()
	canvas := render.NewCanvas(vp.Width, vp.Height)
	s.RenderBackground(canvas, vp, testStyle())
	nonBlank := 0
	for y := 0; y < vp.Height; y++ {
		for x := 0; x < vp.Width; x++ {
			if !canvas.CellAt(x, y).IsBlank() {
				nonBlank++
			}
		}
	}
	if nonBlank == 0 {
		t.Error("mountain background should produce non-blank cells")
	}
}

func TestSceneHandleResize(t *testing.T) {
	for _, st := range AllSceneTypes() {
		s := New(st)
		mgr := entity.NewManager()
		rng := rand.New(rand.NewSource(42))
		vp := testViewport()
		s.Build(mgr, vp, rng, testStyle())
		mgr.Flush()
		// Resize to a different viewport.
		newVP := Viewport{Width: 100, Height: 30}
		s.HandleResize(mgr, newVP, rand.New(rand.NewSource(42)), testStyle())
		mgr.Flush()
		// Should not crash and should still have entities.
		if mgr.Count() == 0 {
			t.Errorf("%v: no entities after resize", st)
		}
	}
}

func TestSceneDeterministicBuild(t *testing.T) {
	for _, st := range AllSceneTypes() {
		s1 := New(st)
		s2 := New(st)
		mgr1 := entity.NewManager()
		mgr2 := entity.NewManager()
		rng1 := rand.New(rand.NewSource(99))
		rng2 := rand.New(rand.NewSource(99))
		vp := testViewport()
		s1.Build(mgr1, vp, rng1, testStyle())
		s2.Build(mgr2, vp, rng2, testStyle())
		mgr1.Flush()
		mgr2.Flush()
		// Both scenes should spawn the same number of entities.
		if mgr1.Count() != mgr2.Count() {
			t.Errorf("%v: entity count differs: %d vs %d", st, mgr1.Count(), mgr2.Count())
		}
		// Render backgrounds should be identical.
		c1 := render.NewCanvas(vp.Width, vp.Height)
		c2 := render.NewCanvas(vp.Width, vp.Height)
		s1.RenderBackground(c1, vp, testStyle())
		s2.RenderBackground(c2, vp, testStyle())
		for y := 0; y < vp.Height; y++ {
			for x := 0; x < vp.Width; x++ {
				if c1.CellAt(x, y) != c2.CellAt(x, y) {
					t.Errorf("%v: background differs at (%d,%d)", st, x, y)
					break
				}
			}
		}
	}
}

func TestSceneUpdateDoesNotCrash(t *testing.T) {
	for _, st := range AllSceneTypes() {
		s := New(st)
		mgr := entity.NewManager()
		rng := rand.New(rand.NewSource(42))
		vp := testViewport()
		s.Build(mgr, vp, rng, testStyle())
		mgr.Flush()
		// Update should not crash.
		s.Update(mgr, vp, 0.016)
	}
}
