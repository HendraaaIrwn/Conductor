package train

// TrainType identifies a locomotive family. It determines which locomotive
// sprite is used and which carriage categories are compatible.
type TrainType int

const (
	// TypeSteam is a classic steam locomotive with a chimney and smoke.
	TypeSteam TrainType = iota
	// TypeDiesel is a modern diesel locomotive with a long body.
	TypeDiesel
	// TypeElectric is a sleek electric commuter train with a pantograph.
	TypeElectric
)

// String returns the human-readable name of the train type.
func (t TrainType) String() string {
	switch t {
	case TypeSteam:
		return "steam"
	case TypeDiesel:
		return "diesel"
	case TypeElectric:
		return "electric"
	default:
		return "unknown"
	}
}

// ParseTrainType converts a string to a TrainType. Returns false if the
// string does not match a known type.
func ParseTrainType(s string) (TrainType, bool) {
	switch s {
	case "steam":
		return TypeSteam, true
	case "diesel":
		return TypeDiesel, true
	case "electric":
		return TypeElectric, true
	default:
		return 0, false
	}
}

// CarriageCategory classifies a carriage for compatibility selection. Each
// locomotive type accepts certain categories in its consist.
type CarriageCategory int

const (
	// CatPassenger is a passenger carriage with windows.
	CatPassenger CarriageCategory = iota
	// CatBoxcar is an enclosed freight car with a sliding door.
	CatBoxcar
	// CatTank is a cylindrical tank car on a flatbed.
	CatTank
	// CatOpenCargo is a low-sided open freight car.
	CatOpenCargo
	// CatCaboose is an end-of-train cabin car.
	CatCaboose
	// CatCommuter is a streamlined commuter carriage for electric trains.
	CatCommuter
)

// String returns the human-readable name of the carriage category.
func (c CarriageCategory) String() string {
	switch c {
	case CatPassenger:
		return "passenger"
	case CatBoxcar:
		return "boxcar"
	case CatTank:
		return "tank"
	case CatOpenCargo:
		return "open-cargo"
	case CatCaboose:
		return "caboose"
	case CatCommuter:
		return "commuter"
	default:
		return "unknown"
	}
}

// AllTrainTypes returns all train types for random selection.
func AllTrainTypes() []TrainType {
	return []TrainType{TypeSteam, TypeDiesel, TypeElectric}
}
