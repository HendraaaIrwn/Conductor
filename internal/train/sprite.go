// Package train contains locomotive and carriage definitions, the train
// consist generator, and the train entity behavior. Sprites are defined in the
// render package; this package focuses on train-specific composition and
// behavior.
package train

// Direction is the travel direction of a train.
type Direction int

const (
	// LeftToRight means the train moves towards the right edge.
	LeftToRight Direction = iota
	// RightToLeft means the train moves towards the left edge.
	RightToLeft
)

// VelocityForDirection returns the sign of the velocity for a direction. A
// left-to-right train has positive velocity; a right-to-left train has
// negative velocity.
func VelocityForDirection(d Direction) float64 {
	if d == RightToLeft {
		return -1
	}
	return 1
}
