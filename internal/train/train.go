package train

import (
	"time"

	"github.com/example/conductor/internal/entity"
	"github.com/example/conductor/internal/render"
)

// Car is a single unit in a train consist: either the locomotive or a
// carriage. OffsetX is the horizontal offset from the train's leading unit,
// measured in cells. A negative offset means the car is behind the leader.
type Car struct {
	Sprite  *render.Sprite
	OffsetX int
}

// SmokeFactory creates a smoke particle entity at the given position. The
// train behavior calls this to emit smoke from the chimney. The world provides
// the real factory (effects.NewSmoke); tests can provide a stub.
type SmokeFactory func(x, y float64) *entity.Entity

// TrainData holds the train-specific state stored in an entity.Entity's Data
// field. It is read by the train's RenderFunc and updated by TrainBehavior.
type TrainData struct {
	Cars           []Car
	Direction      Direction
	WheelFrame     int
	WheelElapsed   time.Duration
	WheelFrameTime time.Duration
	SmokeElapsed   time.Duration
	SmokeInterval  time.Duration
}

// TrainBehavior updates a train entity each frame: it moves the entity, cycles
// the wheel animation, and spawns smoke particles at the chimney.
type TrainBehavior struct {
	Speed        float64      // cells per second (always positive)
	SmokeFactory SmokeFactory // if nil, no smoke is spawned
}

// Update advances the train's position, wheel animation, and smoke timer.
func (b *TrainBehavior) Update(e *entity.Entity, ctx entity.UpdateContext, delta time.Duration) {
	data, ok := e.Data.(*TrainData)
	if !ok || data == nil {
		return
	}

	// Move the train.
	sec := delta.Seconds()
	e.X += e.VX * sec

	// Advance wheel animation.
	data.WheelElapsed += delta
	for data.WheelElapsed >= data.WheelFrameTime {
		data.WheelElapsed -= data.WheelFrameTime
		data.WheelFrame++
	}

	// Spawn smoke if enough time has elapsed and a factory is configured.
	if b.SmokeFactory == nil {
		return
	}
	data.SmokeElapsed += delta
	if data.SmokeElapsed >= data.SmokeInterval {
		data.SmokeElapsed = 0
		smokeX, smokeY := chimneyPosition(e, data)
		ctx.Spawn(b.SmokeFactory(smokeX, smokeY))
	}
}

// chimneyPosition returns the canvas coordinates where smoke should be
// emitted. For a steam locomotive this is the top of the chimney, near the
// front of the locomotive.
func chimneyPosition(e *entity.Entity, data *TrainData) (float64, float64) {
	if len(data.Cars) == 0 {
		return e.X, e.Y
	}
	loco := data.Cars[0]
	spriteW := loco.Sprite.Width
	chimneyOffset := spriteW - 3
	if data.Direction == RightToLeft {
		chimneyOffset = 2
	}
	return e.X + float64(chimneyOffset), e.Y - 1
}

// TrainRenderFunc draws all cars of a train onto the canvas. It is intended to
// be set as the entity's RenderFunc.
func TrainRenderFunc(canvas *render.Canvas, e *entity.Entity) {
	data, ok := e.Data.(*TrainData)
	if !ok || data == nil {
		return
	}
	for _, car := range data.Cars {
		carX := int(e.X) + car.OffsetX
		carY := int(e.Y)
		car.Sprite.Draw(canvas, carX, carY, data.WheelFrame)
	}
}

// IsOffscreen reports whether the entire train entity has left the viewport.
func IsOffscreen(e *entity.Entity, width int) bool {
	data, ok := e.Data.(*TrainData)
	if !ok || data == nil {
		return false
	}
	if e.VX > 0 {
		return leftEdge(e, data) > width
	}
	return rightEdge(e, data) < 0
}

// leftEdge returns the X coordinate of the leftmost cell of the train.
func leftEdge(e *entity.Entity, data *TrainData) int {
	minX := 0
	for _, car := range data.Cars {
		if car.OffsetX < minX {
			minX = car.OffsetX
		}
	}
	return int(e.X) + minX
}

// rightEdge returns the X coordinate just past the rightmost cell of the train.
func rightEdge(e *entity.Entity, data *TrainData) int {
	maxX := 0
	for _, car := range data.Cars {
		if right := car.OffsetX + car.Sprite.Width; right > maxX {
			maxX = right
		}
	}
	return int(e.X) + maxX
}

// TrainWidth returns the total width of a train in cells given its car list.
func TrainWidth(cars []Car) int {
	if len(cars) == 0 {
		return 0
	}
	minX := 0
	maxX := 0
	for _, car := range cars {
		if car.OffsetX < minX {
			minX = car.OffsetX
		}
		if right := car.OffsetX + car.Sprite.Width; right > maxX {
			maxX = right
		}
	}
	return maxX - minX
}

// NewTrainEntity creates an entity.Entity representing a complete train. The
// entity has a TrainBehavior, TrainRenderFunc, and TrainData stored in the
// Data field. The caller should add the returned entity to the manager.
func NewTrainEntity(x float64, y int, speed float64, dir Direction, cars []Car, smokeFactory SmokeFactory) *entity.Entity {
	data := &TrainData{
		Cars:           cars,
		Direction:      dir,
		WheelFrameTime: 150 * time.Millisecond,
		SmokeInterval:  400 * time.Millisecond,
	}
	return &entity.Entity{
		Type:       entity.TypeTrain,
		X:          x,
		Y:          float64(y),
		VX:         speed * VelocityForDirection(dir),
		Layer:      render.LayerTrain,
		Visible:    true,
		Behavior:   &TrainBehavior{Speed: speed, SmokeFactory: smokeFactory},
		RenderFunc: TrainRenderFunc,
		Data:       data,
		// Trains are not auto-removed by the manager; the world checks
		// IsOffscreen explicitly because the train's bounding box is
		// computed from its cars, not from a single sprite.
		RemoveOffscreen: false,
	}
}
