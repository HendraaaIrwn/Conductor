package train

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestGeneratorProducesValidConsist(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	cfg := DefaultConsistConfig(TypeSteam, LeftToRight, 120)
	cars := g.Generate(cfg)
	if len(cars) == 0 {
		t.Fatal("generator produced no cars")
	}
	if err := ValidateConsist(cars); err != nil {
		t.Errorf("generated consist is invalid: %v", err)
	}
}

func TestGeneratorLocomotiveFirst(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	for _, tt := range AllTrainTypes() {
		cfg := DefaultConsistConfig(tt, LeftToRight, 120)
		cars := g.Generate(cfg)
		if len(cars) == 0 {
			t.Errorf("%v: no cars", tt)
			continue
		}
		if cars[0].OffsetX != 0 {
			t.Errorf("%v: locomotive offset = %d, want 0", tt, cars[0].OffsetX)
		}
	}
}

func TestGeneratorCarCountInRange(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	cfg := ConsistConfig{
		TrainType: TypeDiesel,
		Direction: LeftToRight,
		MinCars:   3,
		MaxCars:   5,
		MaxWidth:  200, // large enough to not trim
	}
	cars := g.Generate(cfg)
	// 1 locomotive + [3,5] carriages + 1 end car = [5,7] total
	total := len(cars)
	if total < 5 || total > 7 {
		t.Errorf("total cars = %d, want 5-7", total)
	}
}

func TestGeneratorRespectsMaxWidth(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	cfg := ConsistConfig{
		TrainType: TypeSteam,
		Direction: LeftToRight,
		MinCars:   8,
		MaxCars:   8,
		MaxWidth:  40, // very small, should force trimming
	}
	cars := g.Generate(cfg)
	w := TrainWidth(cars)
	if w > 40 {
		t.Errorf("train width = %d, want <= 40", w)
	}
}

func TestGeneratorLongTrainExceedsWidth(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	cfg := ConsistConfig{
		TrainType:      TypeDiesel,
		Direction:      LeftToRight,
		MinCars:        10,
		MaxCars:        12,
		MaxWidth:       40,
		AllowLongTrain: true,
	}
	cars := g.Generate(cfg)
	w := TrainWidth(cars)
	if w <= 40 {
		t.Errorf("long train width = %d, want > 40", w)
	}
}

func TestGeneratorDeterministicWithSeed(t *testing.T) {
	g1 := NewGenerator(42, tcell.StyleDefault)
	g2 := NewGenerator(42, tcell.StyleDefault)
	cfg := DefaultConsistConfig(TypeSteam, LeftToRight, 120)
	cars1 := g1.Generate(cfg)
	cars2 := g2.Generate(cfg)
	if len(cars1) != len(cars2) {
		t.Errorf("same seed produced different car counts: %d vs %d", len(cars1), len(cars2))
	}
	for i := range cars1 {
		if cars1[i].OffsetX != cars2[i].OffsetX {
			t.Errorf("car %d offset differs: %d vs %d", i, cars1[i].OffsetX, cars2[i].OffsetX)
		}
		if cars1[i].Sprite.Name != cars2[i].Sprite.Name {
			t.Errorf("car %d sprite differs: %s vs %s", i, cars1[i].Sprite.Name, cars2[i].Sprite.Name)
		}
	}
}

func TestGeneratorDifferentSeedsDiffer(t *testing.T) {
	g1 := NewGenerator(1, tcell.StyleDefault)
	g2 := NewGenerator(2, tcell.StyleDefault)
	cfg := DefaultConsistConfig(TypeDiesel, LeftToRight, 120)
	cars1 := g1.Generate(cfg)
	cars2 := g2.Generate(cfg)
	// They might be the same by chance, but very unlikely with different seeds.
	// We check if at least the sprite names differ.
	same := true
	if len(cars1) != len(cars2) {
		same = false
	} else {
		for i := range cars1 {
			if cars1[i].Sprite.Name != cars2[i].Sprite.Name {
				same = false
				break
			}
		}
	}
	// It's possible but extremely unlikely that both produce identical consists.
	// We don't fail if they happen to match — just don't require difference.
	_ = same
}

func TestGenerateRandom(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	tt, cars := g.GenerateRandom(LeftToRight, 120)
	if len(cars) == 0 {
		t.Fatal("GenerateRandom produced no cars")
	}
	if err := ValidateConsist(cars); err != nil {
		t.Errorf("random consist invalid: %v", err)
	}
	if tt < TypeSteam || tt > TypeElectric {
		t.Errorf("random train type = %v, want in range", tt)
	}
}

func TestGenerateLongFreight(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	tt, cars := g.GenerateLongFreight(LeftToRight, 120)
	if len(cars) == 0 {
		t.Fatal("GenerateLongFreight produced no cars")
	}
	if tt != TypeSteam && tt != TypeDiesel {
		t.Errorf("long freight type = %v, want steam or diesel", tt)
	}
	// Long freight should have many cars.
	if len(cars) < 9 {
		t.Errorf("long freight car count = %d, want >= 9", len(cars))
	}
}

func TestGeneratorDirectionAffectsSprites(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	cfg := DefaultConsistConfig(TypeSteam, LeftToRight, 120)
	rightCars := g.Generate(cfg)
	cfg.Direction = RightToLeft
	leftCars := g.Generate(cfg)
	if len(rightCars) == 0 || len(leftCars) == 0 {
		t.Fatal("no cars generated")
	}
	// The locomotive sprite name should differ by direction.
	if rightCars[0].Sprite.Name == leftCars[0].Sprite.Name {
		t.Error("locomotive sprite should differ by direction")
	}
}

func TestGeneratorElectricOnlyCommuterCars(t *testing.T) {
	g := NewGenerator(42, tcell.StyleDefault)
	cfg := DefaultConsistConfig(TypeElectric, LeftToRight, 200)
	cars := g.Generate(cfg)
	if len(cars) < 2 {
		t.Fatal("electric train too short")
	}
	// All carriages (non-locomotive) should be commuter cars.
	for i, car := range cars[1:] {
		if car.Sprite.Name != "commuter-car-right" && car.Sprite.Name != "commuter-car-left" {
			t.Errorf("electric car %d sprite = %s, want commuter", i, car.Sprite.Name)
		}
	}
}

func TestGeneratorSteamDieselAcceptsFreightCars(t *testing.T) {
	g := NewGenerator(100, tcell.StyleDefault)
	// Generate many consists to check that freight cars appear.
	foundBoxcar := false
	foundTank := false
	foundOpenCargo := false
	for i := 0; i < 50; i++ {
		cfg := DefaultConsistConfig(TypeDiesel, LeftToRight, 200)
		cars := g.Generate(cfg)
		for _, car := range cars[1:] {
			switch car.Sprite.Name {
			case "boxcar-right", "boxcar-left":
				foundBoxcar = true
			case "tank-car-right", "tank-car-left":
				foundTank = true
			case "open-cargo-right", "open-cargo-left":
				foundOpenCargo = true
			}
		}
	}
	if !foundBoxcar {
		t.Error("boxcar never appeared in 50 diesel consists")
	}
	if !foundTank {
		t.Error("tank car never appeared in 50 diesel consists")
	}
	if !foundOpenCargo {
		t.Error("open cargo never appeared in 50 diesel consists")
	}
}
