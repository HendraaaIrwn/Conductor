package input

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// mockWorld implements WorldController for testing.
type mockWorld struct {
	paused         bool
	regenerated    bool
	spawned        bool
	sceneChanged   bool
	weatherChanged bool
	timeChanged    bool
	colorToggled   bool
	helpToggled    bool
	helpShown      bool
	speed          float64
}

func (m *mockWorld) Paused() bool          { return m.paused }
func (m *mockWorld) SetPaused(p bool)      { m.paused = p }
func (m *mockWorld) Regenerate()           { m.regenerated = true }
func (m *mockWorld) SpawnNext()            { m.spawned = true }
func (m *mockWorld) ChangeScene()          { m.sceneChanged = true }
func (m *mockWorld) ChangeWeather()        { m.weatherChanged = true }
func (m *mockWorld) ChangeTime()           { m.timeChanged = true }
func (m *mockWorld) ChangeColorPalette()   { m.colorToggled = true }
func (m *mockWorld) ToggleHelp()           { m.helpToggled = true }
func (m *mockWorld) SetHelpVisible(v bool) { m.helpShown = v }
func (m *mockWorld) HelpVisible() bool     { return m.helpShown }
func (m *mockWorld) Speed() float64        { return m.speed }
func (m *mockWorld) SetSpeed(s float64)    { m.speed = s }

// mockClock implements ClockController for testing.
type mockClock struct {
	paused bool
}

func (m *mockClock) SetPaused(p bool) { m.paused = p }

func keyEvent(r rune, key tcell.Key) tcell.Event {
	return tcell.NewEventKey(key, r, tcell.ModNone)
}

func TestQuitKeyQ(t *testing.T) {
	h := NewHandler(&mockWorld{}, &mockClock{})
	if !h.Handle(keyEvent('q', tcell.KeyRune)) {
		t.Error("'q' should return true (quit)")
	}
}

func TestQuitKeyQUpper(t *testing.T) {
	h := NewHandler(&mockWorld{}, &mockClock{})
	if !h.Handle(keyEvent('Q', tcell.KeyRune)) {
		t.Error("'Q' should return true (quit)")
	}
}

func TestQuitCtrlC(t *testing.T) {
	h := NewHandler(&mockWorld{}, &mockClock{})
	if !h.Handle(keyEvent(0, tcell.KeyCtrlC)) {
		t.Error("Ctrl+C should return true (quit)")
	}
}

func TestPauseKeyP(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('p', tcell.KeyRune))
	if !w.paused {
		t.Error("'p' should toggle pause")
	}
}

func TestPauseSpace(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent(' ', tcell.KeyRune))
	if !w.paused {
		t.Error("space should toggle pause")
	}
}

func TestRegenerate(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('r', tcell.KeyRune))
	if !w.regenerated {
		t.Error("'r' should regenerate")
	}
}

func TestSpawnNext(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('n', tcell.KeyRune))
	if !w.spawned {
		t.Error("'n' should spawn next train")
	}
}

func TestChangeScene(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('s', tcell.KeyRune))
	if !w.sceneChanged {
		t.Error("'s' should change scene")
	}
}

func TestChangeWeather(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('w', tcell.KeyRune))
	if !w.weatherChanged {
		t.Error("'w' should change weather")
	}
}

func TestChangeTime(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('d', tcell.KeyRune))
	if !w.timeChanged {
		t.Error("'d' should change time")
	}
}

func TestToggleColor(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('c', tcell.KeyRune))
	if !w.colorToggled {
		t.Error("'c' should toggle color palette")
	}
}

func TestToggleHelp(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('h', tcell.KeyRune))
	if !w.helpToggled {
		t.Error("'h' should toggle help")
	}
}

func TestToggleHelpQuestion(t *testing.T) {
	w := &mockWorld{}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('?', tcell.KeyRune))
	if !w.helpToggled {
		t.Error("'?' should toggle help")
	}
}

