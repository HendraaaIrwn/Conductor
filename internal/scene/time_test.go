package scene

import "testing"

func TestParseTimePeriod(t *testing.T) {
	cases := []struct {
		input string
		want  TimePeriod
		ok    bool
	}{
		{"morning", TimeMorning, true},
		{"day", TimeDay, true},
		{"sunset", TimeSunset, true},
		{"night", TimeNight, true},
		{"unknown", 0, false},
		{"", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseTimePeriod(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseTimePeriod(%q) ok = %v, want %v", tc.input, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Errorf("ParseTimePeriod(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestTimePeriodString(t *testing.T) {
	if TimeMorning.String() != "morning" {
		t.Errorf("morning string = %q", TimeMorning.String())
	}
	if TimeDay.String() != "day" {
		t.Errorf("day string = %q", TimeDay.String())
	}
	if TimeSunset.String() != "sunset" {
		t.Errorf("sunset string = %q", TimeSunset.String())
	}
	if TimeNight.String() != "night" {
		t.Errorf("night string = %q", TimeNight.String())
	}
}

func TestAllTimePeriods(t *testing.T) {
	periods := AllTimePeriods()
	if len(periods) != 4 {
		t.Errorf("AllTimePeriods length = %d, want 4", len(periods))
	}
}

func TestPaletteFor(t *testing.T) {
	for _, tp := range AllTimePeriods() {
		pal := PaletteFor(tp)
		// Each palette should have a defined state for sun/moon/stars.
		_ = pal
	}
}

func TestPaletteNightShowsStars(t *testing.T) {
	pal := PaletteFor(TimeNight)
	if !pal.ShowStars {
		t.Error("night palette should show stars")
	}
	if !pal.ShowMoon {
		t.Error("night palette should show moon")
	}
	if pal.ShowSun {
		t.Error("night palette should not show sun")
	}
}

func TestPaletteDayShowsSun(t *testing.T) {
	pal := PaletteFor(TimeDay)
	if !pal.ShowSun {
		t.Error("day palette should show sun")
	}
	if pal.ShowMoon {
		t.Error("day palette should not show moon")
	}
	if pal.ShowStars {
		t.Error("day palette should not show stars")
	}
}

func TestPaletteSunsetIsDimmed(t *testing.T) {
	pal := PaletteFor(TimeSunset)
	if !pal.DimFactor {
		t.Error("sunset palette should have DimFactor=true")
	}
}

func TestPaletteNightIsDimmed(t *testing.T) {
	pal := PaletteFor(TimeNight)
	if !pal.DimFactor {
		t.Error("night palette should have DimFactor=true")
	}
}

func TestPaletteMorningIsNotDimmed(t *testing.T) {
	pal := PaletteFor(TimeMorning)
	if pal.DimFactor {
		t.Error("morning palette should not be dimmed")
	}
}
