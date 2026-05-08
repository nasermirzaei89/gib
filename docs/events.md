# Events

Events are delivered to:

```lua
function game.event(event)
    -- event.type is a string
end
```

## Payload

Common fields:

- `event.type`

Event-specific fields:

- `key_down`
  - `event.key`: key name in GIB canonical format (examples: `left`, `enter`, `left_shift`, `a`)
  - `event.is_repeat`: `true` when the key-down is an OS key-repeat event
- `key_up`
  - `event.key`: key name in GIB canonical format (examples: `left`, `enter`, `left_shift`, `a`)
- `mouse_motion`
  - `event.x`: mouse x position in window coordinates
  - `event.y`: mouse y position in window coordinates
  - `event.dx`: relative x motion since previous mouse motion event
  - `event.dy`: relative y motion since previous mouse motion event
- `mouse_button_down`
  - `event.button`: canonical button name (`left`, `middle`, `right`, `x1`, `x2`)
  - `event.clicks`: click count (1 for single click, 2 for double click)
  - `event.x`: mouse x position in window coordinates
  - `event.y`: mouse y position in window coordinates
- `mouse_button_up`
  - `event.button`: canonical button name (`left`, `middle`, `right`, `x1`, `x2`)
  - `event.clicks`: click count (usually 1)
  - `event.x`: mouse x position in window coordinates
  - `event.y`: mouse y position in window coordinates
- `mouse_wheel`
  - `event.x`: horizontal wheel delta (positive is right)
  - `event.y`: vertical wheel delta (positive is away from user)
  - `event.mouse_x`: mouse x position in window coordinates during wheel event
  - `event.mouse_y`: mouse y position in window coordinates during wheel event

Other event metadata is intentionally not exposed yet.

## Event Name Examples

The engine maps runtime events to snake_case names. Common examples:

- System: `quit`, `low_memory`, `locale_changed`, `system_theme_changed`
- Window: `window_shown`, `window_resized`, `window_focus_gained`, `window_focus_lost`, `window_close_requested`
- Keyboard: `key_down`, `key_up`, `keyboard_added`, `keyboard_removed`
- Mouse: `mouse_motion`, `mouse_button_down`, `mouse_button_up`, `mouse_wheel`
- Drag and drop: `drop_file`, `drop_text`, `drop_begin`, `drop_complete`
- Render/device: `render_device_reset`, `render_device_lost`

For the full current mapping, see source in `event.go`.

## TBD

- Expose detailed fields for more non-keyboard and non-mouse event types.
- Add optional event filtering helpers.
