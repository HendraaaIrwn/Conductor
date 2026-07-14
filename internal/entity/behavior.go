package entity

import "time"

// LinearMovement updates an entity's position based on its velocity and the
// frame delta. This is the most common behavior for entities that move in a
// straight line (trains, birds, clouds, particles).
type LinearMovement struct{}

// Update advances position: X += VX * delta, Y += VY * delta.
func (LinearMovement) Update(e *Entity, _ UpdateContext, delta time.Duration) {
	sec := delta.Seconds()
	e.X += e.VX * sec
	e.Y += e.VY * sec
}

// Animation cycles through an entity's sprite frames at a fixed interval. It
// is typically composed with another behavior (e.g. LinearMovement) via
// Composite.
type Animation struct {
	FrameTime time.Duration // duration of each frame
	Elapsed   time.Duration // accumulated time since last frame change
}

// Update advances the animation timer and wraps the frame index.
func (a *Animation) Update(e *Entity, _ UpdateContext, delta time.Duration) {
	if e.Sprite == nil || e.Sprite.FrameCount() <= 1 {
		return
	}
	a.Elapsed += delta
	for a.Elapsed >= a.FrameTime {
		a.Elapsed -= a.FrameTime
		e.Frame++
		if e.Frame >= e.Sprite.FrameCount() {
			e.Frame = 0
		}
	}
}

// Composite runs multiple behaviors in order on the same entity. This allows
// combining movement and animation without writing a custom behavior.
type Composite struct {
	Behaviors []Behavior
}

// Update calls each child behavior in sequence.
func (c Composite) Update(e *Entity, ctx UpdateContext, delta time.Duration) {
	for _, b := range c.Behaviors {
		b.Update(e, ctx, delta)
	}
}

// FuncBehavior wraps a function as a Behavior. It is useful for one-off
// behaviors that don't warrant a dedicated type.
type FuncBehavior func(e *Entity, ctx UpdateContext, delta time.Duration)

// Update calls the wrapped function.
func (f FuncBehavior) Update(e *Entity, ctx UpdateContext, delta time.Duration) {
	f(e, ctx, delta)
}
