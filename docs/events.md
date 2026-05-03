# Events

Events are delivered to:

```lua
function game.event(event)
    -- event.type is a string
end
```

## Payload

Currently exposed fields:

- `event.type`

No other event metadata is exposed yet.

## Event Name Examples

The engine maps many SDL3 events to snake_case names. Common examples:

- System: `quit`, `low_memory`, `locale_changed`, `system_theme_changed`
- Window: `window_shown`, `window_resized`, `window_focus_gained`, `window_focus_lost`, `window_close_requested`
- Keyboard: `key_down`, `key_up`, `keyboard_added`, `keyboard_removed`
- Mouse: `mouse_motion`, `mouse_button_down`, `mouse_button_up`, `mouse_wheel`
- Drag and drop: `drop_file`, `drop_text`, `drop_begin`, `drop_complete`
- Render/device: `render_device_reset`, `render_device_lost`

For the full current mapping, see source in `event.go`.

## TBD

- Expose detailed event fields (`key`, `scancode`, mouse coordinates, wheel deltas).
- Add optional event filtering helpers.
