// Package config defines the application configuration model, default values,
// TOML file parsing, and precedence resolution (defaults → config file → CLI
// flags).
package config

// defaults holds the built-in default configuration values. These are the
// lowest-precedence source; config file values override them, and CLI flags
// override both.
var defaults = Config{
	Train:         "random",
	Scene:         "random",
	Weather:       "random",
	Time:          "auto",
	FPS:           20,
	Speed:         1.0,
	Color:         "auto",
	NoColor:       false,
	ReducedMotion: false,
	ShowStatus:    false,
	RandomEvents:  true,
	Seed:          0,
}

// DefaultConfig returns a copy of the built-in default configuration.
func DefaultConfig() Config {
	return defaults
}

// DefaultConfigPath returns the default configuration file path. On Windows
// this uses the user's AppData directory; on other platforms it uses
// ~/.config/conductor/config.toml.
func DefaultConfigPath() string {
	return defaultConfigPath()
}
