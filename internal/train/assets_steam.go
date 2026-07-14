package train

import (
	"github.com/example/conductor/internal/render"
	"github.com/gdamore/tcell/v2"
)

// All ASCII art in this file is original work created specifically for
// Conductor. It is not derived from Asciiquarium or any other third-party
// source.
//
// All sprites are 6 rows tall to ensure coupling alignment across all train
// types. Right-facing sprites face the direction of travel for
// LeftToRight trains; left-facing sprites face the direction of travel for
// RightToLeft trains.

// --- Steam Locomotive (right-facing) ---

var steamLocomotiveRightF1 = []string{
	`  ___           _`,
	` |   |         / \`,
	` |___|________|   |`,
	` | o  o  o  o     |`,
	` |________________|`,
	`  O   O   O   O`,
}

var steamLocomotiveRightF2 = []string{
	`  ___           _`,
	` |   |         / \`,
	` |___|________|   |`,
	` | o  o  o  o     |`,
	` |________________|`,
	`  o   o   o   o`,
}

// --- Steam Locomotive (left-facing) ---
// The chimney is at the front (left side) and the cab is at the rear (right).

var steamLocomotiveLeftF1 = []string{
	`  _           ___`,
	` / \         |   |`,
	`|   |________|___|`,
	`|     o  o  o  o  |`,
	`|________________|`,
	`  O   O   O   O`,
}

var steamLocomotiveLeftF2 = []string{
	`  _           ___`,
	` / \         |   |`,
	`|   |________|___|`,
	`|     o  o  o  o  |`,
	`|________________|`,
	`  o   o   o   o`,
}

// --- Passenger Carriage (right-facing) ---

var passengerCarRightF1 = []string{
	``,
	``,
	`  _______________`,
	` |  o  o  o  o  |`,
	` |_______________|`,
	`  O   O   O   O`,
}

var passengerCarRightF2 = []string{
	``,
	``,
	`  _______________`,
	` |  o  o  o  o  |`,
	` |_______________|`,
	`  o   o   o   o`,
}

// --- Passenger Carriage (left-facing) ---

var passengerCarLeftF1 = []string{
	``,
	``,
	` _______________  `,
	`|  o  o  o  o  | `,
	`|_______________| `,
	`  O   O   O   O  `,
}

var passengerCarLeftF2 = []string{
	``,
	``,
	` _______________  `,
	`|  o  o  o  o  | `,
	`|_______________| `,
	`  o   o   o   o  `,
}

// SteamLocomotiveRight returns the right-facing steam locomotive sprite.
func SteamLocomotiveRight(style tcell.Style) *render.Sprite {
	return buildSprite("steam-locomotive-right", style,
		steamLocomotiveRightF1, steamLocomotiveRightF2)
}

// SteamLocomotiveLeft returns the left-facing steam locomotive sprite.
func SteamLocomotiveLeft(style tcell.Style) *render.Sprite {
	return buildSprite("steam-locomotive-left", style,
		steamLocomotiveLeftF1, steamLocomotiveLeftF2)
}

// PassengerCarRight returns the right-facing passenger carriage sprite.
func PassengerCarRight(style tcell.Style) *render.Sprite {
	return buildSprite("passenger-car-right", style,
		passengerCarRightF1, passengerCarRightF2)
}

// PassengerCarLeft returns the left-facing passenger carriage sprite.
func PassengerCarLeft(style tcell.Style) *render.Sprite {
	return buildSprite("passenger-car-left", style,
		passengerCarLeftF1, passengerCarLeftF2)
}
