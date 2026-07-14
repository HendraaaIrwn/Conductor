# Contributing to Conductor

Thank you for considering contributing to Conductor! This document outlines
the contribution process, code rules, and guidelines for adding new assets.

## Code of Conduct

Be respectful, constructive, and inclusive. We're all here to build something
fun.

## Getting Started

1. Fork the repository.
2. Clone your fork:
   ```bash
   git clone https://github.com/<your-username>/conductor.git
   ```
3. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. Install Go 1.23 or later.
5. Run the tests to verify your setup:
   ```bash
   go test ./...
   ```

## Code Rules

- All code must be formatted with `gofmt`.
- All code must pass `go vet ./...` without warnings.
- All tests must pass: `go test ./...` must succeed.
- Prefer simple, explicit Go code. Avoid premature abstraction.
- Avoid global mutable state.
- Avoid one goroutine per entity.
- Keep simulation and rendering separate.
- Errors must be wrapped with context using `fmt.Errorf("...: %w", err)`.
- Do not expose stack traces for ordinary configuration errors.
- Do not introduce databases, web servers, networking, or plugins.
- Add comments only where they explain non-obvious behavior.

## Adding a Locomotive

1. **Create the ASCII art** in `internal/train/` as a new asset file
   (e.g. `assets_mytype.go`). Each locomotive needs:
   - A right-facing variant with two wheel animation frames.
   - A left-facing variant with two wheel animation frames.
   - All sprites must be **6 rows tall** for coupling alignment.
   - All sprites must have a `Name` that matches the pattern
     `mytype-locomotive-right` and `mytype-locomotive-left`.

2. **Register the locomotive** in `internal/train/registry.go`:
   - Add a new `TrainType` constant in `types.go`.
   - Add a `registerLocomotive` call in `init()` with the sprite factories,
     compatible carriage categories, and optional end car.
   - Add the type to `AllTrainTypes()`.

3. **Add sprite factories** in the same asset file:
   ```go
   func MyTypeLocomotiveRight(style tcell.Style) *render.Sprite { ... }
   func MyTypeLocomotiveLeft(style tcell.Style) *render.Sprite { ... }
   ```

4. **Test the locomotive**:
   ```go
   go test -v -run TestAllSpritesValidate ./internal/train/
   ```

## Adding a Carriage

1. **Create the ASCII art** — same rules as locomotives (6 rows tall,
   two wheel animation frames, left and right variants).

2. **Register the carriage** in `internal/train/registry.go`:
   - Add a new `CarriageCategory` constant in `types.go`.
   - Add a `registerCarriage` call in `init()`.
   - Add the category to the compatible list of any locomotive that should
     use it.

3. **Add sprite factories**:
   ```go
   func MyCarriageRight(style tcell.Style) *render.Sprite { ... }
   func MyCarriageLeft(style tcell.Style) *render.Sprite { ... }
   ```

## Adding a Scene

1. Create a new file in `internal/scene/` (e.g. `desert.go`).
2. Implement the `Scene` interface from `internal/scene/scene.go`.
3. Add the scene type to `SceneType` enum in `scene.go`.
4. Add the scene to the `New()` function and `AllSceneTypes()`.
5. Add `RenderBackgroundWithPalette` for time-of-day-aware rendering.
6. Test with `go test -v ./internal/scene/`.

## Sprite Formatting Rules

- All sprites must be **6 rows tall**.
- All frames within a sprite must have identical dimensions.
- Trailing whitespace must be preserved (use blank cells for padding).
- Do not depend on font ligatures.
- Prefer ASCII characters. Unicode may be used only with an ASCII fallback.
- Sprites must remain understandable in monochrome mode.
- Coupling points between locomotives and carriages must align (same height).
- Left-facing and right-facing sprites must be validated.

## Sprite Validation

Run sprite validation to check for structural problems:

```bash
go test -v -run TestAllSpritesValidate ./internal/train/
```

## Snapshot Tests

Snapshot tests generate deterministic rendered frames and compare them
against golden files. To update golden files after making changes:

```bash
UPDATE_SNAPSHOTS=1 go test -v -run TestSnapshot ./internal/tests/
```

Commit the updated golden files along with your changes.

## Testing

### Unit Tests

```bash
go test ./...
```

### Snapshot Tests

```bash
go test -v -run TestSnapshot ./internal/tests/
```

### Integration Tests

```bash
go test -v -run TestIntegration ./internal/tests/
```

### Soak Tests (long running)

```bash
go test -run TestSoakShort -timeout 5m ./internal/tests/
```

## Pull Request Process

1. Ensure all tests pass: `go test ./...`
2. Ensure formatting is clean: `gofmt -w . && go vet ./...`
3. Commit your changes with a clear message.
4. Push to your fork and open a pull request.
5. In the PR description, explain what you changed and why.
6. If adding new ASCII art, confirm that it is original work or properly
   licensed. Add attribution in the asset file header.

## Asset Licensing

All ASCII art in Conductor must be original work. If you contribute sprites
from another source, you must:

- Verify the source has a compatible license (MIT, Apache 2.0, CC0, or public domain).
- Add a comment in the asset file noting the source and license.
- Do not use ASCII art from Asciiquarium or other terminal screensaver projects
  without explicit permission, as their licensing status is often unclear.