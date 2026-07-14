package engine

import (
	"testing"
	"time"
)

func TestNewClock(t *testing.T) {
	c := NewClock(20)
	if c.FPS() != 20 {
		t.Errorf("FPS = %d, want 20", c.FPS())
	}
	if c.Speed() != 1.0 {
		t.Errorf("Speed = %f, want 1.0", c.Speed())
	}
	if c.Paused() {
		t.Error("new clock should not be paused")
	}
}

func TestClockClampsFPS(t *testing.T) {
	c := NewClock(0)
	if c.FPS() != 1 {
		t.Errorf("FPS = %d, want 1 (clamped)", c.FPS())
	}
}

func TestClockSetSpeed(t *testing.T) {
	c := NewClock(20)
	c.SetSpeed(2.0)
	if c.Speed() != 2.0 {
		t.Errorf("Speed = %f, want 2.0", c.Speed())
	}
	// Speed should be clamped to [0.25, 3.0].
	c.SetSpeed(10.0)
	if c.Speed() != 3.0 {
		t.Errorf("Speed = %f, want 3.0 (clamped)", c.Speed())
	}
	c.SetSpeed(0.1)
	if c.Speed() != 0.25 {
		t.Errorf("Speed = %f, want 0.25 (clamped)", c.Speed())
	}
}

func TestClockTickReturnsZeroOnFirstCall(t *testing.T) {
	c := NewClock(20)
	delta := c.Tick()
	if delta != 0 {
		t.Errorf("first Tick delta = %f, want 0", delta)
	}
}

func TestClockTickReturnsPositiveDelta(t *testing.T) {
	c := NewClock(20)
	c.Tick() // first call returns 0
	time.Sleep(10 * time.Millisecond)
	delta := c.Tick()
	if delta <= 0 {
		t.Errorf("delta = %f, want > 0", delta)
	}
}

func TestClockTickReturnsZeroWhenPaused(t *testing.T) {
	c := NewClock(20)
	c.Tick()
	time.Sleep(10 * time.Millisecond)
	c.SetPaused(true)
	delta := c.Tick()
	if delta != 0 {
		t.Errorf("paused delta = %f, want 0", delta)
	}
}

func TestClockFrameDuration(t *testing.T) {
	c := NewClock(20)
	want := time.Second / 20
	if d := c.FrameDuration(); d != want {
		t.Errorf("FrameDuration = %v, want %v", d, want)
	}
}
