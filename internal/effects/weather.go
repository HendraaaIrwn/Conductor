package effects

import (
	"fmt"

	"github.com/example/conductor/internal/entity"
)

// WeatherType represents the current weather state.
type WeatherType int

const (
	// WeatherClear means no weather particles.
	WeatherClear WeatherType = iota
	// WeatherRain spawns rain particles.
	WeatherRain
	// WeatherSnow spawns snow particles.
	WeatherSnow
)

// String returns the human-readable name of the weather type.
func (w WeatherType) String() string {
	switch w {
	case WeatherClear:
		return "clear"
	case WeatherRain:
		return "rain"
	case WeatherSnow:
		return "snow"
	default:
		return "unknown"
	}
}

// ParseWeatherType converts a string to a WeatherType.
func ParseWeatherType(s string) (WeatherType, bool) {
	switch s {
	case "clear":
		return WeatherClear, true
	case "rain":
		return WeatherRain, true
	case "snow":
		return WeatherSnow, true
	default:
		return 0, false
	}
}

// AllWeatherTypes returns all weather types for cycling.
func AllWeatherTypes() []WeatherType {
	return []WeatherType{WeatherClear, WeatherRain, WeatherSnow}
}

// WeatherSystem manages weather particles. It creates and removes particle
// entities when the weather type changes, and scales particle counts with
// the viewport size.
type WeatherSystem struct {
	current       WeatherType
	reducedMotion bool
	flash         *LightningFlash
}

// NewWeatherSystem creates a WeatherSystem with the given initial weather.
func NewWeatherSystem(initial WeatherType, reducedMotion bool) *WeatherSystem {
	return &WeatherSystem{
		current:       initial,
		reducedMotion: reducedMotion,
		flash:         NewLightningFlash(),
	}
}

// Current returns the current weather type.
func (ws *WeatherSystem) Current() WeatherType { return ws.current }

// SetWeather changes the weather type. It removes old particles and spawns
// new ones appropriate for the new weather. Returns a description of the
// change for status display.
func (ws *WeatherSystem) SetWeather(w WeatherType, manager *entity.Manager, width, height int) {
	if w == ws.current {
		return
	}
	// Remove all existing weather particles.
	ws.clearParticles(manager)
	ws.current = w
	ws.spawnParticles(manager, width, height)
}

// CycleWeather advances to the next weather type (clear → rain → snow → clear).
func (ws *WeatherSystem) CycleWeather(manager *entity.Manager, width, height int) WeatherType {
	types := AllWeatherTypes()
	nextIdx := 0
	for i, t := range types {
		if t == ws.current {
			nextIdx = (i + 1) % len(types)
			break
		}
	}
	ws.SetWeather(types[nextIdx], manager, width, height)
	return ws.current
}

// clearParticles removes all rain and snow entities from the manager.
func (ws *WeatherSystem) clearParticles(manager *entity.Manager) {
	for _, e := range manager.ByType(entity.TypeRain) {
		e.Dead = true
	}
	for _, e := range manager.ByType(entity.TypeSnow) {
		e.Dead = true
	}
	manager.Flush()
}

// spawnParticles creates weather particles for the current weather type.
func (ws *WeatherSystem) spawnParticles(manager *entity.Manager, width, height int) {
	switch ws.current {
	case WeatherRain:
		count := RainParticleCount(width, height)
		if ws.reducedMotion {
			count /= 3
		}
		SpawnRain(manager, width, height, count)
	case WeatherSnow:
		count := SnowParticleCount(width, height)
		if ws.reducedMotion {
			count /= 3
		}
		SpawnSnow(manager, width, height, count)
	}
}

// HandleResize adjusts particle counts for the new viewport size.
func (ws *WeatherSystem) HandleResize(manager *entity.Manager, width, height int) {
	ws.clearParticles(manager)
	ws.spawnParticles(manager, width, height)
}

// Flash returns the lightning flash state for lightning effects during rain.
func (ws *WeatherSystem) Flash() *LightningFlash { return ws.flash }

// MaybeLightning triggers a lightning flash if the weather is rain and
// reduced motion is not enabled. The rng parameter determines probability.
func (ws *WeatherSystem) MaybeLightning(rng interface{ Intn(int) int }) bool {
	if ws.current != WeatherRain || ws.reducedMotion {
		return false
	}
	// 1 in 200 chance per check (called every few seconds by the scheduler).
	if rng.Intn(200) == 0 {
		ws.flash.Trigger()
		return true
	}
	return false
}

// Status returns a human-readable status string for the weather.
func (ws *WeatherSystem) Status() string {
	return fmt.Sprintf("weather: %s", ws.current)
}
