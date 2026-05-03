# input API

Keyboard-only API for now.

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

## TBD

See [Roadmap](../roadmap.md#input).
