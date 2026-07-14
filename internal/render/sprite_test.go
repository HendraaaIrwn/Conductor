package render

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestSpriteValidateValid(t *testing.T) {
	s := NewSprite("test",
		ParseFrame("AB\nCD", tcell.StyleDefault),
		ParseFrame("EF\nGH", tcell.StyleDefault),
	)
	if err := s.Validate(); err != nil {
		t.Errorf("valid sprite failed: %v", err)
	}
}

func TestSpriteValidateNoFrames(t *testing.T) {
	s := NewSprite("empty")
	if err := s.Validate(); err == nil {
		t.Error("sprite with no frames should fail validation")
	}
}

func TestSpriteValidateZeroWidth(t *testing.T) {
	s := NewSprite("zero-width",
		ParseFrame("", tcell.StyleDefault),
	)
	if err := s.Validate(); err == nil {
		t.Error("sprite with zero width should fail validation")
	}
}

func TestSpriteValidateAllBlankFrame(t *testing.T) {
	s := NewSprite("blank",
		ParseFrame("   \n   ", tcell.StyleDefault),
	)
	if err := s.Validate(); err == nil {
		t.Error("all-blank frame should fail validation")
	}
}

func TestSpriteValidateFrameHeightMismatch(t *testing.T) {
	s := &Sprite{
		Name:   "mismatch",
		Width:  2,
		Height: 2,
		Frames: []CellGrid{
			ParseFrame("AB\nCD", tcell.StyleDefault),
			ParseFrame("EF", tcell.StyleDefault), // only 1 row
		},
	}
	if err := s.Validate(); err == nil {
		t.Error("frame height mismatch should fail validation")
	}
}

func TestSpriteValidateFrameWidthMismatch(t *testing.T) {
	// ParseFrame pads each frame to its own max width. "ABC" → width 3,
	// "AB" → width 2. When combined in a sprite with Width=3, the second
	// frame's rows are only 2 cells wide, which should fail validation.
	s := &Sprite{
		Name:   "mismatch-w",
		Width:  3,
		Height: 1,
		Frames: []CellGrid{
			ParseFrame("ABC", tcell.StyleDefault),
			ParseFrame("AB", tcell.StyleDefault),
		},
	}
	if err := s.Validate(); err == nil {
		t.Error("frame width mismatch should fail validation")
	}
}

func TestValidateCouplingCompatible(t *testing.T) {
	a := NewSprite("loco",
		ParseFrame("AB\nCD", tcell.StyleDefault))
	b := NewSprite("car",
		ParseFrame("EF\nGH", tcell.StyleDefault))
	if err := ValidateCoupling(a, b); err != nil {
		t.Errorf("same-height sprites should be compatible: %v", err)
	}
}

func TestValidateCouplingMismatch(t *testing.T) {
	a := NewSprite("loco",
		ParseFrame("AB\nCD\nEF", tcell.StyleDefault))
	b := NewSprite("car",
		ParseFrame("AB\nCD", tcell.StyleDefault))
	if err := ValidateCoupling(a, b); err == nil {
		t.Error("different-height sprites should fail coupling validation")
	}
}

func TestSpriteDraw(t *testing.T) {
	s := NewSprite("test",
		ParseFrame("AB\nCD", tcell.StyleDefault))
	canvas := NewCanvas(10, 10)
	s.Draw(canvas, 2, 3, 0)
	if canvas.CellAt(2, 3).Rune != 'A' {
		t.Errorf("cell (2,3) = %q, want 'A'", canvas.CellAt(2, 3).Rune)
	}
	if canvas.CellAt(3, 3).Rune != 'B' {
		t.Errorf("cell (3,3) = %q, want 'B'", canvas.CellAt(3, 3).Rune)
	}
	if canvas.CellAt(2, 4).Rune != 'C' {
		t.Errorf("cell (2,4) = %q, want 'C'", canvas.CellAt(2, 4).Rune)
	}
}

func TestSpriteDrawOutOfBounds(t *testing.T) {
	s := NewSprite("test",
		ParseFrame("AB\nCD", tcell.StyleDefault))
	canvas := NewCanvas(10, 10)
	// Drawing at negative position should not panic.
	s.Draw(canvas, -1, -1, 0)
}

func TestSpriteDrawSkipsBlank(t *testing.T) {
	// Sprite "A \n B" has a blank at (col=1,row=0) and (col=0,row=1).
	s := NewSprite("test",
		ParseFrame("A \n B", tcell.StyleDefault))
	canvas := NewCanvas(10, 10)
	// Place an 'X' at (3,3) — this is where the blank cell in the sprite
	// would land when drawn at (2,3). The blank should NOT overwrite 'X'.
	canvas.SetRune(3, 3, 'X', tcell.StyleDefault)
	s.Draw(canvas, 2, 3, 0)
	if canvas.CellAt(3, 3).Rune != 'X' {
		t.Errorf("blank cell should not overwrite, got %q", canvas.CellAt(3, 3).Rune)
	}
	// 'A' at (2,3) should be drawn.
	if canvas.CellAt(2, 3).Rune != 'A' {
		t.Errorf("cell (2,3) = %q, want 'A'", canvas.CellAt(2, 3).Rune)
	}
}
