// Package entity provides the entity system used by the animation engine.
// Every visual object in the simulation — trains, smoke particles, signals,
// clouds — is represented as an Entity. The Manager owns all entities for a
// scene and is responsible for updating, culling, and rendering them in layer
// order.
package entity

import (
	"time"

	"github.com/example/conductor/internal/render"
)

// Type classifies an entity for querying and dispatching.
type Type string

const (
	TypeTrain    Type = "train"
	TypeSmoke    Type = "smoke"
	TypeRain     Type = "rain"
	TypeSnow     Type = "snow"
	TypeSignal   Type = "signal"
	TypeCrossing Type = "crossing"
	TypeBird     Type = "bird"
	TypeCloud    Type = "cloud"
	TypeScenery  Type = "scenery"
	TypeTrack    Type = "track"
	TypeOverlay  Type = "overlay"
)

// RenderFunc is a custom rendering function for entities that need more than
// simple sprite blitting (e.g. trains with multiple cars). If an entity's
// RenderFunc is nil, the manager falls back to Sprite-based rendering.
type RenderFunc func(canvas *render.Canvas, e *Entity)

// Entity is a single visual object in the simulation. Entities are updated by
// their Behavior, culled by the Manager (lifetime / off-screen), and rendered
// in layer order.
type Entity struct {
	ID     string
	Type   Type
	X, Y   float64 // logical position in cells (float for smooth motion)
	VX, VY float64 // velocity in cells per second
	Layer  int     // render layer (see render.Layer* constants)

	Sprite *render.Sprite // used when RenderFunc is nil
	Frame  int            // current sprite frame index

	Visible bool

	Lifetime time.Duration // 0 means infinite
	Age      time.Duration

	Behavior   Behavior
	RenderFunc RenderFunc

	// Data holds type-specific state (e.g. *TrainData for trains). The
	// manager never inspects this field; it is owned by the Behavior and
	// RenderFunc.
	Data any

	// RemoveOffscreen controls whether the manager automatically removes
	// the entity when it leaves the viewport. Set to false for entities
	// that recycle themselves (e.g. rain particles).
	RemoveOffscreen bool

	// Dead is set by behaviors or the manager to signal that the entity
	// should be removed during the next cleanup pass.
	Dead bool
}

// Behavior updates an entity's state each frame. Behaviors must not depend on
// terminal rendering — they operate purely on simulation state.
type Behavior interface {
	Update(e *Entity, ctx UpdateContext, delta time.Duration)
}

// UpdateContext is passed to Behavior.Update. It provides viewport dimensions
// and a way to spawn new entities without exposing the full Manager.
type UpdateContext struct {
	Width  int
	Height int
	Spawn  func(*Entity)
}
