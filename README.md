# Conductor

> A tiny railway journey inside your terminal.

Conductor is an ambient terminal animation that displays an animated railway
scene: locomotives crossing the screen, wheels turning, and a track rolling
by. It is inspired by terminal screensavers like Asciiquarium, but all
architecture, source code, and ASCII art are original.

## Current Status

**Milestone 6 — CLI and Configuration** is complete. Conductor currently provides:

- Single binary, no runtime dependencies.
- tcell-based terminal rendering with alternate screen and cursor hiding.
- Double-buffered canvas with cell-level diffing (no full-screen flicker).
- Delta-time animation clock (frame-rate independent movement).
- Entity manager with add/remove/update/query, layer-sorted rendering.
- Layer system (sky=0, celestial=10, distant=20, nearby=30, platform=40,
  track=50, train=60, foreground=70, particles=80, overlay=90).
- Behavior interface with composable built-in behaviors.
- Entity lifecycle: creation, active, lifetime expiration, off-screen removal.
- Sprite system with multi-frame animation, validation, and coupling checks.
- Three train types with six carriage types, directional variants, and
  procedural consist generation.
- Animated wheels and smoke particles from steam locomotives.
- Time system: morning, day, sunset, night — affects palette, sun/moon/stars.
- Weather system: clear, rain, snow — particle entities with recycling.
- Random event scheduler: long freight, smoke burst, birds, lightning, signal
  change — with cooldowns and deterministic seeding.
- Scene interface with countryside, station, and mountain scenes.
- Signal entities with red/green state machine (shape-based accessibility).
- Cloud entities that drift and wrap around the viewport.
- **Full CLI with Cobra**: --train, --scene, --weather, --time, --speed,
  --fps, --seed, --color, --no-color, --reduced-motion, --config, --version,
  --help. All flags validated with clear error messages.
- **TOML configuration file**: ~/.config/conductor/config.toml with
  precedence: defaults → config file → CLI flags. Missing file and unknown
  fields handled gracefully.
- **No-color mode**: removes all color from rendering, preserving readability
  via shape-based distinctions (signals, celestial elements).
- **Reduced-motion mode**: lowers FPS, reduces particle counts, disables
  lightning, configurable via CLI, config file, or 'c' key toggle.
- **Help overlay**: press 'h' or '?' to show/hide keyboard controls. Esc
  closes it. Centered overlay with border, readable in any terminal.
- **Color palette cycling**: 'c' key toggles color/no-color mode.
- Scene switching via 's' key, weather via 'w' key, time via 'd' key.
- Keyboard controls: quit, pause, regenerate, next train, change scene,
  change weather, change time, toggle color, toggle help, speed control.
- Terminal resize handling, graceful terminal restoration, panic recovery.
- Minimum terminal size check (80x24) with a readable message.
- 261 automated tests covering all systems, including input (100% coverage),
  snapshot, integration, and soak tests.

## Installation

### From Source

```bash
go install github.com/example/conductor/cmd/conductor@latest
```

### Build Manually

```bash
git clone https://github.com/example/conductor.git
cd conductor
go build -o conductor ./cmd/conductor
```

## Usage

```bash
conductor
```

### CLI Options

```
--train <type>       Train type: steam, diesel, electric, or random
--scene <name>       Scene: countryside, station, mountain, or random
--weather <type>     Weather: clear, rain, snow, or random
--time <period>      Time of day: morning, day, sunset, night, or auto
--speed <float>      Animation speed multiplier (0.25-3.0, default: 1.0)
--fps <int>          Target frame rate (5-60, default: 20)
--seed <int>         Random seed for deterministic scenes
--color <mode>       Color mode: auto, 16, 256, or truecolor
--no-color           Disable all colors
--reduced-motion     Reduce particles and frame rate
--config <path>      Path to configuration file
--version            Print version and exit
--help               Print help and exit
```

### Configuration File

Default location: `~/.config/conductor/config.toml`

```toml
train = "random"
scene = "random"
weather = "random"
time = "auto"

fps = 20
speed = 1.0
color = "auto"

show_status = false
reduced_motion = false
random_events = true
```

### Configuration Precedence

```
Built-in defaults → Config file → CLI flags
```

Missing config file is not an error. Unknown TOML fields are ignored.

Full CLI flags (--train, --scene, --weather, --time, --speed, --color,
--no-color, --reduced-motion, --config) will be added in Milestone 6.

## Keyboard Controls

```
q or Ctrl+C    Quit
p or Space     Pause or resume
r              Regenerate the scene
n              Spawn or schedule the next train
s              Change scene
w              Change weather (clear → rain → snow)
d              Change time of day (morning → day → sunset → night)
c              Toggle color / no-color mode
h or ?         Toggle help overlay
Esc            Close help overlay
+              Increase animation speed
-              Decrease animation speed
0              Reset animation speed
```

While paused, all movement and animation stops. A small `[PAUSED]` indicator
appears in the top-right corner. Resize events still work while paused.

## Development

### Prerequisites

- Go 1.23 or later

### Build

```bash
go build -o conductor ./cmd/conductor
```

### Test

```bash
go test ./...
```

### Format and Vet

```bash
gofmt -w .
go vet ./...
```

## Project Structure

```
conductor/
├── cmd/conductor/        Entry point
├── internal/
│   ├── app/              Application lifecycle and signal handling
│   ├── engine/           Animation loop, clock, world state, event scheduler
│   ├── entity/           Entity, manager, and behavior system
│   ├── effects/          Smoke, rain, snow, lightning, bird, weather system
│   ├── render/           Canvas, cell, sprite, layers (double buffering)
│   ├── scene/            Scene interface, countryside, station, mountain, signals
│   ├── train/            Train types, consist generator, ASCII art assets
│   ├── input/            Keyboard input handler
│   └── terminal/         tcell wrapper and terminal cleanup
├── assets/               (Reserved for embedded sprite files in later milestones)
├── go.mod
├── LICENSE
└── README.md
```

## Platform Support

- macOS (tested)
- Linux
- Windows

## License

MIT. See [LICENSE](LICENSE).

All ASCII art in this project is original work created specifically for
Conductor. It is not derived from Asciiquarium or any other third-party
source.
