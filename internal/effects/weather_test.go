package effects

import (
	"testing"
	"time"

	"github.com/example/conductor/internal/entity"
)

func TestParseWeatherType(t *testing.T) {
	cases := []struct {
		input string
		want  WeatherType
		ok    bool
	}{
		{"clear", WeatherClear, true},
		{"rain", WeatherRain, true},
		{"snow", WeatherSnow, true},
		{"unknown", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseWeatherType(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseWeatherType(%q) ok = %v, want %v", tc.input, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Errorf("ParseWeatherType(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestWeatherTypeString(t *testing.T) {
	if WeatherClear.String() != "clear" {
		t.Errorf("clear string = %q", WeatherClear.String())
	}
	if WeatherRain.String() != "rain" {
		t.Errorf("rain string = %q", WeatherRain.String())
	}
	if WeatherSnow.String() != "snow" {
		t.Errorf("snow string = %q", WeatherSnow.String())
	}
}

func TestAllWeatherTypes(t *testing.T) {
	types := AllWeatherTypes()
	if len(types) != 3 {
		t.Errorf("AllWeatherTypes length = %d, want 3", len(types))
	}
}

func TestNewWeatherSystem(t *testing.T) {
	ws := NewWeatherSystem(WeatherClear, false)
	if ws.Current() != WeatherClear {
		t.Error("initial weather should be clear")
	}
}

func TestWeatherSystemSetWeather(t *testing.T) {
	ws := NewWeatherSystem(WeatherClear, false)
	mgr := entity.NewManager()
	ws.SetWeather(WeatherRain, mgr, 120, 40)
	mgr.Flush()
	if ws.Current() != WeatherRain {
		t.Error("weather should be rain after SetWeather")
	}
	rain := mgr.ByType(entity.TypeRain)
	if len(rain) == 0 {
		t.Error("rain particles should be spawned")
	}
}

func TestWeatherSystemSetWeatherClearRemovesParticles(t *testing.T) {
	ws := NewWeatherSystem(WeatherRain, false)
	mgr := entity.NewManager()
	ws.SetWeather(WeatherRain, mgr, 120, 40)
	mgr.Flush()
	ws.SetWeather(WeatherClear, mgr, 120, 40)
	mgr.Flush()
	rain := mgr.ByType(entity.TypeRain)
	if len(rain) != 0 {
		t.Error("rain particles should be removed when weather is clear")
	}
}

func TestWeatherSystemCycleWeather(t *testing.T) {
	ws := NewWeatherSystem(WeatherClear, false)
	mgr := entity.NewManager()
	// clear → rain
	ws.CycleWeather(mgr, 120, 40)
	if ws.Current() != WeatherRain {
		t.Errorf("after cycle: %v, want rain", ws.Current())
	}
	// rain → snow
	ws.CycleWeather(mgr, 120, 40)
	if ws.Current() != WeatherSnow {
		t.Errorf("after cycle: %v, want snow", ws.Current())
	}
	// snow → clear
	ws.CycleWeather(mgr, 120, 40)
	if ws.Current() != WeatherClear {
		t.Errorf("after cycle: %v, want clear", ws.Current())
	}
}

func TestWeatherSystemReducedMotionFewerParticles(t *testing.T) {
	mgr1 := entity.NewManager()
	ws1 := NewWeatherSystem(WeatherClear, true)
	ws1.SetWeather(WeatherRain, mgr1, 120, 40)
	mgr1.Flush()
	count1 := len(mgr1.ByType(entity.TypeRain))

	mgr2 := entity.NewManager()
	ws2 := NewWeatherSystem(WeatherClear, false)
	ws2.SetWeather(WeatherRain, mgr2, 120, 40)
	mgr2.Flush()
	count2 := len(mgr2.ByType(entity.TypeRain))

	if count1 == 0 {
		t.Error("reduced motion should still spawn some rain particles")
	}
	if count1 >= count2 {
		t.Errorf("reduced motion should have fewer particles: %d >= %d", count1, count2)
	}
}

func TestRainParticleCount(t *testing.T) {
	count := RainParticleCount(120, 40)
	if count < 20 || count > 150 {
		t.Errorf("rain count = %d, want 20-150", count)
	}
}

func TestSnowParticleCount(t *testing.T) {
	count := SnowParticleCount(120, 40)
	if count < 10 || count > 80 {
		t.Errorf("snow count = %d, want 10-80", count)
	}
}

func TestSnowParticleCountLessThanRain(t *testing.T) {
	rain := RainParticleCount(120, 40)
	snow := SnowParticleCount(120, 40)
	if snow >= rain {
		t.Errorf("snow count %d should be less than rain count %d", snow, rain)
	}
}

func TestRainBehaviorFallsDown(t *testing.T) {
	e := NewRain(10, 5, 10)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	initialY := e.Y
	e.Behavior.Update(e, ctx, 100*time.Millisecond)
	if e.Y <= initialY {
		t.Errorf("rain should fall: Y = %f, initial = %f", e.Y, initialY)
	}
}

func TestRainBehaviorRecycles(t *testing.T) {
	e := NewRain(10, 39, 10)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	// Move past the bottom edge.
	e.Behavior.Update(e, ctx, 200*time.Millisecond)
	if e.Y >= 40 {
		t.Errorf("rain should have recycled: Y = %f", e.Y)
	}
}

func TestSnowBehaviorFallsSlowly(t *testing.T) {
	e := NewSnow(10, 5, 10, 0)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	initialY := e.Y
	e.Behavior.Update(e, ctx, 100*time.Millisecond)
	if e.Y <= initialY {
		t.Errorf("snow should fall: Y = %f, initial = %f", e.Y, initialY)
	}
	// Snow should fall slower than rain.
	rainE := NewRain(10, 5, 10)
	rainE.Behavior.Update(rainE, ctx, 100*time.Millisecond)
	snowFall := e.Y - initialY
	rainFall := rainE.Y - initialY
	if snowFall >= rainFall {
		t.Errorf("snow fall %f should be slower than rain fall %f", snowFall, rainFall)
	}
}

func TestSnowBehaviorRecycles(t *testing.T) {
	e := NewSnow(10, 39, 10, 0)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	// Snow is slow, so we need multiple updates to reach the bottom.
	for i := 0; i < 20; i++ {
		e.Behavior.Update(e, ctx, 200*time.Millisecond)
		e.Age += 200 * time.Millisecond
	}
	if e.Y >= 40 {
		t.Errorf("snow should have recycled: Y = %f", e.Y)
	}
}

func TestLightningFlashTrigger(t *testing.T) {
	flash := NewLightningFlash()
	if flash.Active {
		t.Error("flash should not be active initially")
	}
	flash.Trigger()
	if !flash.Active {
		t.Error("flash should be active after Trigger")
	}
}

func TestLightningFlashFadesOut(t *testing.T) {
	flash := NewLightningFlash()
	flash.Trigger()
	flash.Update(50 * time.Millisecond)
	if !flash.Active {
		t.Error("flash should still be active after 50ms")
	}
	if flash.Intensity() <= 0 {
		t.Error("flash intensity should be > 0")
	}
	flash.Update(200 * time.Millisecond)
	if flash.Active {
		t.Error("flash should be inactive after 250ms total (exceeds duration)")
	}
	if flash.Intensity() != 0 {
		t.Errorf("inactive flash intensity = %f, want 0", flash.Intensity())
	}
}

func TestLightningFlashIntensityDecreases(t *testing.T) {
	flash := NewLightningFlash()
	flash.Trigger()
	i1 := flash.Intensity()
	flash.Update(100 * time.Millisecond)
	i2 := flash.Intensity()
	if i2 >= i1 {
		t.Errorf("intensity should decrease: %f -> %f", i1, i2)
	}
}

func TestBirdFliesHorizontally(t *testing.T) {
	e := NewBird(10, 5, 15.0)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	initialX := e.X
	e.Behavior.Update(e, ctx, 100*time.Millisecond)
	if e.X == initialX {
		t.Error("bird should move horizontally")
	}
}

func TestBirdWingAnimation(t *testing.T) {
	e := NewBird(10, 5, 15.0)
	ctx := entity.UpdateContext{Width: 120, Height: 40, Spawn: func(*entity.Entity) {}}
	initialFrame := e.Frame
	e.Behavior.Update(e, ctx, 250*time.Millisecond)
	if e.Frame == initialFrame {
		t.Error("bird wing frame should advance")
	}
}

func TestMaybeLightningOnlyDuringRain(t *testing.T) {
	ws := NewWeatherSystem(WeatherClear, false)
	// Use a mock RNG that always returns 0 (would trigger if weather were rain).
	rng := &alwaysZeroRNG{}
	if ws.MaybeLightning(rng) {
		t.Error("lightning should not trigger during clear weather")
	}
}

func TestMaybeLightningNotWithReducedMotion(t *testing.T) {
	ws := NewWeatherSystem(WeatherRain, true)
	rng := &alwaysZeroRNG{}
	if ws.MaybeLightning(rng) {
		t.Error("lightning should not trigger with reduced motion")
	}
}

func TestMaybeLightningDuringRain(t *testing.T) {
	ws := NewWeatherSystem(WeatherRain, false)
	rng := &alwaysZeroRNG{}
	// With rng returning 0, MaybeLightning should trigger.
	if !ws.MaybeLightning(rng) {
		t.Error("lightning should trigger during rain with rng=0")
	}
	if !ws.Flash().Active {
		t.Error("flash should be active after lightning trigger")
	}
}

// alwaysZeroRNG always returns 0 from Intn, ensuring lightning triggers.
type alwaysZeroRNG struct{}

func (r *alwaysZeroRNG) Intn(_ int) int { return 0 }
