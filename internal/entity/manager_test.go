package entity

import (
	"testing"
	"time"

	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

func TestManagerAddAndCount(t *testing.T) {
	m := NewManager()
	e1 := &Entity{Type: TypeSmoke, X: 0, Y: 0, Visible: true}
	e2 := &Entity{Type: TypeSmoke, X: 10, Y: 0, Visible: true}
	m.Add(e1)
	m.Add(e2)
	// Entities are in pendingAdd; flush happens during Update.
	m.Update(80, 24, 1*time.Millisecond)
	if m.Count() != 2 {
		t.Errorf("Count = %d, want 2", m.Count())
	}
}

func TestManagerAssignsIDs(t *testing.T) {
	m := NewManager()
	e := &Entity{Type: TypeSmoke, Visible: true}
	m.Add(e)
	m.Update(80, 24, 1*time.Millisecond)
	if e.ID == "" {
		t.Error("entity should have been assigned an ID")
	}
}

func TestManagerByType(t *testing.T) {
	m := NewManager()
	m.Add(&Entity{Type: TypeSmoke, Visible: true})
	m.Add(&Entity{Type: TypeTrain, Visible: true})
	m.Add(&Entity{Type: TypeSmoke, Visible: true})
	m.Update(80, 24, 1*time.Millisecond)
	smoke := m.ByType(TypeSmoke)
	if len(smoke) != 2 {
		t.Errorf("smoke count = %d, want 2", len(smoke))
	}
	trains := m.ByType(TypeTrain)
	if len(trains) != 1 {
		t.Errorf("train count = %d, want 1", len(trains))
	}
}

func TestManagerFirstByType(t *testing.T) {
	m := NewManager()
	m.Add(&Entity{Type: TypeTrain, Visible: true})
	m.Update(80, 24, 1*time.Millisecond)
	if m.FirstByType(TypeTrain) == nil {
		t.Error("expected a train entity")
	}
	if m.FirstByType(TypeSmoke) != nil {
		t.Error("expected nil for smoke")
	}
}

func TestManagerRemove(t *testing.T) {
	m := NewManager()
	e := &Entity{Type: TypeSmoke, Visible: true}
	m.Add(e)
	m.Update(80, 24, 1*time.Millisecond)
	m.Remove(e.ID)
	m.Update(80, 24, 1*time.Millisecond)
	if m.Count() != 0 {
		t.Errorf("after remove Count = %d, want 0", m.Count())
	}
}

func TestManagerLifetimeExpiration(t *testing.T) {
	m := NewManager()
	e := &Entity{
		Type:     TypeSmoke,
		Visible:  true,
		Lifetime: 100 * time.Millisecond,
	}
	m.Add(e)
	m.Update(80, 24, 50*time.Millisecond)
	if e.Dead {
		t.Error("entity should not be dead before lifetime expires")
	}
	m.Update(80, 24, 60*time.Millisecond)
	if !e.Dead {
		t.Error("entity should be dead after lifetime expires")
	}
	m.Update(80, 24, 1*time.Millisecond) // flush
	if m.Count() != 0 {
		t.Errorf("expired entity should be removed, Count = %d", m.Count())
	}
}

func TestManagerOffscreenRemoval(t *testing.T) {
	m := NewManager()
	e := &Entity{
		Type:            TypeSmoke,
		X:               -10,
		Y:               0,
		Visible:         true,
		RemoveOffscreen: true,
	}
	m.Add(e)
	m.Update(80, 24, 1*time.Millisecond)
	m.Update(80, 24, 1*time.Millisecond) // flush
	if m.Count() != 0 {
		t.Errorf("offscreen entity should be removed, Count = %d", m.Count())
	}
}

func TestManagerKeepsNonOffscreenEntity(t *testing.T) {
	m := NewManager()
	e := &Entity{
		Type:            TypeSmoke,
		X:               40,
		Y:               12,
		Visible:         true,
		RemoveOffscreen: true,
	}
	m.Add(e)
	m.Update(80, 24, 1*time.Millisecond)
	m.Update(80, 24, 1*time.Millisecond)
	if m.Count() != 1 {
		t.Errorf("on-screen entity should be kept, Count = %d", m.Count())
	}
}

func TestManagerRenderByLayer(t *testing.T) {
	m := NewManager()
	// Add entities on different layers. We can't directly verify draw order
	// from the manager, but we can verify sortedByLayer returns them in
	// ascending layer order.
	m.Add(&Entity{Type: TypeTrain, Layer: 60, Visible: true})
	m.Add(&Entity{Type: TypeScenery, Layer: 20, Visible: true})
	m.Add(&Entity{Type: TypeSmoke, Layer: 80, Visible: true})
	m.Update(80, 24, 1*time.Millisecond)
	sorted := m.sortedByLayer()
	if sorted[0].Layer != 20 {
		t.Errorf("first layer = %d, want 20", sorted[0].Layer)
	}
	if sorted[1].Layer != 60 {
		t.Errorf("second layer = %d, want 60", sorted[1].Layer)
	}
	if sorted[2].Layer != 80 {
		t.Errorf("third layer = %d, want 80", sorted[2].Layer)
	}
}

func TestManagerRenderUsesSprite(t *testing.T) {
	m := NewManager()
	sprite := render.NewSprite("test",
		render.ParseFrame("X", tcell.StyleDefault))
	m.Add(&Entity{
		Type: TypeSmoke,
		X:    5, Y: 5,
		Layer:   80,
		Sprite:  sprite,
		Visible: true,
	})
	m.Update(80, 24, 1*time.Millisecond)
	canvas := render.NewCanvas(80, 24)
	m.Render(canvas)
	// Cell at (5,5) should have the 'X' rune.
	cell := canvas.CellAt(5, 5)
	if cell.Rune != 'X' {
		t.Errorf("cell (5,5) rune = %q, want 'X'", cell.Rune)
	}
}

func TestManagerRenderUsesRenderFunc(t *testing.T) {
	m := NewManager()
	called := false
	m.Add(&Entity{
		Type: TypeTrain,
		X:    0, Y: 0,
		Layer:   60,
		Visible: true,
		RenderFunc: func(c *render.Canvas, e *Entity) {
			called = true
		},
	})
	m.Update(80, 24, 1*time.Millisecond)
	canvas := render.NewCanvas(80, 24)
	m.Render(canvas)
	if !called {
		t.Error("RenderFunc was not called")
	}
}

func TestManagerPauseSkipsUpdate(t *testing.T) {
	m := NewManager()
	// When delta is 0 (paused), no update should occur.
	m.Add(&Entity{
		Type:     TypeSmoke,
		X:        0,
		Y:        0,
		VX:       10,
		Visible:  true,
		Behavior: LinearMovement{},
	})
	m.Update(80, 24, 0) // zero delta = paused
	if m.FirstByType(TypeSmoke).X != 0 {
		t.Error("entity should not move with zero delta")
	}
}

func TestLinearMovement(t *testing.T) {
	e := &Entity{X: 0, Y: 0, VX: 10, VY: 5}
	LinearMovement{}.Update(e, UpdateContext{}, 1*time.Second)
	if e.X != 10 {
		t.Errorf("X = %f, want 10", e.X)
	}
	if e.Y != 5 {
		t.Errorf("Y = %f, want 5", e.Y)
	}
}

func TestAnimation(t *testing.T) {
	sprite := render.NewSprite("test",
		render.ParseFrame("A", tcell.StyleDefault),
		render.ParseFrame("B", tcell.StyleDefault),
		render.ParseFrame("C", tcell.StyleDefault),
	)
	e := &Entity{Sprite: sprite, Frame: 0}
	anim := &Animation{FrameTime: 100 * time.Millisecond}
	anim.Update(e, UpdateContext{}, 100*time.Millisecond)
	if e.Frame != 1 {
		t.Errorf("Frame = %d, want 1", e.Frame)
	}
	anim.Update(e, UpdateContext{}, 100*time.Millisecond)
	if e.Frame != 2 {
		t.Errorf("Frame = %d, want 2", e.Frame)
	}
	anim.Update(e, UpdateContext{}, 100*time.Millisecond)
	if e.Frame != 0 {
		t.Errorf("Frame = %d, want 0 (wrapped)", e.Frame)
	}
}

func TestCompositeBehavior(t *testing.T) {
	e := &Entity{X: 0, Y: 0, VX: 10, VY: 0}
	sprite := render.NewSprite("test",
		render.ParseFrame("A", tcell.StyleDefault),
		render.ParseFrame("B", tcell.StyleDefault),
	)
	e.Sprite = sprite
	anim := &Animation{FrameTime: 100 * time.Millisecond}
	c := Composite{Behaviors: []Behavior{LinearMovement{}, anim}}
	c.Update(e, UpdateContext{}, 100*time.Millisecond)
	if e.X != 1 {
		t.Errorf("X = %f, want 1", e.X)
	}
	if e.Frame != 1 {
		t.Errorf("Frame = %d, want 1", e.Frame)
	}
}
