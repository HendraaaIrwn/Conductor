package train

import (
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// buildSprite constructs a render.Sprite from frame art strings. Each string
// in the frames slice is joined with newlines and parsed into a CellGrid.
func buildSprite(name string, style tcell.Style, frames ...[]string) *render.Sprite {
	grids := make([]render.CellGrid, len(frames))
	for i, art := range frames {
		grids[i] = render.ParseFrame(joinLines(art), style)
	}
	return render.NewSprite(name, grids...)
}

// joinLines concatenates lines with newline separators.
func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

// ValidateConsist checks that all cars in a consist have compatible heights
// (coupling alignment) and that their sprites are valid. Returns nil if the
// consist is well-formed.
func ValidateConsist(cars []Car) error {
	if len(cars) == 0 {
		return nil
	}
	for i, car := range cars {
		if car.Sprite == nil {
			return &consistError{car: i, msg: "nil sprite"}
		}
		if err := car.Sprite.Validate(); err != nil {
			return &consistError{car: i, msg: err.Error()}
		}
	}
	for i := 1; i < len(cars); i++ {
		if err := render.ValidateCoupling(cars[0].Sprite, cars[i].Sprite); err != nil {
			return &consistError{car: i, msg: err.Error()}
		}
	}
	return nil
}

type consistError struct {
	car int
	msg string
}

func (e *consistError) Error() string {
	return "consist car " + itoa(e.car) + ": " + e.msg
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// SimpleConsist builds a minimal train consist for testing: one steam
// locomotive followed by one passenger carriage. Production code should use
// the Generator in generator.go instead.
func SimpleConsist(style tcell.Style) []Car {
	loco := SteamLocomotiveRight(style)
	carriage := PassengerCarRight(style)
	gap := 1
	carriageOffset := -(carriage.Width + gap)
	return []Car{
		{Sprite: loco, OffsetX: 0},
		{Sprite: carriage, OffsetX: carriageOffset},
	}
}
