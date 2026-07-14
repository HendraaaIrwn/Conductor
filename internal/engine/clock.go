// Package engine contains the core animation loop, timing, and world state.
// The engine is intentionally decoupled from the terminal: the Loop drives the
// Clock and World, and only the Loop talks to the terminal through the render
// Canvas.
package engine

import "time"

// Clock tracks frame timing and produces delta-time values used for
// frame-rate-independent movement. When paused, Tick returns a zero delta so
// that entities freeze in place.
type Clock struct {
	fps       int
	speed     float64
	paused    bool
	lastTick  time.Time
	firstTick bool
}

// NewClock creates a Clock targeting the given FPS. The speed multiplier
// defaults to 1.0.
func NewClock(fps int) *Clock {
	if fps < 1 {
		fps = 1
	}
	return &Clock{
		fps:       fps,
		speed:     1.0,
		firstTick: true,
	}
}

// FPS returns the configured target frame rate.
func (c *Clock) FPS() int { return c.fps }

// SetFPS changes the target frame rate.
func (c *Clock) SetFPS(fps int) {
	if fps < 1 {
		fps = 1
	}
	c.fps = fps
}

// Speed returns the current animation speed multiplier.
func (c *Clock) Speed() float64 { return c.speed }

// SetSpeed sets the animation speed multiplier. The value is clamped to the
// range [0.25, 3.0].
func (c *Clock) SetSpeed(speed float64) {
	if speed < 0.25 {
		speed = 0.25
	}
	if speed > 3.0 {
		speed = 3.0
	}
	c.speed = speed
}

// Paused reports whether the clock is currently paused.
func (c *Clock) Paused() bool { return c.paused }

// SetPaused pauses or resumes the clock.
func (c *Clock) SetPaused(paused bool) {
	c.paused = paused
	// Reset the last tick when unpausing so that the pause duration is not
	// counted as elapsed game time.
	if !paused {
		c.lastTick = time.Now()
	}
}

// Tick advances the clock and returns the delta time in seconds scaled by the
// speed multiplier. When paused it returns 0. On the very first call it
// returns 0 and records the start time.
func (c *Clock) Tick() float64 {
	now := time.Now()
	if c.firstTick {
		c.firstTick = false
		c.lastTick = now
		return 0
	}
	raw := now.Sub(c.lastTick).Seconds()
	c.lastTick = now
	if c.paused {
		return 0
	}
	return raw * c.speed
}

// FrameDuration returns the target duration of a single frame based on the
// configured FPS. It is used by the loop to sleep between frames.
func (c *Clock) FrameDuration() time.Duration {
	return time.Second / time.Duration(c.fps)
}
