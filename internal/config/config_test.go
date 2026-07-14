package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.FPS != 20 {
		t.Errorf("default FPS = %d, want 20", cfg.FPS)
	}
	if cfg.Speed != 1.0 {
		t.Errorf("default Speed = %f, want 1.0", cfg.Speed)
	}
	if cfg.Train != "random" {
		t.Errorf("default Train = %q, want 'random'", cfg.Train)
	}
	if cfg.Scene != "random" {
		t.Errorf("default Scene = %q, want 'random'", cfg.Scene)
	}
	if cfg.Weather != "random" {
		t.Errorf("default Weather = %q, want 'random'", cfg.Weather)
	}
	if cfg.Time != "auto" {
		t.Errorf("default Time = %q, want 'auto'", cfg.Time)
	}
	if cfg.Color != "auto" {
		t.Errorf("default Color = %q, want 'auto'", cfg.Color)
	}
	if cfg.NoColor {
		t.Error("default NoColor should be false")
	}
	if cfg.ReducedMotion {
		t.Error("default ReducedMotion should be false")
	}
	if !cfg.RandomEvents {
		t.Error("default RandomEvents should be true")
	}
}

func TestValidateValidConfig(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid: %v", err)
	}
}

func TestValidateInvalidTrain(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Train = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid train should fail validation")
	}
}

func TestValidateInvalidScene(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Scene = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid scene should fail validation")
	}
}

func TestValidateInvalidWeather(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Weather = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid weather should fail validation")
	}
}

func TestValidateInvalidTime(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Time = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid time should fail validation")
	}
}

func TestValidateInvalidColor(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Color = "invalid"
	if err := cfg.Validate(); err == nil {
		t.Error("invalid color should fail validation")
	}
}

func TestValidateFPSTooLow(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FPS = 4
	if err := cfg.Validate(); err == nil {
		t.Error("FPS 4 should fail validation")
	}
}

func TestValidateFPSTooHigh(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FPS = 61
	if err := cfg.Validate(); err == nil {
		t.Error("FPS 61 should fail validation")
	}
}

func TestValidateSpeedTooLow(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Speed = 0.24
	if err := cfg.Validate(); err == nil {
		t.Error("speed 0.24 should fail validation")
	}
}

func TestValidateSpeedTooHigh(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Speed = 3.1
	if err := cfg.Validate(); err == nil {
		t.Error("speed 3.1 should fail validation")
	}
}

func TestValidateBoundaryValues(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FPS = 5
	cfg.Speed = 0.25
	if err := cfg.Validate(); err != nil {
		t.Errorf("boundary FPS=5 Speed=0.25 should be valid: %v", err)
	}
	cfg.FPS = 60
	cfg.Speed = 3.0
	if err := cfg.Validate(); err != nil {
		t.Errorf("boundary FPS=60 Speed=3.0 should be valid: %v", err)
	}
}

func TestLoadMissingConfigFile(t *testing.T) {
	// A nonexistent config path should not error.
	result, err := Load("/nonexistent/path/config.toml", nil)
	if err != nil {
		t.Errorf("missing config file should not error: %v", err)
	}
	if result.Config.FPS != 20 {
		t.Errorf("defaults should be used: FPS = %d, want 20", result.Config.FPS)
	}
}

func TestLoadWithConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `train = "steam"
scene = "station"
weather = "rain"
time = "night"
fps = 30
speed = 1.5
color = "256"
no_color = true
reduced_motion = true
random_events = false
seed = 42
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Load(path, nil)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	cfg := result.Config
	if cfg.Train != "steam" {
		t.Errorf("Train = %q, want 'steam'", cfg.Train)
	}
	if cfg.Scene != "station" {
		t.Errorf("Scene = %q, want 'station'", cfg.Scene)
	}
	if cfg.Weather != "rain" {
		t.Errorf("Weather = %q, want 'rain'", cfg.Weather)
	}
	if cfg.Time != "night" {
		t.Errorf("Time = %q, want 'night'", cfg.Time)
	}
	if cfg.FPS != 30 {
		t.Errorf("FPS = %d, want 30", cfg.FPS)
	}
	if cfg.Speed != 1.5 {
		t.Errorf("Speed = %f, want 1.5", cfg.Speed)
	}
	if cfg.Color != "256" {
		t.Errorf("Color = %q, want '256'", cfg.Color)
	}
	if !cfg.NoColor {
		t.Error("NoColor should be true")
	}
	if !cfg.ReducedMotion {
		t.Error("ReducedMotion should be true")
	}
	if cfg.RandomEvents {
		t.Error("RandomEvents should be false")
	}
	if cfg.Seed != 42 {
		t.Errorf("Seed = %d, want 42", cfg.Seed)
	}
}

func TestLoadConfigFileWithUnknownFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `train = "diesel"
unknown_field = "should be ignored"
another_unknown = 123
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Load(path, nil)
	if err != nil {
		t.Errorf("unknown fields should not cause error: %v", err)
	}
	if result.Config.Train != "diesel" {
		t.Errorf("Train = %q, want 'diesel'", result.Config.Train)
	}
}

func TestLoadConfigFileEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Load(path, nil)
	if err != nil {
		t.Errorf("empty config file should not error: %v", err)
	}
	// Should use defaults.
	if result.Config.FPS != 20 {
		t.Errorf("empty config FPS = %d, want 20 (default)", result.Config.FPS)
	}
}

func TestLoadCLIOverridesConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `train = "steam"
fps = 30
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	trainCLI := "diesel"
	fpsCLI := 40
	cliFlags := &CLIFlags{
		Train: &trainCLI,
		FPS:   &fpsCLI,
	}

	result, err := Load(path, cliFlags)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	// CLI should override config file.
	if result.Config.Train != "diesel" {
		t.Errorf("Train = %q, want 'diesel' (CLI override)", result.Config.Train)
	}
	if result.Config.FPS != 40 {
		t.Errorf("FPS = %d, want 40 (CLI override)", result.Config.FPS)
	}
}

func TestLoadDefaultsWhenNoConfigAndNoCLI(t *testing.T) {
	result, err := Load("/nonexistent/path.toml", nil)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if result.Config.Train != "random" {
		t.Errorf("Train = %q, want 'random' (default)", result.Config.Train)
	}
	if result.Config.FPS != 20 {
		t.Errorf("FPS = %d, want 20 (default)", result.Config.FPS)
	}
}

func TestLoadCLIOverridesDefaults(t *testing.T) {
	trainCLI := "electric"
	cliFlags := &CLIFlags{
		Train: &trainCLI,
	}

	result, err := Load("/nonexistent/path.toml", cliFlags)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if result.Config.Train != "electric" {
		t.Errorf("Train = %q, want 'electric' (CLI override)", result.Config.Train)
	}
}

func TestLoadInvalidConfigFileValues(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `train = "invalid_type"
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path, nil)
	if err == nil {
		t.Error("invalid config file values should fail validation")
	}
}

func TestCLIFlagsApplyTo(t *testing.T) {
	cfg := DefaultConfig()
	train := "electric"
	weather := "snow"
	noColor := true
	cliFlags := &CLIFlags{
		Train:   &train,
		Weather: &weather,
		NoColor: &noColor,
	}
	cliFlags.applyTo(&cfg)
	if cfg.Train != "electric" {
		t.Errorf("Train = %q, want 'electric'", cfg.Train)
	}
	if cfg.Weather != "snow" {
		t.Errorf("Weather = %q, want 'snow'", cfg.Weather)
	}
	if !cfg.NoColor {
		t.Error("NoColor should be true")
	}
}

func TestCLIFlagsNilDoesNotOverride(t *testing.T) {
	cfg := DefaultConfig()
	originalTrain := cfg.Train
	cliFlags := &CLIFlags{
		Train: nil,
	}
	cliFlags.applyTo(&cfg)
	if cfg.Train != originalTrain {
		t.Errorf("nil CLI flag should not override: Train = %q, want %q", cfg.Train, originalTrain)
	}
}

func TestConfigString(t *testing.T) {
	cfg := DefaultConfig()
	s := cfg.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Error("DefaultConfigPath should not be empty")
	}
}
