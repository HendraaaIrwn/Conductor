// Package scene defines the Scene interface and provides concrete scene
// implementations. A scene is responsible for generating the background
// scenery, positioning the track, and populating the world with environment
// entities (signals, platforms, clouds, etc.).
//
// Scenes are deterministic: given the same viewport dimensions and seed, they
// produce identical layouts. This is critical for reproducible snapshots and
// for the 'r' (regenerate) key to work correctly.
package scene

import (
	"math/rand"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// Viewport describes the visible terminal area. It is passed to scenes so they
// can lay out content responsively.
type Viewport struct {
	Width  int
	Height int
}

// TrackY returns the Y coordinate where the track should be placed for this
// viewport. Scenes use this to align scenery relative to the rails.
func (v Viewport) TrackY() int {
	y := v.Height * 7 / 10
	if y < 3 {
		y = 3
	}
	return y
}

// SceneType identifies a scene for CLI selection and the 's' key cycle.
type SceneType int

const (
	// SceneCountryside is a rural scene with hills, trees, and houses.
	SceneCountryside SceneType = iota
	// SceneStation is a small station with a platform and signals.
	SceneStation
	// SceneMountain is a mountain route with tunnels and bridges.
	SceneMountain
)

// String returns the human-readable name of the scene type.
func (s SceneType) String() string {
	switch s {
	case SceneCountryside:
		return "countryside"
	case SceneStation:
		return "station"
	case SceneMountain:
		return "mountain"
	default:
		return "unknown"
	}
}

// ParseSceneType converts a string to a SceneType. Returns false if the
// string does not match a known scene.
func ParseSceneType(s string) (SceneType, bool) {
	switch s {
	case "countryside":
		return SceneCountryside, true
	case "station":
		return SceneStation, true
	case "mountain":
		return SceneMountain, true
	default:
		return 0, false
	}
}

// AllSceneTypes returns all scene types for random selection and cycling.
func AllSceneTypes() []SceneType {
	return []SceneType{SceneCountryside, SceneStation, SceneMountain}
}

// Scene is the interface implemented by all scene types. The world calls
// Build once at creation (and on regenerate/resize), Update each frame, and
// Render each frame to draw the background.
type Scene interface {
	// Name returns the scene name string (e.g. "countryside").
	Name() string

	// Type returns the SceneType enum value.
	Type() SceneType

	// Build generates the scene layout and populates the entity manager with
	// scenery entities. It is called on creation, regeneration, and resize.
	// The RNG ensures deterministic layout for a given seed.
	Build(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style)

	// Update is called each frame to advance scene-specific state (e.g. cloud
	// movement, signal state changes). delta is the frame delta in seconds;
	// it is 0 when paused.
	Update(manager *entity.Manager, vp Viewport, delta float64)

	// RenderBackground draws the static background (sky, distant mountains)
	// directly onto the canvas. This is called before entities are rendered.
	RenderBackground(canvas *render.Canvas, vp Viewport, style tcell.Style)

	// RenderBackgroundWithPalette draws the background using a time-of-day
	// palette for celestial elements and color styling.
	RenderBackgroundWithPalette(canvas *render.Canvas, vp Viewport, pal Palette, style tcell.Style)

	// RenderTrack draws the railway track onto the canvas. This is called
	// after RenderBackground but before entities are rendered.
	RenderTrack(canvas *render.Canvas, vp Viewport, style tcell.Style)

	// HandleResize is called when the terminal is resized. The scene should
	// recalculate positions for responsive elements. Entities that no longer
	// fit should be marked Dead.
	HandleResize(manager *entity.Manager, vp Viewport, rng *rand.Rand, style tcell.Style)
}

// New creates a Scene of the given type.
func New(t SceneType) Scene {
	switch t {
	case SceneCountryside:
		return &Countryside{}
	case SceneStation:
		return &Station{}
	case SceneMountain:
		return &Mountain{}
	default:
		return &Countryside{}
	}
}

// NewRandom selects a random scene type from the RNG and returns a new Scene.
func NewRandom(rng *rand.Rand) Scene {
	types := AllSceneTypes()
	return New(types[rng.Intn(len(types))])
}
