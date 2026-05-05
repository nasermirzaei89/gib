# Examples

## hello-world

Path: `examples/hello-world/main.lua`

Shows the minimal `game.render()` flow with `debug.print`.

## render-image

Path: `examples/render-image/main.lua`

Shows image loading with `graphics.load_image` and drawing with `graphics.draw_image`.

## draw-shapes

Path: `examples/draw-shapes/main.lua`

Shows primitive shape rendering using:

- `graphics.draw_rect`
- `graphics.draw_line`

## animation

Path: `examples/animation/main.lua`

Shows sprite-sheet animation by changing source rect (`sx/sy/sw/sh`) over time.

## input-keyboard

Path: `examples/input-keyboard/main.lua`

Shows keyboard movement using:

- `input.is_key_down`
- source-rect drawing for animated sprites

## input-mouse

Path: `examples/input-mouse/main.lua`

Shows mouse position and edge-triggered button input using:

- `input.get_mouse_position`
- `input.is_mouse_button_down`
- `input.is_mouse_button_pressed`
- `input.is_mouse_button_released`

## events

Path: `examples/events/main.lua`

Shows `game.event(event)` usage and logs each incoming event type.

## TBD

- Add an example for `is_key_pressed` / `is_key_released` edge-triggered input.
- Add an example for scaled/mirrored `graphics.draw_image` options.
