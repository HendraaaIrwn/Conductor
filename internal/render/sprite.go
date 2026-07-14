package render

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// CellGrid is a 2D grid of cells indexed [y][x].
type CellGrid [][]Cell

// Sprite is a small multi-frame image made of terminal cells. All frames
// within a sprite share the same width and height.
type Sprite struct {
	Name   string
	Width  int
	Height int
	Frames []CellGrid
}

// ParseFrame converts a multiline string into a CellGrid. Each line becomes a
// row and each rune becomes a cell. Lines shorter than the longest line are
// padded with blank cells so that every row has equal width. The style is
// applied to every non-blank cell.
func ParseFrame(art string, style tcell.Style) CellGrid {
	lines := strings.Split(art, "\n")
	height := len(lines)
	width := 0
	for _, line := range lines {
		if w := runeWidth(line); w > width {
			width = w
		}
	}
	grid := make(CellGrid, height)
	for y, line := range lines {
		row := make([]Cell, width)
		for x := 0; x < width; x++ {
			row[x] = Blank()
		}
		x := 0
		for _, r := range line {
			if r == ' ' {
				x++
				continue
			}
			if x < width {
				row[x] = Cell{Rune: r, Style: style}
			}
			x++
		}
		grid[y] = row
	}
	return grid
}

// runeWidth counts the number of rune cells a string will occupy.
func runeWidth(s string) int {
	return len([]rune(s))
}

// NewSprite builds a Sprite from one or more frame grids. All frames must have
// identical dimensions.
func NewSprite(name string, frames ...CellGrid) *Sprite {
	if len(frames) == 0 {
		return &Sprite{Name: name}
	}
	height := len(frames[0])
	width := 0
	if height > 0 {
		width = len(frames[0][0])
	}
	return &Sprite{
		Name:   name,
		Width:  width,
		Height: height,
		Frames: frames,
	}
}

// FrameCount returns the number of animation frames in the sprite.
func (s *Sprite) FrameCount() int { return len(s.Frames) }

// Draw blits frame index f onto the canvas at position (x, y). The position is
// the top-left corner of the sprite. Transparent (blank) cells in the sprite
// are skipped so that the background shows through.
func (s *Sprite) Draw(canvas *Canvas, x, y, frame int) {
	if frame < 0 || frame >= len(s.Frames) {
		return
	}
	grid := s.Frames[frame]
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			cell := grid[row][col]
			if cell.IsBlank() {
				continue
			}
			canvas.Set(x+col, y+row, cell)
		}
	}
}

// Validate checks the sprite for common structural problems. It returns nil if
// the sprite is well-formed, or an error describing the first problem found.
//
// Validated properties:
//   - At least one frame exists.
//   - Width and height are positive.
//   - All frames have identical dimensions.
//   - No frame is entirely blank (empty).
func (s *Sprite) Validate() error {
	if len(s.Frames) == 0 {
		return errSpriteEmpty(s.Name, "no frames")
	}
	if s.Width <= 0 {
		return errSpriteEmpty(s.Name, "width is zero")
	}
	if s.Height <= 0 {
		return errSpriteEmpty(s.Name, "height is zero")
	}
	for i, grid := range s.Frames {
		if len(grid) != s.Height {
			return errFrameDim(s.Name, i, len(grid), s.Height, "height")
		}
		for _, row := range grid {
			if len(row) != s.Width {
				return errFrameDim(s.Name, i, len(row), s.Width, "width")
			}
		}
		if isGridBlank(grid) {
			return errFrameBlank(s.Name, i)
		}
	}
	return nil
}

// isGridBlank reports whether every cell in the grid is blank.
func isGridBlank(grid CellGrid) bool {
	for _, row := range grid {
		for _, cell := range row {
			if !cell.IsBlank() {
				return false
			}
		}
	}
	return true
}

// ValidateCoupling checks that two sprites have compatible coupling heights:
// their heights must be equal so that the rail line aligns when they are
// placed on the same track Y.
func ValidateCoupling(a, b *Sprite) error {
	if a.Height != b.Height {
		return errCouplingHeight(a.Name, b.Name, a.Height, b.Height)
	}
	return nil
}

type spriteError struct {
	sprite  string
	message string
}

func (e *spriteError) Error() string {
	return "sprite " + e.sprite + ": " + e.message
}

func errSpriteEmpty(name, reason string) error {
	return &spriteError{sprite: name, message: reason}
}

func errFrameDim(name string, frame, got, want int, dim string) error {
	return &spriteError{
		sprite: name,
		message: "frame " + itoa(frame) + " " + dim +
			" is " + itoa(got) + ", expected " + itoa(want),
	}
}

func errFrameBlank(name string, frame int) error {
	return &spriteError{
		sprite:  name,
		message: "frame " + itoa(frame) + " is entirely blank",
	}
}

func errCouplingHeight(a, b string, ha, hb int) error {
	return &spriteError{
		sprite: a + " <-> " + b,
		message: "coupling height mismatch: " + itoa(ha) +
			" vs " + itoa(hb),
	}
}

// itoa is a minimal int-to-string to avoid importing strconv in this hot path.
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
