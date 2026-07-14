package engine

import (
	"math/rand"
	"time"
)

// EventType identifies a random event.
type EventType int

const (
	EventNone EventType = iota
	EventLongFreight
	EventSmokeBurst
	EventBirdsCrossing
	EventLightning
	EventSignalChange
)

// String returns the human-readable name of the event type.
func (e EventType) String() string {
	switch e {
	case EventLongFreight:
		return "long-freight"
	case EventSmokeBurst:
		return "smoke-burst"
	case EventBirdsCrossing:
		return "birds"
	case EventLightning:
		return "lightning"
	case EventSignalChange:
		return "signal-change"
	default:
		return "none"
	}
}

// EventCallback is called when an event fires. It receives the event type
// and the world's RNG for any randomness needed during execution.
type EventCallback func(et EventType, rng *rand.Rand)

// ScheduledEvent is a pending event with a trigger time.
type ScheduledEvent struct {
	Type      EventType
	TriggerAt time.Duration
}

// Scheduler manages random events with cooldowns and probabilities. It is
// deterministic when given a fixed seed.
type Scheduler struct {
	rng       *rand.Rand
	cooldown  time.Duration
	lastEvent time.Duration
	elapsed   time.Duration
	minDelay  time.Duration
	maxDelay  time.Duration
	enabled   bool
	callbacks map[EventType]EventCallback
	pending   []ScheduledEvent
}

// NewScheduler creates a Scheduler with the given seed.
func NewScheduler(seed int64) *Scheduler {
	return &Scheduler{
		rng:       rand.New(rand.NewSource(seed)),
		cooldown:  15 * time.Second,
		minDelay:  10 * time.Second,
		maxDelay:  45 * time.Second,
		enabled:   true,
		callbacks: make(map[EventType]EventCallback),
	}
}

// SetEnabled enables or disables the scheduler. When disabled, no events
// are scheduled or fired.
func (s *Scheduler) SetEnabled(enabled bool) {
	s.enabled = enabled
}

// On registers a callback for a specific event type.
func (s *Scheduler) On(et EventType, cb EventCallback) {
	s.callbacks[et] = cb
}

// Reset resets the scheduler state (e.g. on regenerate).
func (s *Scheduler) Reset(seed int64) {
	s.rng = rand.New(rand.NewSource(seed))
	s.elapsed = 0
	s.lastEvent = 0
	s.pending = s.pending[:0]
}

// Update advances the scheduler by the given delta. It checks pending events
// and schedules new ones. When an event fires, its callback is invoked.
func (s *Scheduler) Update(delta time.Duration) {
	if !s.enabled || delta <= 0 {
		return
	}
	s.elapsed += delta

	// Check pending events.
	remaining := s.pending[:0]
	for _, ev := range s.pending {
		if s.elapsed >= ev.TriggerAt {
			s.fire(ev.Type)
		} else {
			remaining = append(remaining, ev)
		}
	}
	s.pending = remaining

	// Maybe schedule a new event if cooldown has passed.
	if s.elapsed-s.lastEvent >= s.cooldown && len(s.pending) == 0 {
		s.maybeSchedule()
	}
}

// fire invokes the callback for the given event type and updates the last
// event time.
func (s *Scheduler) fire(et EventType) {
	s.lastEvent = s.elapsed
	if cb, ok := s.callbacks[et]; ok {
		cb(et, s.rng)
	}
}

// maybeSchedule probabilistically schedules a new random event.
func (s *Scheduler) maybeSchedule() {
	// 30% chance per check to schedule an event.
	if s.rng.Intn(10) >= 3 {
		return
	}
	delay := s.minDelay + time.Duration(s.rng.Int63n(int64(s.maxDelay-s.minDelay)))
	et := s.pickEvent()
	s.pending = append(s.pending, ScheduledEvent{
		Type:      et,
		TriggerAt: s.elapsed + delay,
	})
}

// pickEvent selects a random event type with weighted probabilities.
func (s *Scheduler) pickEvent() EventType {
	events := []EventType{
		EventLongFreight,
		EventSmokeBurst,
		EventBirdsCrossing,
		EventLightning,
		EventSignalChange,
	}
	return events[s.rng.Intn(len(events))]
}

// PendingCount returns the number of pending events.
func (s *Scheduler) PendingCount() int { return len(s.pending) }

// Elapsed returns the total elapsed time since the scheduler started.
func (s *Scheduler) Elapsed() time.Duration { return s.elapsed }

// HasFiredSince reports whether any event has fired since the given time.
func (s *Scheduler) HasFiredSince(since time.Duration) bool {
	return s.lastEvent > since
}
