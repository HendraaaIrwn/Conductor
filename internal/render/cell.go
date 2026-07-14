package render

import "github.com/gdamore/tcell/v2"

// Cell represents a single terminal cell with a rune and a style. The zero
// value is a blank cell with the default style.
type Cell struct {
	Rune  rune
	Style tcell.Style
}

// IsBlank reports whether the cell is empty (a space rune with the default
// style). Blank cells are skipped when diffing frames so that the renderer
// does not waste time repainting empty space.
func (c Cell) IsBlank() bool {
	return c.Rune == 0 || c.Rune == ' '
}

// Equal reports whether two cells are visually identical.
func (c Cell) Equal(other Cell) bool {
	return c.Rune == other.Rune && c.Style == other.Style
}

// Blank returns a blank cell.
func Blank() Cell {
	return Cell{Rune: ' ', Style: tcell.StyleDefault}
}
