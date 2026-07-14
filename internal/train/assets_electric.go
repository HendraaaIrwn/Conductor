package train

import (
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// --- Electric Commuter Locomotive (right-facing) ---
// A sleek electric train with a pantograph on top and a streamlined nose.
// 6 rows tall for coupling alignment.

var electricLocomotiveRightF1 = []string{
	`       _`,
	`      |_|`,
	`  ____|_______`,
	` /  E  E  E  E \`,
	` \_____________/`,
	`  O   O   O   O`,
}

var electricLocomotiveRightF2 = []string{
	`       _`,
	`      |_|`,
	`  ____|_______`,
	` /  E  E  E  E \`,
	` \_____________/`,
	`  o   o   o   o`,
}

// --- Electric Commuter Locomotive (left-facing) ---

var electricLocomotiveLeftF1 = []string{
	`  _      `,
	` |_|     `,
	`  ____|_______`,
	` /  E  E  E  E \`,
	` \_____________/`,
	`  O   O   O   O`,
}

var electricLocomotiveLeftF2 = []string{
	`  _      `,
	` |_|     `,
	`  ____|_______`,
	` /  E  E  E  E \`,
	` \_____________/`,
	`  o   o   o   o`,
}

// --- Commuter Carriage (right-facing) ---
// A streamlined commuter carriage matching the electric locomotive.

var commuterCarRightF1 = []string{
	``,
	``,
	`  _______________`,
	` /  M  M  M  M  \`,
	` \_______________/`,
	`  O   O   O   O`,
}

var commuterCarRightF2 = []string{
	``,
	``,
	`  _______________`,
	` /  M  M  M  M  \`,
	` \_______________/`,
	`  o   o   o   o`,
}

// --- Commuter Carriage (left-facing) ---

var commuterCarLeftF1 = []string{
	``,
	``,
	` _______________  `,
	`/  M  M  M  M  \ `,
	`\_______________/ `,
	`  O   O   O   O  `,
}

var commuterCarLeftF2 = []string{
	``,
	``,
	` _______________  `,
	`/  M  M  M  M  \ `,
	`\_______________/ `,
	`  o   o   o   o  `,
}

// --- Caboose (right-facing) ---
// A small end-of-train cabin car with a cupola on top.

var cabooseRightF1 = []string{
	`      _`,
	`     | |`,
	`  ___| |___`,
	` |  C  C  C  |`,
	` |____________|`,
	`  O   O   O`,
}

var cabooseRightF2 = []string{
	`      _`,
	`     | |`,
	`  ___| |___`,
	` |  C  C  C  |`,
	` |____________|`,
	`  o   o   o`,
}

// --- Caboose (left-facing) ---

var cabooseLeftF1 = []string{
	`  _      `,
	` | |     `,
	` ___| |___  `,
	`|  C  C  C  | `,
	`|____________| `,
	`  O   O   O  `,
}

var cabooseLeftF2 = []string{
	`  _      `,
	` | |     `,
	` ___| |___  `,
	`|  C  C  C  | `,
	`|____________| `,
	`  o   o   o  `,
}

// ElectricLocomotiveRight returns the right-facing electric locomotive sprite.
func ElectricLocomotiveRight(style tcell.Style) *render.Sprite {
	return buildSprite("electric-locomotive-right", style,
		electricLocomotiveRightF1, electricLocomotiveRightF2)
}

// ElectricLocomotiveLeft returns the left-facing electric locomotive sprite.
func ElectricLocomotiveLeft(style tcell.Style) *render.Sprite {
	return buildSprite("electric-locomotive-left", style,
		electricLocomotiveLeftF1, electricLocomotiveLeftF2)
}

// CommuterCarRight returns the right-facing commuter carriage sprite.
func CommuterCarRight(style tcell.Style) *render.Sprite {
	return buildSprite("commuter-car-right", style,
		commuterCarRightF1, commuterCarRightF2)
}

// CommuterCarLeft returns the left-facing commuter carriage sprite.
func CommuterCarLeft(style tcell.Style) *render.Sprite {
	return buildSprite("commuter-car-left", style,
		commuterCarLeftF1, commuterCarLeftF2)
}

// CabooseRight returns the right-facing caboose sprite.
func CabooseRight(style tcell.Style) *render.Sprite {
	return buildSprite("caboose-right", style,
		cabooseRightF1, cabooseRightF2)
}

// CabooseLeft returns the left-facing caboose sprite.
func CabooseLeft(style tcell.Style) *render.Sprite {
	return buildSprite("caboose-left", style,
		cabooseLeftF1, cabooseLeftF2)
}
