package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

// Config holds all application configuration values. The TOML tags define the
// keys used in the configuration file.
type Config struct {
	Train         string  `toml:"train"`
	Scene         string  `toml:"scene"`
	Weather       string  `toml:"weather"`
	Time          string  `toml:"time"`
	FPS           int     `toml:"fps"`
	Speed         float64 `toml:"speed"`
	Color         string  `toml:"color"`
	NoColor       bool    `toml:"no_color"`
	ReducedMotion bool    `toml:"reduced_motion"`
	ShowStatus    bool    `toml:"show_status"`
	RandomEvents  bool    `toml:"random_events"`
	Seed          int64   `toml:"seed"`
}

// LoadResult holds the final resolved configuration and the path of the config
// file that was loaded (if any).
type LoadResult struct {
	Config     Config
	ConfigPath string
}

// Load resolves the configuration in three steps:
//  1. Start with built-in defaults.
//  2. Overlay values from the TOML config file (if it exists).
//  3. Overlay values from CLI flags (those that are non-empty/non-zero).
//
// Missing config file is not an error. Unknown TOML fields are ignored.
// Invalid values are reported with a clear message.
func Load(configPath string, cliFlags *CLIFlags) (*LoadResult, error) {
	cfg := DefaultConfig()

	// If no explicit config path, use the default.
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	// Try to read the config file. Missing file is not an error.
	if _, err := os.Stat(configPath); err == nil {
		if err := loadTOML(configPath, &cfg); err != nil {
			return nil, fmt.Errorf("config file %s: %w", configPath, err)
		}
	}

	// Apply CLI overrides.
	if cliFlags != nil {
		cliFlags.applyTo(&cfg)
	}

	// Validate the final configuration.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &LoadResult{
		Config:     cfg,
		ConfigPath: configPath,
	}, nil
}

// loadTOML reads and parses a TOML config file into the given config struct.
// Unknown fields are ignored (BurntSushi/toml ignores them by default unless
// the Unmarshaler interface is implemented).
func loadTOML(path string, cfg *Config) error {
	_, err := toml.DecodeFile(path, cfg)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	return nil
}

// CLIFlags holds the CLI flag values. Fields that are nil or empty were not
// set on the command line and should not override config file values.
type CLIFlags struct {
	Train         *string
	Scene         *string
	Weather       *string
	Time          *string
	FPS           *int
	Speed         *float64
	Color         *string
	NoColor       *bool
	ReducedMotion *bool
	ConfigPath    *string
	Seed          *int64
}

// applyTo overlays non-nil CLI flag values onto the config struct.
func (f *CLIFlags) applyTo(cfg *Config) {
	if f.Train != nil && *f.Train != "" {
		cfg.Train = *f.Train
	}
	if f.Scene != nil && *f.Scene != "" {
		cfg.Scene = *f.Scene
	}
	if f.Weather != nil && *f.Weather != "" {
		cfg.Weather = *f.Weather
	}
	if f.Time != nil && *f.Time != "" {
		cfg.Time = *f.Time
	}
	if f.FPS != nil && *f.FPS != 0 {
		cfg.FPS = *f.FPS
	}
	if f.Speed != nil && *f.Speed != 0 {
		cfg.Speed = *f.Speed
	}
	if f.Color != nil && *f.Color != "" {
		cfg.Color = *f.Color
	}
	if f.NoColor != nil && *f.NoColor {
		cfg.NoColor = *f.NoColor
	}
	if f.ReducedMotion != nil && *f.ReducedMotion {
		cfg.ReducedMotion = *f.ReducedMotion
	}
	if f.Seed != nil && *f.Seed != 0 {
		cfg.Seed = *f.Seed
	}
}

// String returns a human-readable representation of the config for debugging.
func (c *Config) String() string {
	return fmt.Sprintf("train=%s scene=%s weather=%s time=%s fps=%d speed=%.2f color=%s no-color=%v reduced-motion=%v seed=%d",
		c.Train, c.Scene, c.Weather, c.Time, c.FPS, c.Speed, c.Color, c.NoColor, c.ReducedMotion, c.Seed)
}

// SpeedString returns the speed as a clean string for display.
func (c *Config) SpeedString() string {
	return strconv.FormatFloat(c.Speed, 'f', 2, 64)
}
