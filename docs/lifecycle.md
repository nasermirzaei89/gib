# Game Lifecycle

The engine creates a global `game` table and calls callbacks if present.

## Callback Order Per Frame

1. Poll and process events.
2. Call `game.event(event)` for each event (if defined).
3. Run fixed-step updates (`game.fixed_update(fixed_dt)`) at 60 TPS while accumulator allows.
4. Call `game.update(dt)` once with variable frame delta.
5. Call `game.render()` once and present the frame.

## Supported Callbacks

```lua
function game.load() end
function game.fixed_update(fixed_dt) end
function game.update(dt) end
function game.render() end
function game.event(event) end
```

All callbacks are optional.

## Timing

- Fixed tick rate: 60 TPS (`fixed_dt = 1/60`).
- Variable update `dt`: real elapsed seconds for current frame.

## Event Payload Today

`game.event(event)` currently receives:

- `event.type` (string)

See [Events](events.md) for mapped names.

## TBD

- Event payload fields beyond `event.type` (key code, mouse position, etc.) are not exposed yet.
- Pause/time-scale control is not implemented yet.
