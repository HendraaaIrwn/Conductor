// Package render provides the canvas abstraction used by the animation
// engine. The canvas is a double-buffered grid of cells: every frame the
// simulation writes into the current buffer, and the renderer only flushes the
// cells that differ from the previous frame to the terminal.
package render

import (
	"github.com/example/conductor/internal/terminal"
	"github.com/gdamore/tcell/v2"
)

// Canvas is a double-buffered grid of cells. It owns two equally sized
// buffers: the current frame being painted and the previous frame that was
// last flushed to the terminal. Flush only writes cells whose contents
// changed between the two buffers.
type Canvas struct {
	width, height int

	// current is the buffer the simulation paints into during a frame.
	current [][]Cell
	// previous is the last frame that was flushed to the terminal.
	previous [][]Cell
}

// NewCanvas creates a canvas with the given dimensions. Both buffers are
// initialised to blank cells.
func NewCanvas(width, height int) *Canvas {
	c := &Canvas{
		width:  width,
		height: height,
	}
	c.current = newBuffer(width, height)
	c.previous = newBuffer(width, height)
	return c
}

// Width returns the canvas width in cells.
func (c *Canvas) Width() int { return c.width }

// Height returns the canvas height in cells.
func (c *Canvas) Height() int { return c.height }

// Clear resets the current buffer to blank cells. The previous buffer is left
// untouched so that the next Flush can still diff against the last rendered
// frame.
func (c *Canvas) Clear() {
	fillBlank(c.current)
}

// Set writes a cell at the given coordinate. Out-of-bounds writes are ignored.
func (c *Canvas) Set(x, y int, cell Cell) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.current[y][x] = cell
}

// SetRune is a convenience wrapper that sets a cell from a rune and a style.
func (c *Canvas) SetRune(x, y int, r rune, style tcell.Style) {
	c.Set(x, y, Cell{Rune: r, Style: style})
}

// CellAt returns the cell at the given coordinate from the current buffer.
// Out-of-bounds coordinates return a blank cell.
func (c *Canvas) CellAt(x, y int) Cell {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return Blank()
	}
	return c.current[y][x]
}

// FillRect fills a rectangular region with the given cell.
func (c *Canvas) FillRect(x, y, w, h int, cell Cell) {
	for row := y; row < y+h; row++ {
		for col := x; col < x+w; col++ {
			c.Set(col, row, cell)
		}
	}
}

// Resize changes the canvas dimensions. The current buffer is rebuilt blank and
// the previous buffer is rebuilt blank as well, which forces a full repaint on
// the next Flush. This avoids stale-cell artefacts after a terminal resize.
func (c *Canvas) Resize(width, height int) {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	c.width = width
	c.height = height
	c.current = newBuffer(width, height)
	c.previous = newBuffer(width, height)
}

// Flush writes only the cells that changed since the last Flush to the
// terminal screen, then calls Show to make the changes visible. After Flush,
// the current buffer becomes the previous buffer for the next frame.
func (c *Canvas) Flush(screen *terminal.Screen) {
	minW := c.width
	minH := c.height
	// Defensive: the two buffers should always match the canvas size, but
	// guard against any drift caused by a resize racing with a flush.
	if len(c.previous) < minH {
		minH = len(c.previous)
	}
	for y := 0; y < minH; y++ {
		prevRow := c.previous[y]
		curRow := c.current[y]
		rowW := minW
		if len(prevRow) < rowW {
			rowW = len(prevRow)
		}
		if len(curRow) < rowW {
			rowW = len(curRow)
		}
		for x := 0; x < rowW; x++ {
			cur := curRow[x]
			if !cur.Equal(prevRow[x]) {
				screen.SetCell(x, y, cur.Rune, cur.Style)
			}
		}
	}
	screen.Show()

	// Swap: copy current into previous for the next diff.
	for y := 0; y < minH && y < len(c.current) && y < len(c.previous); y++ {
		copy(c.previous[y], c.current[y])
	}
}

// newBuffer allocates a 2D slice of blank cells with the given dimensions.
func newBuffer(width, height int) [][]Cell {
	buf := make([][]Cell, height)
	for y := range buf {
		buf[y] = make([]Cell, width)
		for x := range buf[y] {
			buf[y][x] = Blank()
		}
	}
	return buf
}

// fillBlank resets every cell in the buffer to a blank cell.
func fillBlank(buf [][]Cell) {
	for y := range buf {
		for x := range buf[y] {
			buf[y][x] = Blank()
		}
	}
}
