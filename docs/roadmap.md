# Roadmap

This page tracks the planned direction and priorities for Gib.

Gib focuses on a small, clean Lua API with a native runtime backend. The goal is to keep iteration fast while avoiding unnecessary engine complexity.

---

## Current Baseline

The following systems are already available:

- Graphics:
  image loading, image drawing with scaling/source rects, and basic shape rendering.

- Input:
  keyboard and mouse polling with edge-triggered pressed/released helpers.

- Audio:
  asset loading, playback instances, pause/resume/stop, per-instance volume, master volume.

- Window:
  size/title/fullscreen controls and close requests.

- Lifecycle:
  `config`, `load`, `fixed_update`, `update`, `render`, and `event` callbacks.

---

## Priority 1: CLI and Development Workflow

The CLI should remain small and workflow-oriented.

### Core Commands

```bash
gib run [game-dir]
gib build [game-dir] --target <web|windows|macos|linux>
````

#### `gib run`

Run games locally for desktop iteration.

#### `gib build`

Export builds for supported targets.

Initial web output baseline:

```text
dist/
  index.html
  game.wasm
  game.data
```

(or equivalent JS runtime output)

### Packaging

Optional release packaging helpers:

```bash
gib package
```

Used for:

- zip archives
- app bundles
- installers

This may later merge into `build`.

### Web Runtime Requirements

Web support requires:

- WebAssembly runtime support
- browser input/event backend
- browser-compatible asset loading

### Not Planned

The following are intentionally out of scope for now:

- plugin system
- package manager
- cloud publishing
- editor tooling
- account/login systems

---

## Priority 2: Core Runtime Features

These systems unlock a large amount of gameplay and UI functionality.

### Event System Expansion

Richer `game.event(event)` payloads:

- canonical key data
- mouse position/delta/button data
- wheel deltas
- window events
- text input payloads
- file drop payloads
- device identifiers

### Text Input

Production-ready text input support:

- composition-aware input
- IME support
- text editing helpers
- UI/chat/name-entry workflows

### Gamepad Support

- connected device management
- button/axis polling
- deadzone helpers
- hot-plug handling

### Save/Load Utilities

Simple persistence helpers:

```lua
save.write_json(...)
save.read_json(...)
```

### Debugging and Error Handling

Developer-facing runtime diagnostics:

- Lua stack trace overlay
- runtime error screen
- debug text helpers
- FPS/debug metrics

### Asset Hot Reload

Development-time asset reload support for:

- images
- shaders
- audio
- Lua scripts where possible

---

## Priority 3: Rendering Foundations

These systems reduce repetitive rendering and camera scaffolding.

### Background Clear API

Implemented baseline:

- engine clears every frame by default
- startup config supports `conf.graphics.auto_clear` and `conf.graphics.clear_color`
- runtime controls: `graphics.clear()`, `graphics.set_clear_color(...)`, `graphics.set_auto_clear(...)`

### Image Transforms

Unified transform options are now available through `opts` for:

- `graphics.draw_image`
- `graphics.draw_rect`
- `graphics.draw_ellipse`
- `graphics.draw_arc`

Transform fields:

- rotation
- origin/pivot
- scale (including mirroring with negative scale)

Notes:

- Rotation uses radians.
- Scale and rotation share one pivot (`origin`).
- Transform order is `scale -> rotate -> translate`.
- `draw_polygon` transform support is deferred to a follow-up phase.

### Camera and Viewport Helpers

- camera position
- zoom
- world-to-screen transforms
- viewport/scissor utilities

### Text Rendering

Production-ready text rendering:

- font loading
- font sizing
- alignment
- multiline layout
- color styling

---

## Priority 4: Project Structure Helpers

These systems improve maintainability as projects grow.

### Scene and State Helpers

Optional helpers for:

- scene transitions
- stack/switch flows
- lifecycle coordination

### Gameplay Organization Helpers

Lightweight helpers for:

- update/render grouping
- ownership patterns
- entity collections

Gib does not plan to enforce a mandatory ECS architecture.

### Time Controls

- pause helpers
- time scaling
- slow motion support

---

## Priority 5: Multiplayer Enablement

Lua networking libraries are already usable externally, but engine-level support may later improve multiplayer workflows.

### Potential Future Additions

- optional networking wrappers/helpers
- deterministic tick helpers
- sync-oriented utilities
- serialization helpers
- sample multiplayer architectures

---

## Documentation Priorities

## Examples and Recipes

Focused examples for:

- animation systems
- cameras
- UI widgets
- save systems
- scene flow
- gamepad usage

## Project Structure Guidance

Recommended Lua project organization for:

- small projects
- medium projects
- larger productions

## Contributor Documentation

Separate contributor-focused documentation for:

- engine architecture
- backend systems
- rendering internals
- platform/runtime details

---

## Non-Goals

Gib intentionally avoids:

- built-in visual editors
- visual scripting systems
- mandatory ECS architecture
- large framework abstractions
- engine-managed gameplay structure
- cloud/account ecosystems