func TestEscapeClosesHelp(t *testing.T) {
	w := &mockWorld{helpShown: true}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent(0, tcell.KeyEscape))
	if w.helpShown {
		t.Error("Esc should close help")
	}
}

func TestEscapeNoopWhenHelpNotVisible(t *testing.T) {
	w := &mockWorld{helpShown: false}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent(0, tcell.KeyEscape))
	// Should not crash, should not show help.
	if w.helpShown {
		t.Error("Esc should not show help when already hidden")
	}
}

func TestIncreaseSpeed(t *testing.T) {
	w := &mockWorld{speed: 25.0}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('+', tcell.KeyRune))
	if w.speed != 30.0 {
		t.Errorf("speed = %f, want 30.0", w.speed)
	}
}

func TestDecreaseSpeed(t *testing.T) {
	w := &mockWorld{speed: 25.0}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('-', tcell.KeyRune))
	if w.speed != 20.0 {
		t.Errorf("speed = %f, want 20.0", w.speed)
	}
}

func TestResetSpeed(t *testing.T) {
	w := &mockWorld{speed: 50.0}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('0', tcell.KeyRune))
	if w.speed != 25.0 {
		t.Errorf("speed = %f, want 25.0", w.speed)
	}
}

func TestSpeedClampMin(t *testing.T) {
	w := &mockWorld{speed: 5.0}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('-', tcell.KeyRune))
	if w.speed != 5.0 {
		t.Errorf("speed = %f, want 5.0 (clamped)", w.speed)
	}
}

func TestSpeedClampMax(t *testing.T) {
	w := &mockWorld{speed: 80.0}
	h := NewHandler(w, &mockClock{})
	h.Handle(keyEvent('+', tcell.KeyRune))
	if w.speed != 80.0 {
		t.Errorf("speed = %f, want 80.0 (clamped)", w.speed)
	}
}

func TestPauseAlsoPausesClock(t *testing.T) {
	w := &mockWorld{}
	c := &mockClock{}
	h := NewHandler(w, c)
	h.Handle(keyEvent('p', tcell.KeyRune))
	if !c.paused {
		t.Error("clock should be paused when world is paused")
	}
}

func TestUnpauseAlsoUnpausesClock(t *testing.T) {
	w := &mockWorld{paused: true}
	c := &mockClock{paused: true}
	h := NewHandler(w, c)
	h.Handle(keyEvent('p', tcell.KeyRune))
	if w.paused {
		t.Error("world should be unpaused")
	}
	if c.paused {
		t.Error("clock should be unpaused")
	}
}

func TestUnhandledKeyDoesNotCrash(t *testing.T) {
	h := NewHandler(&mockWorld{}, &mockClock{})
	// Various keys that should be ignored.
	for _, r := range []rune{'x', 'z', '1', '9', '&', '\n', '\t'} {
		if h.Handle(keyEvent(r, tcell.KeyRune)) {
			t.Errorf("key %q should not quit", r)
		}
	}
}

func TestUnhandledSpecialKeyDoesNotCrash(t *testing.T) {
	h := NewHandler(&mockWorld{}, &mockClock{})
	keys := []tcell.Key{
		tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight,
		tcell.KeyHome, tcell.KeyEnd, tcell.KeyPgUp, tcell.KeyPgDn,
		tcell.KeyF1, tcell.KeyF2, tcell.KeyF3, tcell.KeyF4, tcell.KeyF5,
		tcell.KeyTab, tcell.KeyBackspace, tcell.KeyDelete, tcell.KeyEnter,
	}
	for _, k := range keys {
		if h.Handle(keyEvent(0, k)) {
			t.Errorf("key %v should not quit", k)
		}
	}
}

func TestNonKeyEventDoesNotCrash(t *testing.T) {
	h := NewHandler(&mockWorld{}, &mockClock{})
	// A resize event should not be handled by the input handler.
	ev := tcell.NewEventResize(80, 24)
	if h.Handle(ev) {
		t.Error("resize event should not return true")
	}
}
