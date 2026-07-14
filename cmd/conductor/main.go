// Command conductor launches an ambient railway animation in the terminal.
//
// Usage:
//
//	conductor [flags]
//
// Flags:
//
//	--train <steam|diesel|electric|random>
//	--scene <countryside|station|mountain|random>
//	--weather <clear|rain|snow|random>
//	--time <morning|day|sunset|night|auto>
//	--speed <float>
//	--fps <integer>
//	--seed <integer>
//	--color <auto|16|256|truecolor>
//	--no-color
//	--reduced-motion
//	--config <path>
//	--version
//	--help
//
// Configuration precedence: built-in defaults → config file → CLI flags.
package main

import (
	"fmt"
	"os"

	"github.com/example/conductor/internal/app"
	"github.com/example/conductor/internal/config"
	"github.com/spf13/cobra"
)

// version is the application version. It is overridden at build time with
// -ldflags "-X main.version=...".
var version = "0.1.0-dev"

func main() {
	var (
		train         string
		scene         string
		weather       string
		timePeriod    string
		speed         float64
		fps           int
		seed          int64
		color         string
		noColor       bool
		reducedMotion bool
		configPath    string
		showVersion   bool
	)

	rootCmd := &cobra.Command{
		Use:   "conductor",
		Short: "A tiny railway journey inside your terminal",
		Long: "Conductor is an ambient terminal animation that displays an " +
			"animated railway scene with trains, scenery, weather, and " +
			"occasional random events.\n\n" +
			"Run it with no arguments to start the animation:\n\n" +
			"  conductor\n\n" +
			"Press 'h' or '?' inside the animation for keyboard controls.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("conductor %s\n", version)
				return nil
			}

			// Build CLI flags struct from parsed values.
			cliFlags := &config.CLIFlags{
				Train:         strPtr(train),
				Scene:         strPtr(scene),
				Weather:       strPtr(weather),
				Time:          strPtr(timePeriod),
				FPS:           intPtr(fps),
				Speed:         floatPtr(speed),
				Color:         strPtr(color),
				NoColor:       boolPtr(noColor),
				ReducedMotion: boolPtr(reducedMotion),
				Seed:          int64Ptr(seed),
				ConfigPath:    strPtr(configPath),
			}

			// Load configuration with precedence.
			result, err := config.Load(configPath, cliFlags)
			if err != nil {
				return err
			}

			// Create and run the application.
			application := app.New(result.Config)
			if err := application.Run(); err != nil {
				return err
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Register flags.
	rootCmd.Flags().StringVar(&train, "train", "", "train type: steam, diesel, electric, or random")
	rootCmd.Flags().StringVar(&scene, "scene", "", "scene: countryside, station, mountain, or random")
	rootCmd.Flags().StringVar(&weather, "weather", "", "weather: clear, rain, snow, or random")
	rootCmd.Flags().StringVar(&timePeriod, "time", "", "time of day: morning, day, sunset, night, or auto")
	rootCmd.Flags().Float64Var(&speed, "speed", 0, "animation speed multiplier (0.25-3.0)")
	rootCmd.Flags().IntVar(&fps, "fps", 0, "target frame rate (5-60)")
	rootCmd.Flags().Int64Var(&seed, "seed", 0, "random seed for deterministic scenes")
	rootCmd.Flags().StringVar(&color, "color", "", "color mode: auto, 16, 256, or truecolor")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "disable all colors")
	rootCmd.Flags().BoolVar(&reducedMotion, "reduced-motion", false, "reduce particles and frame rate")
	rootCmd.Flags().StringVar(&configPath, "config", "", "path to configuration file")
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "print version and exit")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "conductor: %v\n", err)
		os.Exit(1)
	}
}

// strPtr returns a pointer to the given string, or nil if it's empty.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// intPtr returns a pointer to the given int, or nil if it's zero.
func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// floatPtr returns a pointer to the given float, or nil if it's zero.
func floatPtr(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

// boolPtr returns a pointer to the given bool. Only non-nil true values
// override config.
func boolPtr(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

// int64Ptr returns a pointer to the given int64, or nil if it's zero.
func int64Ptr(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}
