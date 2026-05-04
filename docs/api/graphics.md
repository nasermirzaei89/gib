# graphics API

## load_image(path)

```lua
local image = graphics.load_image("../assets/leo.png")
```

Loads an image and returns image userdata.

Notes:

- Supports BMP, PNG, and JPEG.
- Relative paths are resolved from script base directory.

The returned image userdata provides this method:

```lua
local width, height = image:get_size()
```

Returns image dimensions in integer pixels.

Calling `image:get_size()` after `graphics.unload_image(image)` raises an error.

## unload_image(image)

```lua
graphics.unload_image(image)
```

Releases texture resources for that image.

## draw_image(image, x, y, opts?)

```lua
graphics.draw_image(image, x, y)

graphics.draw_image(image, x, y, {
  sx = 0, sy = 0, sw = 32, sh = 32,
  scale_x = 1,
  scale_y = 1,
})
```

Draws image at destination position.

Options:

- Source rect: `sx`, `sy`, `sw`, `sh`
- Scaling: `scale_x`, `scale_y` (default `1`)
- Negative scale mirrors on that axis.
- `scale_x == 0` or `scale_y == 0` draws nothing.

Validation behavior:

- If any source key is present, all `sx/sy/sw/sh` are required.
- Source rect must be in bounds and use integer values.

## TBD

See [Roadmap](../roadmap.md#graphics).
