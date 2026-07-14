package render

import "testing"

func TestNewCanvasDimensions(t *testing.T) {
	c := NewCanvas(80, 24)
	if c.Width() != 80 {
		t.Errorf("Width = %d, want 80", c.Width())
	}
	if c.Height() != 24 {
		t.Errorf("Height = %d, want 24", c.Height())
	}
}

func TestCanvasSetOutOfBounds(t *testing.T) {
	c := NewCanvas(10, 5)
	// Setting cells outside the canvas should not panic.
	c.Set(-1, 0, Cell{Rune: 'x'})
	c.Set(0, -1, Cell{Rune: 'x'})
	c.Set(10, 0, Cell{Rune: 'x'})
	c.Set(0, 5, Cell{Rune: 'x'})
}

func TestCanvasSetAndGet(t *testing.T) {
	c := NewCanvas(10, 5)
	cell := Cell{Rune: 'A'}
	c.Set(3, 2, cell)
	got := c.current[2][3]
	if got.Rune != 'A' {
		t.Errorf("cell at (3,2) rune = %q, want 'A'", got.Rune)
	}
}

func TestCanvasClear(t *testing.T) {
	c := NewCanvas(10, 5)
	c.Set(0, 0, Cell{Rune: 'X'})
	c.Clear()
	got := c.current[0][0]
	if !got.IsBlank() {
		t.Errorf("after Clear, cell (0,0) rune = %q, want blank", got.Rune)
	}
}

func TestCanvasResize(t *testing.T) {
	c := NewCanvas(80, 24)
	c.Resize(120, 40)
	if c.Width() != 120 {
		t.Errorf("after Resize Width = %d, want 120", c.Width())
	}
	if c.Height() != 40 {
		t.Errorf("after Resize Height = %d, want 40", c.Height())
	}
	// Buffers should match new dimensions.
	if len(c.current) != 40 {
		t.Errorf("current buffer height = %d, want 40", len(c.current))
	}
	if len(c.current[0]) != 120 {
		t.Errorf("current buffer width = %d, want 120", len(c.current[0]))
	}
}

func TestCanvasResizeToZero(t *testing.T) {
	c := NewCanvas(80, 24)
	// Resizing to zero should not panic.
	c.Resize(0, 0)
	if c.Width() != 0 || c.Height() != 0 {
		t.Errorf("after Resize(0,0) = %dx%d, want 0x0", c.Width(), c.Height())
	}
}

func TestCellEqual(t *testing.T) {
	a := Cell{Rune: 'x'}
	b := Cell{Rune: 'x'}
	c := Cell{Rune: 'y'}
	if !a.Equal(b) {
		t.Error("identical cells should be equal")
	}
	if a.Equal(c) {
		t.Error("different runes should not be equal")
	}
}

func TestCellIsBlank(t *testing.T) {
	if !Blank().IsBlank() {
		t.Error("Blank() should be blank")
	}
	if (Cell{Rune: 'x'}).IsBlank() {
		t.Error("'x' should not be blank")
	}
	if !(Cell{Rune: ' '}).IsBlank() {
		t.Error("' ' should be blank")
	}
}
