package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// defaultConfigPath returns the platform-appropriate default config file path.
func defaultConfigPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		return filepath.Join(appData, "conductor", "config.toml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.toml"
	}
	return filepath.Join(home, ".config", "conductor", "config.toml")
}

// validTrainTypes is the set of accepted --train values.
var validTrainTypes = map[string]bool{
	"steam": true, "diesel": true, "electric": true, "random": true,
}

// validSceneTypes is the set of accepted --scene values.
var validSceneTypes = map[string]bool{
	"countryside": true, "station": true, "mountain": true, "random": true,
}

// validWeatherTypes is the set of accepted --weather values.
var validWeatherTypes = map[string]bool{
	"clear": true, "rain": true, "snow": true, "random": true,
}

// validTimePeriods is the set of accepted --time values.
var validTimePeriods = map[string]bool{
	"morning": true, "day": true, "sunset": true, "night": true, "auto": true,
}

// validColorModes is the set of accepted --color values.
var validColorModes = map[string]bool{
	"auto": true, "16": true, "256": true, "truecolor": true,
}

// isValid checks a string against a map of valid values.
func isValid(s string, valid map[string]bool) bool {
	return valid[s]
}

// Validate checks the config for invalid values and returns an error with a
// clear, user-facing message if any field is out of range. This function does
// not return stack traces or internal errors.
func (c *Config) Validate() error {
	if !isValid(c.Train, validTrainTypes) {
		return fmt.Errorf("invalid --train value %q: expected steam, diesel, electric, or random", c.Train)
	}
	if !isValid(c.Scene, validSceneTypes) {
		return fmt.Errorf("invalid --scene value %q: expected countryside, station, mountain, or random", c.Scene)
	}
	if !isValid(c.Weather, validWeatherTypes) {
		return fmt.Errorf("invalid --weather value %q: expected clear, rain, snow, or random", c.Weather)
	}
	if !isValid(c.Time, validTimePeriods) {
		return fmt.Errorf("invalid --time value %q: expected morning, day, sunset, night, or auto", c.Time)
	}
	if !isValid(c.Color, validColorModes) {
		return fmt.Errorf("invalid --color value %q: expected auto, 16, 256, or truecolor", c.Color)
	}
	if c.FPS < 5 || c.FPS > 60 {
		return fmt.Errorf("invalid fps %d: must be between 5 and 60", c.FPS)
	}
	if c.Speed < 0.25 || c.Speed > 3.0 {
		return fmt.Errorf("invalid speed %.2f: must be between 0.25 and 3.0", c.Speed)
	}
	return nil
}
