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

## draw_rect(x, y, w, h, opts?)

```lua
graphics.draw_rect(50, 50, 120, 80)

graphics.draw_rect(50, 50, 120, 80, {
  color = {0.2, 0.7, 1.0, 1.0},
  filled = true,
})
```

Draws a rectangle.

- `filled = true` uses filled rendering.
- `filled = false` (default) draws outline.

Options:

- `color = {r, g, b, a}` where each value is in `[0, 1]`.
- `filled` boolean (default `false`).

Validation behavior:

- `w` and `h` must be `> 0`.
- If provided, `color` must contain exactly 4 numeric values in `[0, 1]`.

## draw_line(x1, y1, x2, y2, opts?)

```lua
graphics.draw_line(40, 40, 300, 120)

graphics.draw_line(40, 40, 300, 120, {
  color = {1.0, 0.3, 0.3, 1.0},
})
```

Draws a single line segment.

Options:

- `color = {r, g, b, a}` where each value is in `[0, 1]`.

Validation behavior:

- If provided, `color` must contain exactly 4 numeric values in `[0, 1]`.

## draw_polygon(points, opts?)

```lua
graphics.draw_polygon({
  {80, 340},
  {160, 280},
  {260, 340},
})

graphics.draw_polygon({
  {320, 340},
  {420, 260},
  {520, 340},
  {470, 420},
  {350, 420},
}, {
  color = {0.8, 0.9, 0.2, 1.0},
  filled = true,
})
```

Draws a polygon from point pairs.

Points format:

- `points = {{x, y}, {x, y}, ...}`

Options:

- `color = {r, g, b, a}` where each value is in `[0, 1]`.
- `filled` boolean (default `false`).
- `closed` boolean (default `true`) for outline mode.

Validation behavior:

- `points` must contain at least 3 points.
- Each point must contain exactly two numeric values.
- If provided, `color` must contain exactly 4 numeric values in `[0, 1]`.
- Filled polygons currently require convex points.
- Filled polygons require `closed = true`.

## draw_circle(x, y, r, opts?)

```lua
graphics.draw_circle(120, 500, 36)

graphics.draw_circle(240, 500, 36, {
  color = {0.2, 0.8, 1.0, 1.0},
  filled = true,
  segments = 48,
})
```

Draws a circle.

Options:

- `color = {r, g, b, a}` where each value is in `[0, 1]`.
- `filled` boolean (default `false`).
- `segments` integer (default `48`, minimum `3`) controlling smoothness.

Validation behavior:

- `r` must be `> 0`.
- If provided, `color` must contain exactly 4 numeric values in `[0, 1]`.
- If provided, `segments` must be an integer `>= 3`.

## draw_ellipse(x, y, rx, ry, opts?)

```lua
graphics.draw_ellipse(380, 500, 56, 28)

graphics.draw_ellipse(520, 500, 56, 28, {
  color = {1.0, 0.5, 0.2, 1.0},
  filled = true,
  segments = 56,
})
```

Draws an ellipse.

Options:

- `color = {r, g, b, a}` where each value is in `[0, 1]`.
- `filled` boolean (default `false`).
- `segments` integer (default `48`, minimum `3`) controlling smoothness.

Validation behavior:

- `rx` and `ry` must be `> 0`.
- If provided, `color` must contain exactly 4 numeric values in `[0, 1]`.
- If provided, `segments` must be an integer `>= 3`.

## draw_arc(x, y, r, start_angle, end_angle, opts?)

```lua
graphics.draw_arc(680, 500, 44, 0.0, math.pi * 1.5)

graphics.draw_arc(680, 500, 32, math.pi * 1.25, math.pi * 0.5, {
  color = {1.0, 0.3, 0.4, 1.0},
  segments = 24,
})
```

Draws a circular arc using radian angles.

Arc semantics:

- Draw direction is counter-clockwise (CCW).
- If `end_angle < start_angle`, the arc wraps across `2π`.

Options:

- `color = {r, g, b, a}` where each value is in `[0, 1]`.
- `segments` integer (default `32`, minimum `1`) controlling smoothness.

Validation behavior:

- `r` must be `> 0`.
- Filled arcs are currently unsupported.
- If provided, `color` must contain exactly 4 numeric values in `[0, 1]`.
- If provided, `segments` must be an integer `>= 1`.

## TBD

See [Roadmap](../roadmap.md#graphics).
