package train

import (
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// --- Diesel Locomotive (right-facing) ---
// A long, boxy diesel locomotive with a headlight at the front (right side).
// 6 rows tall for coupling alignment.

var dieselLocomotiveRightF1 = []string{
	``,
	``,
	`  _________________`,
	` |  D  D  D  D  D  |`,
	` |_________________|`,
	`  O   O   O   O   O`,
}

var dieselLocomotiveRightF2 = []string{
	``,
	``,
	`  _________________`,
	` |  D  D  D  D  D  |`,
	` |_________________|`,
	`  o   o   o   o   o`,
}

// --- Diesel Locomotive (left-facing) ---

var dieselLocomotiveLeftF1 = []string{
	``,
	``,
	`_________________  `,
	`|  D  D  D  D  D  | `,
	`|_________________| `,
	`  O   O   O   O   O`,
}

var dieselLocomotiveLeftF2 = []string{
	``,
	``,
	`_________________  `,
	`|  D  D  D  D  D  | `,
	`|_________________| `,
	`  o   o   o   o   o`,
}

// --- Boxcar (right-facing) ---
// An enclosed freight car with a sliding door marked 'B'.

var boxcarRightF1 = []string{
	``,
	``,
	`  _______________`,
	` |  B  B  B  B  |`,
	` |_______________|`,
	`  O   O   O   O`,
}

var boxcarRightF2 = []string{
	``,
	``,
	`  _______________`,
	` |  B  B  B  B  |`,
	` |_______________|`,
	`  o   o   o   o`,
}

// --- Boxcar (left-facing) ---

var boxcarLeftF1 = []string{
	``,
	``,
	` _______________  `,
	`|  B  B  B  B  | `,
	`|_______________| `,
	`  O   O   O   O  `,
}

var boxcarLeftF2 = []string{
	``,
	``,
	` _______________  `,
	`|  B  B  B  B  | `,
	`|_______________| `,
	`  o   o   o   o  `,
}

// --- Tank Car (right-facing) ---
// A cylindrical tank on a flatbed, marked 'T'.

var tankCarRightF1 = []string{
	``,
	`  ~~~~~~~~~~~~~~~`,
	` /  T  T  T  T  \`,
	` \_______________/`,
	`  _______________`,
	`  O   O   O   O`,
}

var tankCarRightF2 = []string{
	``,
	`  ~~~~~~~~~~~~~~~`,
	` /  T  T  T  T  \`,
	` \_______________/`,
	`  _______________`,
	`  o   o   o   o`,
}

// --- Tank Car (left-facing) ---

var tankCarLeftF1 = []string{
	``,
	`~~~~~~~~~~~~~~~  `,
	`/  T  T  T  T  \ `,
	`\_______________/ `,
	`  _______________ `,
	`  O   O   O   O  `,
}

var tankCarLeftF2 = []string{
	``,
	`~~~~~~~~~~~~~~~  `,
	`/  T  T  T  T  \ `,
	`\_______________/ `,
	`  _______________ `,
	`  o   o   o   o  `,
}

// --- Open Cargo Car (right-facing) ---
// A low-sided open car with visible cargo, marked 'C'.

var openCargoRightF1 = []string{
	``,
	``,
	`  _______________`,
	` | C C C C C C C |`,
	` |_______________|`,
	`  O   O   O   O`,
}

var openCargoRightF2 = []string{
	``,
	``,
	`  _______________`,
	` | C C C C C C C |`,
	` |_______________|`,
	`  o   o   o   o`,
}

// --- Open Cargo Car (left-facing) ---

var openCargoLeftF1 = []string{
	``,
	``,
	` _______________  `,
	`| C C C C C C C | `,
	`|_______________| `,
	`  O   O   O   O  `,
}

var openCargoLeftF2 = []string{
	``,
	``,
	` _______________  `,
	`| C C C C C C C | `,
	`|_______________| `,
	`  o   o   o   o  `,
}

// DieselLocomotiveRight returns the right-facing diesel locomotive sprite.
func DieselLocomotiveRight(style tcell.Style) *render.Sprite {
	return buildSprite("diesel-locomotive-right", style,
		dieselLocomotiveRightF1, dieselLocomotiveRightF2)
}

// DieselLocomotiveLeft returns the left-facing diesel locomotive sprite.
func DieselLocomotiveLeft(style tcell.Style) *render.Sprite {
	return buildSprite("diesel-locomotive-left", style,
		dieselLocomotiveLeftF1, dieselLocomotiveLeftF2)
}

// BoxcarRight returns the right-facing boxcar sprite.
func BoxcarRight(style tcell.Style) *render.Sprite {
	return buildSprite("boxcar-right", style,
		boxcarRightF1, boxcarRightF2)
}

// BoxcarLeft returns the left-facing boxcar sprite.
func BoxcarLeft(style tcell.Style) *render.Sprite {
	return buildSprite("boxcar-left", style,
		boxcarLeftF1, boxcarLeftF2)
}

// TankCarRight returns the right-facing tank car sprite.
func TankCarRight(style tcell.Style) *render.Sprite {
	return buildSprite("tank-car-right", style,
		tankCarRightF1, tankCarRightF2)
}

// TankCarLeft returns the left-facing tank car sprite.
func TankCarLeft(style tcell.Style) *render.Sprite {
	return buildSprite("tank-car-left", style,
		tankCarLeftF1, tankCarLeftF2)
}

// OpenCargoRight returns the right-facing open cargo car sprite.
func OpenCargoRight(style tcell.Style) *render.Sprite {
	return buildSprite("open-cargo-right", style,
		openCargoRightF1, openCargoRightF2)
}

// OpenCargoLeft returns the left-facing open cargo car sprite.
func OpenCargoLeft(style tcell.Style) *render.Sprite {
	return buildSprite("open-cargo-left", style,
		openCargoLeftF1, openCargoLeftF2)
}
