package train

import (
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// SpriteFactory returns a sprite for a given style. This indirection allows
// sprites to be lazily created with the correct palette/style at consist
// generation time.
type SpriteFactory func(style tcell.Style) *render.Sprite

// carriageSpec describes a carriage type and its sprite factories for both
// directions.
type carriageSpec struct {
	Category CarriageCategory
	Right    SpriteFactory
	Left     SpriteFactory
}

// locomotiveSpec describes a locomotive type, its sprite factories, and the
// carriage categories it is compatible with.
type locomotiveSpec struct {
	Type           TrainType
	Right          SpriteFactory
	Left           SpriteFactory
	CompatibleCars []CarriageCategory
	EndCar         CarriageCategory // optional, set to -1 for none
	Speed          float64          // base speed in cells per second
	EmitsSmoke     bool             // steam locomotives emit smoke
}

// registry holds all locomotive and carriage definitions. It is populated at
// init time and read by the consist generator.
var registry = struct {
	locomotives  map[TrainType]*locomotiveSpec
	carriages    map[CarriageCategory]*carriageSpec
	carriageList []CarriageCategory // for deterministic iteration
}{
	locomotives: map[TrainType]*locomotiveSpec{},
	carriages:   map[CarriageCategory]*carriageSpec{},
}

func init() {
	registerLocomotive(&locomotiveSpec{
		Type:           TypeSteam,
		Right:          SteamLocomotiveRight,
		Left:           SteamLocomotiveLeft,
		CompatibleCars: []CarriageCategory{CatPassenger, CatBoxcar, CatTank, CatOpenCargo},
		EndCar:         CatCaboose,
		Speed:          22.0,
		EmitsSmoke:     true,
	})
	registerLocomotive(&locomotiveSpec{
		Type:           TypeDiesel,
		Right:          DieselLocomotiveRight,
		Left:           DieselLocomotiveLeft,
		CompatibleCars: []CarriageCategory{CatBoxcar, CatTank, CatOpenCargo, CatPassenger},
		EndCar:         CatCaboose,
		Speed:          28.0,
		EmitsSmoke:     false,
	})
	registerLocomotive(&locomotiveSpec{
		Type:           TypeElectric,
		Right:          ElectricLocomotiveRight,
		Left:           ElectricLocomotiveLeft,
		CompatibleCars: []CarriageCategory{CatCommuter},
		EndCar:         CatCommuter, // electric trains end with another commuter car
		Speed:          32.0,
		EmitsSmoke:     false,
	})

	registerCarriage(&carriageSpec{Category: CatPassenger, Right: PassengerCarRight, Left: PassengerCarLeft})
	registerCarriage(&carriageSpec{Category: CatBoxcar, Right: BoxcarRight, Left: BoxcarLeft})
	registerCarriage(&carriageSpec{Category: CatTank, Right: TankCarRight, Left: TankCarLeft})
	registerCarriage(&carriageSpec{Category: CatOpenCargo, Right: OpenCargoRight, Left: OpenCargoLeft})
	registerCarriage(&carriageSpec{Category: CatCaboose, Right: CabooseRight, Left: CabooseLeft})
	registerCarriage(&carriageSpec{Category: CatCommuter, Right: CommuterCarRight, Left: CommuterCarLeft})
}

func registerLocomotive(spec *locomotiveSpec) {
	registry.locomotives[spec.Type] = spec
}

func registerCarriage(spec *carriageSpec) {
	registry.carriages[spec.Category] = spec
	registry.carriageList = append(registry.carriageList, spec.Category)
}

// LocomotiveSpec returns the spec for the given train type, or nil if the type
// is not registered.
func LocomotiveSpec(t TrainType) *locomotiveSpec {
	return registry.locomotives[t]
}

// CarriageSpec returns the spec for the given carriage category, or nil if the
// category is not registered.
func CarriageSpec(c CarriageCategory) *carriageSpec {
	return registry.carriages[c]
}

// LocomotiveSprite returns the appropriate sprite for a locomotive given the
// direction.
func LocomotiveSprite(t TrainType, dir Direction, style tcell.Style) *render.Sprite {
	spec := LocomotiveSpec(t)
	if spec == nil {
		return nil
	}
	if dir == RightToLeft {
		return spec.Left(style)
	}
	return spec.Right(style)
}

// CarriageSprite returns the appropriate sprite for a carriage category given
// the direction.
func CarriageSprite(c CarriageCategory, dir Direction, style tcell.Style) *render.Sprite {
	spec := CarriageSpec(c)
	if spec == nil {
		return nil
	}
	if dir == RightToLeft {
		return spec.Left(style)
	}
	return spec.Right(style)
}

// CompatibleCarriages returns the carriage categories compatible with the
// given train type.
func CompatibleCarriages(t TrainType) []CarriageCategory {
	spec := LocomotiveSpec(t)
	if spec == nil {
		return nil
	}
	return spec.CompatibleCars
}

// EndCarCategory returns the end carriage category for the given train type,
// or -1 if no end car is used.
func EndCarCategory(t TrainType) CarriageCategory {
	spec := LocomotiveSpec(t)
	if spec == nil {
		return -1
	}
	return spec.EndCar
}

// TrainSpeed returns the base speed for the given train type.
func TrainSpeed(t TrainType) float64 {
	spec := LocomotiveSpec(t)
	if spec == nil {
		return 25.0
	}
	return spec.Speed
}

// EmitsSmoke reports whether the given train type emits smoke.
func EmitsSmoke(t TrainType) bool {
	spec := LocomotiveSpec(t)
	if spec == nil {
		return false
	}
	return spec.EmitsSmoke
}
