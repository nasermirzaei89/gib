# input API

Keyboard and mouse input APIs.

## is_key_down(name)

```lua
if input.is_key_down("left") then
  -- held
end
```

Returns true while key is currently down.

## is_key_pressed(name)

```lua
if input.is_key_pressed("space") then
  -- true only on first frame of press
end
```

Edge-triggered press state for current frame.

## is_key_released(name)

```lua
if input.is_key_released("space") then
  -- true only on release frame
end
```

Edge-triggered release state for current frame.

## Key Names

Aliases include common names such as:

- `left`, `right`, `up`, `down`
- `space`, `enter`/`return`, `esc`/`escape`
- `shift`, `ctrl`, `alt`

Unknown key names raise an error.

## get_mouse_position()

```lua
local x, y = input.get_mouse_position()
```

Returns the mouse cursor position as window-local integer pixel coordinates.

## is_mouse_button_down(name)

```lua
if input.is_mouse_button_down("left") then
  -- held
end
```

Returns true while the mouse button is currently down.

## is_mouse_button_pressed(name)

```lua
if input.is_mouse_button_pressed("left") then
  -- true only on first frame of press
end
```

Edge-triggered mouse button press state for current frame.

## is_mouse_button_released(name)

```lua
if input.is_mouse_button_released("left") then
  -- true only on release frame
end
```

Edge-triggered mouse button release state for current frame.

## Mouse Button Names

Supported names:

- `left`
- `middle`
- `right`
- `x1`
- `x2`

Unknown mouse button names raise an error.

## TBD

See [Roadmap](../roadmap.md#input).
