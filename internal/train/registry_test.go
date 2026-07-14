package train

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestParseTrainType(t *testing.T) {
	cases := []struct {
		input string
		want  TrainType
		ok    bool
	}{
		{"steam", TypeSteam, true},
		{"diesel", TypeDiesel, true},
		{"electric", TypeElectric, true},
		{"unknown", 0, false},
		{"", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseTrainType(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseTrainType(%q) ok = %v, want %v", tc.input, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Errorf("ParseTrainType(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestTrainTypeString(t *testing.T) {
	if TypeSteam.String() != "steam" {
		t.Errorf("steam string = %q", TypeSteam.String())
	}
	if TypeDiesel.String() != "diesel" {
		t.Errorf("diesel string = %q", TypeDiesel.String())
	}
	if TypeElectric.String() != "electric" {
		t.Errorf("electric string = %q", TypeElectric.String())
	}
}

func TestAllTrainTypes(t *testing.T) {
	types := AllTrainTypes()
	if len(types) != 3 {
		t.Errorf("AllTrainTypes length = %d, want 3", len(types))
	}
}

func TestLocomotiveSpec(t *testing.T) {
	for _, tt := range AllTrainTypes() {
		spec := LocomotiveSpec(tt)
		if spec == nil {
			t.Errorf("no spec for %v", tt)
			continue
		}
		if spec.Type != tt {
			t.Errorf("spec type = %v, want %v", spec.Type, tt)
		}
		if len(spec.CompatibleCars) == 0 {
			t.Errorf("%v has no compatible cars", tt)
		}
	}
}

func TestCarriageSpec(t *testing.T) {
	categories := []CarriageCategory{
		CatPassenger, CatBoxcar, CatTank, CatOpenCargo, CatCaboose, CatCommuter,
	}
	for _, cat := range categories {
		spec := CarriageSpec(cat)
		if spec == nil {
			t.Errorf("no spec for %v", cat)
			continue
		}
		if spec.Category != cat {
			t.Errorf("spec category = %v, want %v", spec.Category, cat)
		}
		if spec.Right == nil {
			t.Errorf("%v has no right sprite factory", cat)
		}
		if spec.Left == nil {
			t.Errorf("%v has no left sprite factory", cat)
		}
	}
}

func TestLocomotiveSpriteBothDirections(t *testing.T) {
	style := tcell.StyleDefault
	for _, tt := range AllTrainTypes() {
		right := LocomotiveSprite(tt, LeftToRight, style)
		if right == nil {
			t.Errorf("no right sprite for %v", tt)
		}
		left := LocomotiveSprite(tt, RightToLeft, style)
		if left == nil {
			t.Errorf("no left sprite for %v", tt)
		}
		if right != nil && left != nil && right == left {
			t.Errorf("%v: left and right sprites should differ", tt)
		}
	}
}

func TestCarriageSpriteBothDirections(t *testing.T) {
	style := tcell.StyleDefault
	categories := []CarriageCategory{
		CatPassenger, CatBoxcar, CatTank, CatOpenCargo, CatCaboose, CatCommuter,
	}
	for _, cat := range categories {
		right := CarriageSprite(cat, LeftToRight, style)
		if right == nil {
			t.Errorf("no right sprite for %v", cat)
		}
		left := CarriageSprite(cat, RightToLeft, style)
		if left == nil {
			t.Errorf("no left sprite for %v", cat)
		}
	}
}

func TestCompatibleCarriages(t *testing.T) {
	steam := CompatibleCarriages(TypeSteam)
	if len(steam) == 0 {
		t.Error("steam should have compatible carriages")
	}
	diesel := CompatibleCarriages(TypeDiesel)
	if len(diesel) == 0 {
		t.Error("diesel should have compatible carriages")
	}
	electric := CompatibleCarriages(TypeElectric)
	if len(electric) == 0 {
		t.Error("electric should have compatible carriages")
	}
}

func TestEmitsSmoke(t *testing.T) {
	if !EmitsSmoke(TypeSteam) {
		t.Error("steam should emit smoke")
	}
	if EmitsSmoke(TypeDiesel) {
		t.Error("diesel should not emit smoke")
	}
	if EmitsSmoke(TypeElectric) {
		t.Error("electric should not emit smoke")
	}
}

func TestTrainSpeed(t *testing.T) {
	for _, tt := range AllTrainTypes() {
		speed := TrainSpeed(tt)
		if speed <= 0 {
			t.Errorf("%v speed = %f, want > 0", tt, speed)
		}
	}
}

func TestAllSpritesValidate(t *testing.T) {
	style := tcell.StyleDefault
	// Locomotives
	for _, tt := range AllTrainTypes() {
		right := LocomotiveSprite(tt, LeftToRight, style)
		if err := right.Validate(); err != nil {
			t.Errorf("%v right locomotive: %v", tt, err)
		}
		left := LocomotiveSprite(tt, RightToLeft, style)
		if err := left.Validate(); err != nil {
			t.Errorf("%v left locomotive: %v", tt, err)
		}
	}
	// Carriages
	categories := []CarriageCategory{
		CatPassenger, CatBoxcar, CatTank, CatOpenCargo, CatCaboose, CatCommuter,
	}
	for _, cat := range categories {
		right := CarriageSprite(cat, LeftToRight, style)
		if err := right.Validate(); err != nil {
			t.Errorf("%v right carriage: %v", cat, err)
		}
		left := CarriageSprite(cat, RightToLeft, style)
		if err := left.Validate(); err != nil {
			t.Errorf("%v left carriage: %v", cat, err)
		}
	}
}

func TestAllSpritesHaveTwoFrames(t *testing.T) {
	style := tcell.StyleDefault
	for _, tt := range AllTrainTypes() {
		right := LocomotiveSprite(tt, LeftToRight, style)
		if right.FrameCount() < 2 {
			t.Errorf("%v right locomotive frames = %d, want >= 2", tt, right.FrameCount())
		}
	}
	categories := []CarriageCategory{
		CatPassenger, CatBoxcar, CatTank, CatOpenCargo, CatCaboose, CatCommuter,
	}
	for _, cat := range categories {
		right := CarriageSprite(cat, LeftToRight, style)
		if right.FrameCount() < 2 {
			t.Errorf("%v right carriage frames = %d, want >= 2", cat, right.FrameCount())
		}
	}
}

func TestAllSpritesSameHeight(t *testing.T) {
	style := tcell.StyleDefault
	// All sprites (locomotives + carriages) must be 6 rows tall for coupling.
	const wantHeight = 6
	for _, tt := range AllTrainTypes() {
		right := LocomotiveSprite(tt, LeftToRight, style)
		if right.Height != wantHeight {
			t.Errorf("%v right height = %d, want %d", tt, right.Height, wantHeight)
		}
		left := LocomotiveSprite(tt, RightToLeft, style)
		if left.Height != wantHeight {
			t.Errorf("%v left height = %d, want %d", tt, left.Height, wantHeight)
		}
	}
	categories := []CarriageCategory{
		CatPassenger, CatBoxcar, CatTank, CatOpenCargo, CatCaboose, CatCommuter,
	}
	for _, cat := range categories {
		right := CarriageSprite(cat, LeftToRight, style)
		if right.Height != wantHeight {
			t.Errorf("%v right height = %d, want %d", cat, right.Height, wantHeight)
		}
		left := CarriageSprite(cat, RightToLeft, style)
		if left.Height != wantHeight {
			t.Errorf("%v left height = %d, want %d", cat, left.Height, wantHeight)
		}
	}
}
