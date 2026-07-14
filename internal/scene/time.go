package scene

import (
	"math/rand"

	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// TimePeriod represents a time of day that affects the scene's palette and
// celestial elements (sun, moon, stars).
type TimePeriod int

const (
	// TimeMorning is early morning with warm light.
	TimeMorning TimePeriod = iota
	// TimeDay is full daylight.
	TimeDay
	// TimeSunset is evening with orange tones.
	TimeSunset
	// TimeNight is dark with stars and moon.
	TimeNight
)

// String returns the human-readable name of the time period.
func (t TimePeriod) String() string {
	switch t {
	case TimeMorning:
		return "morning"
	case TimeDay:
		return "day"
	case TimeSunset:
		return "sunset"
	case TimeNight:
		return "night"
	default:
		return "unknown"
	}
}

// ParseTimePeriod converts a string to a TimePeriod. Returns false if the
// string does not match a known period.
func ParseTimePeriod(s string) (TimePeriod, bool) {
	switch s {
	case "morning":
		return TimeMorning, true
	case "day":
		return TimeDay, true
	case "sunset":
		return TimeSunset, true
	case "night":
		return TimeNight, true
	default:
		return 0, false
	}
}

// AllTimePeriods returns all time periods for cycling and random selection.
func AllTimePeriods() []TimePeriod {
	return []TimePeriod{TimeMorning, TimeDay, TimeSunset, TimeNight}
}

// Palette returns the color palette for this time period. The palette
// contains sky, celestial, and scenery colors that scenes use when rendering
// backgrounds.
type Palette struct {
	SkyColor    tcell.Color
	SunColor    tcell.Color
	MoonColor   tcell.Color
	StarColor   tcell.Color
	HillColor   tcell.Color
	GroundColor tcell.Color
	TreeColor   tcell.Color
	DimFactor   bool // if true, styles are dimmed (night/sunset)
	ShowStars   bool
	ShowSun     bool
	ShowMoon    bool
}

// PaletteFor returns the color palette for the given time period.
func PaletteFor(t TimePeriod) Palette {
	switch t {
	case TimeMorning:
		return Palette{
			SkyColor:    tcell.ColorLightYellow,
			SunColor:    tcell.ColorYellow,
			HillColor:   tcell.ColorGreen,
			GroundColor: tcell.ColorDarkGreen,
			TreeColor:   tcell.ColorGreen,
			DimFactor:   false,
			ShowStars:   false,
			ShowSun:     true,
			ShowMoon:    false,
		}
	case TimeDay:
		return Palette{
			SkyColor:    tcell.ColorWhite,
			SunColor:    tcell.ColorYellow,
			HillColor:   tcell.ColorGreen,
			GroundColor: tcell.ColorDarkGreen,
			TreeColor:   tcell.ColorGreen,
			DimFactor:   false,
			ShowStars:   false,
			ShowSun:     true,
			ShowMoon:    false,
		}
	case TimeSunset:
		return Palette{
			SkyColor:    tcell.ColorDarkOrange,
			SunColor:    tcell.ColorOrange,
			MoonColor:   tcell.ColorSilver,
			HillColor:   tcell.ColorDarkGreen,
			GroundColor: tcell.ColorBlack,
			TreeColor:   tcell.ColorDarkGreen,
			DimFactor:   true,
			ShowStars:   false,
			ShowSun:     true,
			ShowMoon:    false,
		}
	case TimeNight:
		return Palette{
			SkyColor:    tcell.ColorDarkBlue,
			MoonColor:   tcell.ColorSilver,
			StarColor:   tcell.ColorWhite,
			HillColor:   tcell.ColorDarkBlue,
			GroundColor: tcell.ColorBlack,
			TreeColor:   tcell.ColorNavy,
			DimFactor:   true,
			ShowStars:   true,
			ShowSun:     false,
			ShowMoon:    true,
		}
	default:
		return PaletteFor(TimeDay)
	}
}

// ApplyToStyle returns a new tcell.Style with the palette's foreground color
// applied, and dimmed if the palette's DimFactor is true.
func (p Palette) ApplyToStyle(style tcell.Style) tcell.Style {
	if p.DimFactor {
		style = style.Dim(true)
	}
	return style
}

// drawCelestial draws the sun or moon and stars based on the time period's
// palette. This is called by scenes during RenderBackground.
func drawCelestial(canvas *render.Canvas, vp Viewport, pal Palette, rng *rand.Rand) {
	// Draw stars at night.
	if pal.ShowStars {
		density := vp.Width / 5
		if density < 10 {
			density = 10
		}
		starStyle := tcell.StyleDefault.Foreground(pal.StarColor)
		starY := vp.TrackY() - 2
		if starY < 1 {
			starY = 1
		}
		for i := 0; i < density; i++ {
			x := rng.Intn(vp.Width)
			y := rng.Intn(starY)
			canvas.SetRune(x, y, '*', starStyle)
		}
	}

	// Draw sun.
	if pal.ShowSun {
		sunX := vp.Width * 3 / 4
		sunY := 3
		if sunY >= vp.TrackY() {
			sunY = vp.TrackY() - 2
		}
		sunStyle := tcell.StyleDefault.Foreground(pal.SunColor).Bold(true)
		canvas.SetRune(sunX, sunY, '@', sunStyle)
		canvas.SetRune(sunX-1, sunY, '-', sunStyle)
		canvas.SetRune(sunX+1, sunY, '-', sunStyle)
	}

	// Draw moon.
	if pal.ShowMoon {
		moonX := vp.Width * 3 / 4
		moonY := 3
		if moonY >= vp.TrackY() {
			moonY = vp.TrackY() - 2
		}
		moonStyle := tcell.StyleDefault.Foreground(pal.MoonColor).Bold(true)
		canvas.SetRune(moonX, moonY, 'C', moonStyle)
		canvas.SetRune(moonX+1, moonY, ')', moonStyle)
	}
}
