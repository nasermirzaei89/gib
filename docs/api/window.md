# window API

## get_size()

```lua
local w, h = window.get_size()
```

Returns current window width and height.

## set_size(w, h)

```lua
window.set_size(1024, 768)
```

Sets window size. `w` and `h` must be positive integers.

## set_title(title)

```lua
window.set_title("My Game")
```

Sets window title.

## set_fullscreen(enabled)

```lua
window.set_fullscreen(true)
```

Toggles fullscreen mode.

## is_fullscreen()

```lua
local full = window.is_fullscreen()
```

Returns fullscreen state.

## close()

```lua
window.close()
```

Requests engine loop shutdown.

## TBD

See [Roadmap](../roadmap.md#engine--gameplay).
