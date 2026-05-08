# Game Lifecycle

The engine creates a global `game` table and calls callbacks if present.

## Startup Callback

Before initialization, the engine calls `game.config(conf)` if present.

```lua
function game.config(conf)
    conf.window.width = 1280
    conf.window.height = 720
    conf.window.title = "My Game"
    conf.window.resizable = true
    conf.window.fullscreen = false
    conf.graphics.auto_clear = true
    conf.graphics.clear_color = {0.0, 0.0, 0.0, 1.0}
    conf.tps = 60
end
```

`conf` defaults:

- `conf.window.title = "Game"`
- `conf.window.width = 800`
- `conf.window.height = 600`
- `conf.window.resizable = false`
- `conf.window.fullscreen = false`
- `conf.graphics.auto_clear = true`
- `conf.graphics.clear_color = {0.0, 0.0, 0.0, 1.0}`
- `conf.tps = 60`

Validation rules:

- `conf.tps` must be a number greater than 0.
- `conf.window.width` and `conf.window.height` must be positive integers.
- `conf.window.title` must be a string (empty string allowed).
- `conf.window.resizable` and `conf.window.fullscreen` must be booleans.
- `conf.graphics.auto_clear` must be a boolean.
- `conf.graphics.clear_color` must be `{r, g, b, a}` with values in `[0, 1]`.

## Callback Order Per Frame

1. Poll and process events.
2. Call `game.event(event)` for each event (if defined).
3. Run fixed-step updates (`game.fixed_update(fixed_dt)`) at 60 TPS while accumulator allows.
4. Call `game.update(dt)` once with variable frame delta.
5. Auto-clear the frame (if enabled).
6. Call `game.render()` once and present the frame.

## Supported Callbacks

```lua
function game.config(conf) end
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

## Event Payload

`game.event(event)` receives `event.type` plus event-specific fields.

See [Events](events.md) for the full payload contract.

## TBD

- Pause/time-scale control is not implemented yet.
