package train

import (
	"testing"
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

func testStyle() tcell.Style {
	return tcell.StyleDefault
}

func TestVelocityForDirection(t *testing.T) {
	if v := VelocityForDirection(LeftToRight); v != 1 {
		t.Errorf("LeftToRight velocity = %f, want 1", v)
	}
	if v := VelocityForDirection(RightToLeft); v != -1 {
		t.Errorf("RightToLeft velocity = %f, want -1", v)
	}
}

func TestNewTrainEntity(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(0, 10, 25.0, LeftToRight, cars, nil)
	if e.X != 0 {
		t.Errorf("X = %f, want 0", e.X)
	}
	if e.Y != 10 {
		t.Errorf("Y = %f, want 10", e.Y)
	}
	if e.VX != 25.0 {
		t.Errorf("VX = %f, want 25.0", e.VX)
	}
	data, ok := e.Data.(*TrainData)
	if !ok {
		t.Fatal("Data is not *TrainData")
	}
	if data.Direction != LeftToRight {
		t.Error("Direction should be LeftToRight")
	}
}

func TestTrainEntityUpdateMovesPosition(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(0, 10, 25.0, LeftToRight, cars, nil)
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	e.Behavior.Update(e, ctx, 1*time.Second)
	if e.X != 25.0 {
		t.Errorf("after 1s X = %f, want 25.0", e.X)
	}
}

func TestTrainEntityUpdateLeftDirection(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(100, 10, 25.0, RightToLeft, cars, nil)
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	e.Behavior.Update(e, ctx, 1*time.Second)
	if e.X != 75.0 {
		t.Errorf("after 1s left X = %f, want 75.0", e.X)
	}
}

func TestTrainEntityUpdateAccumulates(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(0, 10, 25.0, LeftToRight, cars, nil)
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	e.Behavior.Update(e, ctx, 500*time.Millisecond)
	e.Behavior.Update(e, ctx, 500*time.Millisecond)
	if e.X != 25.0 {
		t.Errorf("after 0.5+0.5s X = %f, want 25.0", e.X)
	}
}

func TestTrainIsOffscreenRight(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(50, 10, 25.0, LeftToRight, cars, nil)
	if IsOffscreen(e, 80) {
		t.Error("train at X=50 should not be offscreen on width 80")
	}
	e.X = 100
	if !IsOffscreen(e, 80) {
		t.Error("train at X=100 should be offscreen on width 80")
	}
}

func TestTrainIsOffscreenLeft(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(50, 10, 25.0, RightToLeft, cars, nil)
	if IsOffscreen(e, 80) {
		t.Error("train at X=50 should not be offscreen on width 80")
	}
	e.X = -100
	if !IsOffscreen(e, 80) {
		t.Error("train at X=-100 moving left should be offscreen")
	}
}

func TestTrainWidth(t *testing.T) {
	cars := SimpleConsist(testStyle())
	w := TrainWidth(cars)
	if w <= 0 {
		t.Errorf("Width = %d, want > 0", w)
	}
}

func TestTrainWheelAnimationAdvances(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(0, 10, 25.0, LeftToRight, cars, nil)
	data := e.Data.(*TrainData)
	initial := data.WheelFrame
	ctx := entity.UpdateContext{Width: 80, Height: 24, Spawn: func(*entity.Entity) {}}
	e.Behavior.Update(e, ctx, 200*time.Millisecond) // > 150ms frame time
	if data.WheelFrame <= initial {
		t.Errorf("wheelFrame = %d, want > %d", data.WheelFrame, initial)
	}
}

func TestTrainSmokeSpawns(t *testing.T) {
	cars := SimpleConsist(testStyle())
	spawned := false
	factory := func(x, y float64) *entity.Entity {
		spawned = true
		return &entity.Entity{Type: entity.TypeSmoke, X: x, Y: y}
	}
	e := NewTrainEntity(0, 10, 25.0, LeftToRight, cars, factory)
	ctx := entity.UpdateContext{
		Width: 80, Height: 24,
		Spawn: func(en *entity.Entity) {
			if en.Type != entity.TypeSmoke {
				t.Error("spawned entity should be smoke")
			}
		},
	}
	// 400ms should trigger smoke
	e.Behavior.Update(e, ctx, 400*time.Millisecond)
	if !spawned {
		t.Error("smoke should have been spawned after 400ms")
	}
}

func TestTrainSmokeNilFactoryNoSpawn(t *testing.T) {
	cars := SimpleConsist(testStyle())
	e := NewTrainEntity(0, 10, 25.0, LeftToRight, cars, nil)
	spawnCalled := false
	ctx := entity.UpdateContext{
		Width: 80, Height: 24,
		Spawn: func(*entity.Entity) { spawnCalled = true },
	}
	e.Behavior.Update(e, ctx, 500*time.Millisecond)
	if spawnCalled {
		t.Error("no smoke should be spawned with nil factory")
	}
}

func TestParseFrame(t *testing.T) {
	grid := render.ParseFrame("AB\nCD", testStyle())
	if len(grid) != 2 {
		t.Fatalf("grid height = %d, want 2", len(grid))
	}
	if len(grid[0]) != 2 {
		t.Fatalf("grid width = %d, want 2", len(grid[0]))
	}
	if grid[0][0].Rune != 'A' {
		t.Errorf("cell (0,0) = %q, want 'A'", grid[0][0].Rune)
	}
	if grid[1][1].Rune != 'D' {
		t.Errorf("cell (1,1) = %q, want 'D'", grid[1][1].Rune)
	}
}

func TestParseFramePadsShortLines(t *testing.T) {
	grid := render.ParseFrame("ABC\nA", testStyle())
	if len(grid[1]) != 3 {
		t.Errorf("short row width = %d, want 3 (padded)", len(grid[1]))
	}
	if grid[1][0].Rune != 'A' {
		t.Errorf("cell (1,0) = %q, want 'A'", grid[1][0].Rune)
	}
	if !grid[1][1].IsBlank() {
		t.Errorf("cell (1,1) should be blank padding")
	}
}

func TestSimpleConsist(t *testing.T) {
	cars := SimpleConsist(testStyle())
	if len(cars) != 2 {
		t.Fatalf("consist length = %d, want 2", len(cars))
	}
	if cars[0].OffsetX != 0 {
		t.Errorf("locomotive offset = %d, want 0", cars[0].OffsetX)
	}
	if cars[1].OffsetX >= 0 {
		t.Errorf("carriage offset = %d, want negative", cars[1].OffsetX)
	}
}

func TestValidateConsist(t *testing.T) {
	cars := SimpleConsist(testStyle())
	if err := ValidateConsist(cars); err != nil {
		t.Errorf("valid consist failed validation: %v", err)
	}
}

func TestValidateConsistRejectsNilSprite(t *testing.T) {
	cars := []Car{{Sprite: nil, OffsetX: 0}}
	err := ValidateConsist(cars)
	if err == nil {
		t.Error("consist with nil sprite should fail validation")
	}
}

func TestValidateConsistRejectsHeightMismatch(t *testing.T) {
	// Create two sprites with different heights.
	short := render.NewSprite("short",
		render.ParseFrame("AB", testStyle()))
	tall := render.NewSprite("tall",
		render.ParseFrame("AB\nCD", testStyle()))
	cars := []Car{
		{Sprite: short, OffsetX: 0},
		{Sprite: tall, OffsetX: -3},
	}
	err := ValidateConsist(cars)
	if err == nil {
		t.Error("consist with height mismatch should fail validation")
	}
}
