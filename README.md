# Conductor

> A tiny railway journey inside your terminal.

Conductor is an ambient terminal animation that displays an animated railway
scene: locomotives crossing the screen, wheels turning, and a track rolling
by. It is inspired by terminal screensavers like Asciiquarium, but all
architecture, source code, and ASCII art are original.

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
