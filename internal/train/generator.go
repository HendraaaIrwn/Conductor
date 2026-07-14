package train

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// ConsistConfig controls how a train consist is generated.
type ConsistConfig struct {
	TrainType      TrainType
	Direction      Direction
	MinCars        int  // minimum carriage count (excluding locomotive and end car)
	MaxCars        int  // maximum carriage count (excluding locomotive and end car)
	MaxWidth       int  // maximum total train width in cells (usually terminal width)
	AllowLongTrain bool // if true, may exceed MaxWidth (rare event)
}

// DefaultConsistConfig returns a sensible default configuration for the given
// train type and terminal width.
func DefaultConsistConfig(t TrainType, dir Direction, terminalWidth int) ConsistConfig {
	return ConsistConfig{
		TrainType:      t,
		Direction:      dir,
		MinCars:        3,
		MaxCars:        8,
		MaxWidth:       terminalWidth,
		AllowLongTrain: false,
	}
}

// Generator creates train consists procedurally. It uses a dedicated RNG so
// that consist generation is deterministic given a seed.
type Generator struct {
	rng   *rand.Rand
	style tcell.Style
}

// NewGenerator creates a Generator with the given seed and base style.
func NewGenerator(seed int64, style tcell.Style) *Generator {
	return &Generator{
		rng:   rand.New(rand.NewSource(seed)),
		style: style,
	}
}

// Generate produces a train consist from the given configuration. The
// locomotive is always at offset 0; carriages are placed behind it at
// negative offsets. If the train is too long, carriages are removed from the
// rear until it fits within MaxWidth (unless AllowLongTrain is set).
func (g *Generator) Generate(cfg ConsistConfig) []Car {
	locoSprite := LocomotiveSprite(cfg.TrainType, cfg.Direction, g.style)
	if locoSprite == nil {
		return nil
	}

	// Choose a random number of carriages.
	count := cfg.MinCars
	if cfg.MaxCars > cfg.MinCars {
		count = cfg.MinCars + g.rng.Intn(cfg.MaxCars-cfg.MinCars+1)
	}

	// Select compatible carriages.
	compat := CompatibleCarriages(cfg.TrainType)
	if len(compat) == 0 {
		return []Car{{Sprite: locoSprite, OffsetX: 0}}
	}

	cars := []Car{{Sprite: locoSprite, OffsetX: 0}}
	carCategories := make([]CarriageCategory, 0, count)
	for i := 0; i < count; i++ {
		cat := compat[g.rng.Intn(len(compat))]
		carCategories = append(carCategories, cat)
	}

	// Add end car if configured.
	endCat := EndCarCategory(cfg.TrainType)
	hasEndCar := endCat >= 0 && endCat != cfg.TrainType.toEndCarSelf()
	if hasEndCar {
		carCategories = append(carCategories, endCat)
	}

	// Build car list with offsets.
	gap := 1
	cursor := 0
	for _, cat := range carCategories {
		sprite := CarriageSprite(cat, cfg.Direction, g.style)
		if sprite == nil {
			continue
		}
		offset := cursor - (sprite.Width + gap)
		cars = append(cars, Car{Sprite: sprite, OffsetX: offset})
		cursor = offset
	}

	// Trim if too long (unless long train is allowed).
	if !cfg.AllowLongTrain {
		cars = trimToFit(cars, cfg.MaxWidth)
	}

	return cars
}

// trimToFit removes carriages from the rear (end) of the train until the total
// width fits within maxWidth. The locomotive is always kept.
func trimToFit(cars []Car, maxWidth int) []Car {
	if len(cars) <= 1 || maxWidth <= 0 {
		return cars
	}
	for TrainWidth(cars) > maxWidth && len(cars) > 1 {
		cars = cars[:len(cars)-1]
	}
	return cars
}

// toEndCarSelf is a helper that returns a sentinel value; if the end car
// category equals the train type's own commuter car, we still add it. This
// method exists to handle the electric train case where the "end car" is
// another commuter car rather than a distinct caboose.
func (t TrainType) toEndCarSelf() CarriageCategory {
	return -1
}

// GenerateRandom selects a random train type and generates a consist for it.
func (g *Generator) GenerateRandom(dir Direction, terminalWidth int) (TrainType, []Car) {
	types := AllTrainTypes()
	t := types[g.rng.Intn(len(types))]
	cfg := DefaultConsistConfig(t, dir, terminalWidth)
	return t, g.Generate(cfg)
}

// GenerateLongFreight creates a rare long freight train that may exceed the
// terminal width. Only diesel and steam types are used for long freight.
func (g *Generator) GenerateLongFreight(dir Direction, terminalWidth int) (TrainType, []Car) {
	// Steam or diesel for freight.
	types := []TrainType{TypeSteam, TypeDiesel}
	t := types[g.rng.Intn(len(types))]
	cfg := ConsistConfig{
		TrainType:      t,
		Direction:      dir,
		MinCars:        8,
		MaxCars:        12,
		MaxWidth:       terminalWidth,
		AllowLongTrain: true,
	}
	return t, g.Generate(cfg)
}
