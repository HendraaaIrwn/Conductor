package engine

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	s := NewScheduler(42)
	if s.PendingCount() != 0 {
		t.Errorf("new scheduler should have 0 pending, got %d", s.PendingCount())
	}
}

func TestSchedulerDisabled(t *testing.T) {
	s := NewScheduler(42)
	s.SetEnabled(false)
	called := false
	s.On(EventLightning, func(_ EventType, _ *rand.Rand) {
		called = true
	})
	s.Update(60 * time.Second)
	if called {
		t.Error("disabled scheduler should not fire events")
	}
}

func TestSchedulerFiresAfterDelay(t *testing.T) {
	s := NewScheduler(42)
	s.SetEnabled(true)
	fired := false
	s.On(EventBirdsCrossing, func(_ EventType, _ *rand.Rand) {
		fired = true
	})
	s.On(EventLightning, func(_ EventType, _ *rand.Rand) {
		fired = true
	})
	s.On(EventSignalChange, func(_ EventType, _ *rand.Rand) {
		fired = true
	})
	s.On(EventSmokeBurst, func(_ EventType, _ *rand.Rand) {
		fired = true
	})
	s.On(EventLongFreight, func(_ EventType, _ *rand.Rand) {
		fired = true
	})
	// Advance past the cooldown and scheduling window.
	for i := 0; i < 1000 && !fired; i++ {
		s.Update(100 * time.Millisecond)
	}
	if !fired {
		t.Error("scheduler should have fired an event after 100 seconds")
	}
}

func TestSchedulerReset(t *testing.T) {
	s := NewScheduler(42)
	s.Update(60 * time.Second)
	if s.Elapsed() == 0 {
		t.Error("elapsed should be > 0 after Update")
	}
	s.Reset(99)
	if s.Elapsed() != 0 {
		t.Errorf("after Reset, elapsed = %v, want 0", s.Elapsed())
	}
	if s.PendingCount() != 0 {
		t.Errorf("after Reset, pending = %d, want 0", s.PendingCount())
	}
}

func TestSchedulerCooldownPreventsRapidFire(t *testing.T) {
	s := NewScheduler(42)
	s.SetEnabled(true)
	fireCount := 0
	cb := func(_ EventType, _ *rand.Rand) { fireCount++ }
	for _, et := range []EventType{EventBirdsCrossing, EventLightning, EventSignalChange, EventSmokeBurst, EventLongFreight} {
		s.On(et, cb)
	}
	// Run for 60 seconds.
	for i := 0; i < 600; i++ {
		s.Update(100 * time.Millisecond)
	}
	// With a 15-second cooldown, at most ~5 events in 60 seconds.
	if fireCount > 6 {
		t.Errorf("fire count = %d, should be limited by cooldown", fireCount)
	}
}

func TestSchedulerEventTypeString(t *testing.T) {
	cases := []struct {
		et   EventType
		want string
	}{
		{EventLongFreight, "long-freight"},
		{EventSmokeBurst, "smoke-burst"},
		{EventBirdsCrossing, "birds"},
		{EventLightning, "lightning"},
		{EventSignalChange, "signal-change"},
		{EventNone, "none"},
	}
	for _, tc := range cases {
		if got := tc.et.String(); got != tc.want {
			t.Errorf("%v.String() = %q, want %q", tc.et, got, tc.want)
		}
	}
}

func TestSchedulerHasFiredSince(t *testing.T) {
	s := NewScheduler(42)
	s.SetEnabled(true)
	cb := func(_ EventType, _ *rand.Rand) {}
	for _, et := range []EventType{EventBirdsCrossing, EventLightning, EventSignalChange, EventSmokeBurst, EventLongFreight} {
		s.On(et, cb)
	}
	// Run long enough for at least one event to fire.
	for i := 0; i < 1000; i++ {
		s.Update(100 * time.Millisecond)
	}
	if !s.HasFiredSince(0) {
		t.Error("HasFiredSince(0) should be true after running")
	}
}

func TestSchedulerDeterministic(t *testing.T) {
	s1 := NewScheduler(42)
	s2 := NewScheduler(42)
	s1.SetEnabled(true)
	s2.SetEnabled(true)

	count1 := 0
	count2 := 0
	cb1 := func(_ EventType, _ *rand.Rand) { count1++ }
	cb2 := func(_ EventType, _ *rand.Rand) { count2++ }
	for _, et := range []EventType{EventBirdsCrossing, EventLightning, EventSignalChange, EventSmokeBurst, EventLongFreight} {
		s1.On(et, cb1)
		s2.On(et, cb2)
	}

	for i := 0; i < 1000; i++ {
		s1.Update(100 * time.Millisecond)
		s2.Update(100 * time.Millisecond)
	}
	if count1 != count2 {
		t.Errorf("same seed produced different fire counts: %d vs %d", count1, count2)
	}
}
