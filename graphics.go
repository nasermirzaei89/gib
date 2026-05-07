package gib

import (
	"image"
	"image/draw"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
	lua "github.com/yuin/gopher-lua"

	_ "image/jpeg"
	_ "image/png"
)

const imageMetatableName = "graphics.image"

type Image struct {
	texture  *sdl.Texture
	width    int32
	height   int32
	released bool
}

type shapeOptions struct {
	color  sdl.Color
	filled bool
	closed bool
}

const (
	defaultCircleSegments  = 48
	defaultEllipseSegments = 48
	defaultArcSegments     = 32
)

func checkIntegerField(L *lua.LState, tbl *lua.LTable, key string) int32 {
	val := tbl.RawGetString(key)
	if val.Type() == lua.LTNil {
		L.RaiseError("graphics.draw_image: source rect key %q is required", key)
		return 0
	}

	num, ok := val.(lua.LNumber)
	if !ok {
		L.RaiseError("graphics.draw_image: source rect key %q must be a number", key)
		return 0
	}

	if lua.LNumber(int64(num)) != num {
		L.RaiseError("graphics.draw_image: source rect key %q must be an integer", key)
		return 0
	}

	return int32(int64(num))
}

func getOptionalNumberField(L *lua.LState, tbl *lua.LTable, key string, defaultVal float32) (float32, bool) {
	val := tbl.RawGetString(key)
	if val.Type() == lua.LTNil {
		return defaultVal, false
	}

	num, ok := val.(lua.LNumber)
	if !ok {
		L.RaiseError("graphics.draw_image: option key %q must be a number", key)
		return 0, true
	}

	return float32(num), true
}

func checkLuaImage(L *lua.LState, index int) *Image {
	ud := L.CheckUserData(index)
	img, ok := ud.Value.(*Image)
	if !ok || img == nil {
		L.ArgError(index, "expected image userdata")
		return nil
	}

	return img
}

func releaseLuaImage(img *Image) {
	if img == nil || img.released {
		return
	}

	if img.texture != nil {
		img.texture.Destroy()
		img.texture = nil
	}

	img.released = true
}

func normalizeColorChannel(component float32) uint8 {
	if component <= 0 {
		return 0
	}
	if component >= 1 {
		return 255
	}
	return uint8(math.Round(float64(component * 255)))
}

func parseShapeOptions(L *lua.LState, fnName string, arg lua.LValue, allowFilled bool) shapeOptions {
	opts := shapeOptions{
		color:  sdl.Color{R: 255, G: 255, B: 255, A: 255},
		closed: true,
	}

	if arg.Type() == lua.LTNil {
		return opts
	}

	tbl, ok := arg.(*lua.LTable)
	if !ok {
		L.RaiseError("%s: options argument must be a table", fnName)
		return opts
	}

	if allowFilled {
		filledVal := tbl.RawGetString("filled")
		if filledVal.Type() != lua.LTNil {
			filledBool, ok := filledVal.(lua.LBool)
			if !ok {
				L.RaiseError("%s: options key %q must be a boolean", fnName, "filled")
				return opts
			}
			opts.filled = bool(filledBool)
		}
	}

	closedVal := tbl.RawGetString("closed")
	if closedVal.Type() != lua.LTNil {
		closedBool, ok := closedVal.(lua.LBool)
		if !ok {
			L.RaiseError("%s: options key %q must be a boolean", fnName, "closed")
			return opts
		}
		opts.closed = bool(closedBool)
	}

	colorVal := tbl.RawGetString("color")
	if colorVal.Type() == lua.LTNil {
		return opts
	}

	colorTbl, ok := colorVal.(*lua.LTable)
	if !ok {
		L.RaiseError("%s: options key %q must be a table {r, g, b, a}", fnName, "color")
		return opts
	}

	if colorTbl.Len() != 4 {
		L.RaiseError("%s: options key %q must contain exactly 4 values (r, g, b, a)", fnName, "color")
		return opts
	}

	components := [4]float32{}
	for i := 1; i <= 4; i++ {
		componentVal := colorTbl.RawGetInt(i)
		componentNum, ok := componentVal.(lua.LNumber)
		if !ok {
			L.RaiseError("%s: options key %q value at index %d must be a number", fnName, "color", i)
			return opts
		}

		component := float32(componentNum)
		if component < 0 || component > 1 {
			L.RaiseError("%s: options key %q values must be in range [0, 1]", fnName, "color")
			return opts
		}

		components[i-1] = component
	}

	opts.color = sdl.Color{
		R: normalizeColorChannel(components[0]),
		G: normalizeColorChannel(components[1]),
		B: normalizeColorChannel(components[2]),
		A: normalizeColorChannel(components[3]),
	}

	return opts
}

func renderWithColor(renderer *sdl.Renderer, color sdl.Color, draw func() error) error {
	previousColor, err := renderer.DrawColor()
	if err != nil {
		return err
	}

	if err := renderer.SetDrawColor(color.R, color.G, color.B, color.A); err != nil {
		return err
	}

	drawErr := draw()
	restoreErr := renderer.SetDrawColor(previousColor.R, previousColor.G, previousColor.B, previousColor.A)

	if drawErr != nil {
		return drawErr
	}

	if restoreErr != nil {
		return restoreErr
	}

	return nil
}

func parseSegmentsOption(L *lua.LState, fnName string, arg lua.LValue, defaultSegments, minSegments int) int {
	if arg.Type() == lua.LTNil {
		return defaultSegments
	}

	tbl, ok := arg.(*lua.LTable)
	if !ok {
		L.RaiseError("%s: options argument must be a table", fnName)
		return defaultSegments
	}

	segmentsVal := tbl.RawGetString("segments")
	if segmentsVal.Type() == lua.LTNil {
		return defaultSegments
	}

	segmentsNum, ok := segmentsVal.(lua.LNumber)
	if !ok {
		L.RaiseError("%s: options key %q must be a number", fnName, "segments")
		return defaultSegments
	}

	if lua.LNumber(int64(segmentsNum)) != segmentsNum {
		L.RaiseError("%s: options key %q must be an integer", fnName, "segments")
		return defaultSegments
	}

	segments := int(segmentsNum)
	if segments < minSegments {
		L.RaiseError("%s: options key %q must be >= %d", fnName, "segments", minSegments)
		return defaultSegments
	}

	return segments
}

func ellipsePoint(cx, cy, rx, ry float32, angle float64) sdl.FPoint {
	return sdl.FPoint{
		X: cx + rx*float32(math.Cos(angle)),
		Y: cy + ry*float32(math.Sin(angle)),
	}
}

func buildEllipsePoints(cx, cy, rx, ry float32, segments int) []sdl.FPoint {
	points := make([]sdl.FPoint, 0, segments)
	for i := range segments {
		angle := 2 * math.Pi * float64(i) / float64(segments)
		points = append(points, ellipsePoint(cx, cy, rx, ry, angle))
	}

	return points
}

func normalizeArc(startAngle, endAngle float64) (float64, float64) {
	span := math.Mod(endAngle-startAngle, 2*math.Pi)
	if span < 0 {
		span += 2 * math.Pi
	}

	return startAngle, span
}

func buildArcPoints(cx, cy, r float32, startAngle, endAngle float64, segments int) []sdl.FPoint {
	start, span := normalizeArc(startAngle, endAngle)
	if span == 0 {
		return nil
	}

	points := make([]sdl.FPoint, 0, segments+1)
	for i := 0; i <= segments; i++ {
		t := float64(i) / float64(segments)
		angle := start + span*t
		points = append(points, ellipsePoint(cx, cy, r, r, angle))
	}

	return points
}

func parsePolygonPoints(L *lua.LState, pointsVal lua.LValue) []sdl.FPoint {
	pointsTbl, ok := pointsVal.(*lua.LTable)
	if !ok {
		L.RaiseError("graphics.draw_polygon: points must be a table of points {{x, y}, ...}")
		return nil
	}

	pointCount := pointsTbl.Len()
	if pointCount < 3 {
		L.RaiseError("graphics.draw_polygon: points must contain at least 3 points")
		return nil
	}

	points := make([]sdl.FPoint, 0, pointCount)
	for i := 1; i <= pointCount; i++ {
		pointVal := pointsTbl.RawGetInt(i)
		pointTbl, ok := pointVal.(*lua.LTable)
		if !ok {
			L.RaiseError("graphics.draw_polygon: point at index %d must be a table {x, y}", i)
			return nil
		}

		if pointTbl.Len() != 2 {
			L.RaiseError("graphics.draw_polygon: point at index %d must contain exactly 2 numeric values", i)
			return nil
		}

		xVal := pointTbl.RawGetInt(1)
		yVal := pointTbl.RawGetInt(2)

		xNum, ok := xVal.(lua.LNumber)
		if !ok {
			L.RaiseError("graphics.draw_polygon: point at index %d x must be a number", i)
			return nil
		}
		yNum, ok := yVal.(lua.LNumber)
		if !ok {
			L.RaiseError("graphics.draw_polygon: point at index %d y must be a number", i)
			return nil
		}

		points = append(points, sdl.FPoint{X: float32(xNum), Y: float32(yNum)})
	}

	return points
}

func polygonOrientation(a, b, c sdl.FPoint) float32 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

func isConvexPolygon(points []sdl.FPoint) bool {
	if len(points) < 3 {
		return false
	}

	hasPositive := false
	hasNegative := false
	const epsilon = 1e-6

	for i := range points {
		a := points[i]
		b := points[(i+1)%len(points)]
		c := points[(i+2)%len(points)]

		cross := polygonOrientation(a, b, c)
		if math.Abs(float64(cross)) <= epsilon {
			continue
		}

		if cross > 0 {
			hasPositive = true
		} else {
			hasNegative = true
		}

		if hasPositive && hasNegative {
			return false
		}
	}

	return true
}

func polygonVertices(points []sdl.FPoint, color sdl.Color) []sdl.Vertex {
	vertices := make([]sdl.Vertex, 0, len(points))
	vertexColor := sdl.FColor{
		R: float32(color.R) / 255,
		G: float32(color.G) / 255,
		B: float32(color.B) / 255,
		A: float32(color.A) / 255,
	}

	for _, point := range points {
		vertices = append(vertices, sdl.Vertex{
			Position: point,
			Color:    vertexColor,
			TexCoord: sdl.FPoint{},
		})
	}

	return vertices
}

func polygonIndices(pointCount int) []int32 {
	indices := make([]int32, 0, (pointCount-2)*3)
	for i := 1; i < pointCount-1; i++ {
		indices = append(indices, 0, int32(i), int32(i+1))
	}
	return indices
}

func drawPolygonOutline(renderer *sdl.Renderer, points []sdl.FPoint, closed bool) error {
	if len(points) < 2 {
		return nil
	}

	for i := 1; i < len(points); i++ {
		prev := points[i-1]
		curr := points[i]
		if err := renderer.RenderLine(prev.X, prev.Y, curr.X, curr.Y); err != nil {
			return err
		}
	}

	if closed {
		first := points[0]
		last := points[len(points)-1]
		if err := renderer.RenderLine(last.X, last.Y, first.X, first.Y); err != nil {
			return err
		}
	}

	return nil
}

// InitGraphics registers the graphics module and its functions to the Lua state.
// Graphics module is responsible rendering.
func InitGraphics(l *lua.LState, renderer *sdl.Renderer, baseDir string) {
	mt := l.NewTypeMetatable(imageMetatableName)
	imageMethods := l.NewTable()
	imageMethods.RawSetString("get_size", l.NewFunction(func(L *lua.LState) int {
		img := checkLuaImage(L, 1)
		if img.released || img.texture == nil {
			L.RaiseError("graphics.image:get_size: image has been unloaded")
			return 0
		}

		L.Push(lua.LNumber(img.width))
		L.Push(lua.LNumber(img.height))
		return 2
	}))
	mt.RawSetString("__index", imageMethods)
	mt.RawSetString("__gc", l.NewFunction(func(L *lua.LState) int {
		img := checkLuaImage(L, 1)
		releaseLuaImage(img)
		return 0
	}))

	graphics := l.NewTable()
	l.SetGlobal("graphics", graphics)

	graphics.RawSetString("load_image", l.NewFunction(func(L *lua.LState) int {
		fp := L.CheckString(1)

		if !filepath.IsAbs(fp) {
			fp = filepath.Join(baseDir, fp)
		}

		texture, err := loadTexture(renderer, fp)
		if err != nil {
			L.RaiseError("graphics.load_image(%q) failed: %v", fp, err)
			return 0
		}

		if texture == nil {
			L.RaiseError("graphics.load_image(%q) failed: texture is nil", fp)
			return 0
		}

		if texture.W <= 0 || texture.H <= 0 {
			texture.Destroy()
			L.RaiseError("graphics.load_image(%q) failed: invalid texture size %dx%d", fp, texture.W, texture.H)
			return 0
		}

		ud := L.NewUserData()
		ud.Value = &Image{texture: texture, width: texture.W, height: texture.H}
		L.SetMetatable(ud, L.GetTypeMetatable(imageMetatableName))
		L.Push(ud)

		return 1
	}))

	graphics.RawSetString("unload_image", l.NewFunction(func(L *lua.LState) int {
		img := checkLuaImage(L, 1)
		releaseLuaImage(img)

		return 0
	}))

	graphics.RawSetString("draw_image", l.NewFunction(func(L *lua.LState) int {
		img := checkLuaImage(L, 1)
		x := float32(L.CheckNumber(2))
		y := float32(L.CheckNumber(3))

		if img.released || img.texture == nil {
			L.RaiseError("graphics.draw_image: image has been unloaded")
			return 0
		}

		var src *sdl.FRect
		dstW := float32(img.texture.W)
		dstH := float32(img.texture.H)
		scaleX := float32(1)
		scaleY := float32(1)
		flipMode := sdl.FLIP_NONE

		if L.GetTop() >= 4 {
			arg4 := L.Get(4)
			tbl, ok := arg4.(*lua.LTable)
			if !ok {
				L.RaiseError("graphics.draw_image: 4th argument must be a table with draw options")
				return 0
			}

			scaleXVal, _ := getOptionalNumberField(L, tbl, "scale_x", 1)
			scaleYVal, _ := getOptionalNumberField(L, tbl, "scale_y", 1)
			scaleX = scaleXVal
			scaleY = scaleYVal

			hasSX := tbl.RawGetString("sx").Type() != lua.LTNil
			hasSY := tbl.RawGetString("sy").Type() != lua.LTNil
			hasSW := tbl.RawGetString("sw").Type() != lua.LTNil
			hasSH := tbl.RawGetString("sh").Type() != lua.LTNil
			hasAnySourceKey := hasSX || hasSY || hasSW || hasSH

			if hasAnySourceKey {
				if !(hasSX && hasSY && hasSW && hasSH) {
					L.RaiseError("graphics.draw_image: source rect requires all keys sx, sy, sw, sh when any is provided")
					return 0
				}

				sx := checkIntegerField(L, tbl, "sx")
				sy := checkIntegerField(L, tbl, "sy")
				sw := checkIntegerField(L, tbl, "sw")
				sh := checkIntegerField(L, tbl, "sh")

				if sx < 0 || sy < 0 {
					L.RaiseError("graphics.draw_image: source rect sx/sy must be >= 0")
					return 0
				}

				if sw <= 0 || sh <= 0 {
					L.RaiseError("graphics.draw_image: source rect sw/sh must be > 0")
					return 0
				}

				if sx+sw > img.texture.W || sy+sh > img.texture.H {
					L.RaiseError(
						"graphics.draw_image: source rect (%d,%d,%d,%d) is out of bounds for image size %dx%d",
						sx,
						sy,
						sw,
						sh,
						img.texture.W,
						img.texture.H,
					)
					return 0
				}

				src = &sdl.FRect{
					X: float32(sx),
					Y: float32(sy),
					W: float32(sw),
					H: float32(sh),
				}

				dstW = float32(sw)
				dstH = float32(sh)
			}
		}

		if scaleX < 0 {
			flipMode |= sdl.FLIP_HORIZONTAL
		}
		if scaleY < 0 {
			flipMode |= sdl.FLIP_VERTICAL
		}

		dstW *= float32(math.Abs(float64(scaleX)))
		dstH *= float32(math.Abs(float64(scaleY)))

		if dstW == 0 || dstH == 0 {
			return 0
		}

		dst := sdl.FRect{
			X: x,
			Y: y,
			W: dstW,
			H: dstH,
		}

		err := renderer.RenderTextureRotated(img.texture, src, &dst, 0, nil, flipMode)
		if err != nil {
			L.RaiseError("graphics.draw_image failed: %v", err)
			return 0
		}

		return 0
	}))

	graphics.RawSetString("draw_rect", l.NewFunction(func(L *lua.LState) int {
		x := float32(L.CheckNumber(1))
		y := float32(L.CheckNumber(2))
		w := float32(L.CheckNumber(3))
		h := float32(L.CheckNumber(4))

		if w <= 0 || h <= 0 {
			L.RaiseError("graphics.draw_rect: width and height must be > 0")
			return 0
		}

		optsArg := lua.LValue(lua.LNil)
		if L.GetTop() >= 5 {
			optsArg = L.Get(5)
		}

		opts := parseShapeOptions(L, "graphics.draw_rect", optsArg, true)

		rect := sdl.FRect{X: x, Y: y, W: w, H: h}
		err := renderWithColor(renderer, opts.color, func() error {
			if opts.filled {
				return renderer.RenderFillRect(&rect)
			}
			return renderer.RenderRect(&rect)
		})
		if err != nil {
			L.RaiseError("graphics.draw_rect failed: %v", err)
			return 0
		}

		return 0
	}))

	graphics.RawSetString("draw_line", l.NewFunction(func(L *lua.LState) int {
		x1 := float32(L.CheckNumber(1))
		y1 := float32(L.CheckNumber(2))
		x2 := float32(L.CheckNumber(3))
		y2 := float32(L.CheckNumber(4))

		optsArg := lua.LValue(lua.LNil)
		if L.GetTop() >= 5 {
			optsArg = L.Get(5)
		}

		opts := parseShapeOptions(L, "graphics.draw_line", optsArg, false)

		err := renderWithColor(renderer, opts.color, func() error {
			return renderer.RenderLine(x1, y1, x2, y2)
		})
		if err != nil {
			L.RaiseError("graphics.draw_line failed: %v", err)
			return 0
		}

		return 0
	}))

	graphics.RawSetString("draw_polygon", l.NewFunction(func(L *lua.LState) int {
		points := parsePolygonPoints(L, L.CheckAny(1))

		optsArg := lua.LValue(lua.LNil)
		if L.GetTop() >= 2 {
			optsArg = L.Get(2)
		}

		opts := parseShapeOptions(L, "graphics.draw_polygon", optsArg, true)

		if opts.filled {
			if !opts.closed {
				L.RaiseError("graphics.draw_polygon: filled polygons require opts.closed = true")
				return 0
			}

			if !isConvexPolygon(points) {
				L.RaiseError("graphics.draw_polygon: filled polygons currently require convex points")
				return 0
			}
		}

		err := renderWithColor(renderer, opts.color, func() error {
			if opts.filled {
				vertices := polygonVertices(points, opts.color)
				indices := polygonIndices(len(points))
				return renderer.RenderGeometry(nil, vertices, indices)
			}

			return drawPolygonOutline(renderer, points, opts.closed)
		})
		if err != nil {
			L.RaiseError("graphics.draw_polygon failed: %v", err)
			return 0
		}

		return 0
	}))

	graphics.RawSetString("draw_circle", l.NewFunction(func(L *lua.LState) int {
		x := float32(L.CheckNumber(1))
		y := float32(L.CheckNumber(2))
		r := float32(L.CheckNumber(3))

		if r <= 0 {
			L.RaiseError("graphics.draw_circle: radius must be > 0")
			return 0
		}

		optsArg := lua.LValue(lua.LNil)
		if L.GetTop() >= 4 {
			optsArg = L.Get(4)
		}

		opts := parseShapeOptions(L, "graphics.draw_circle", optsArg, true)
		segments := parseSegmentsOption(L, "graphics.draw_circle", optsArg, defaultCircleSegments, 3)
		points := buildEllipsePoints(x, y, r, r, segments)

		err := renderWithColor(renderer, opts.color, func() error {
			if opts.filled {
				vertices := polygonVertices(points, opts.color)
				indices := polygonIndices(len(points))
				return renderer.RenderGeometry(nil, vertices, indices)
			}

			return drawPolygonOutline(renderer, points, true)
		})
		if err != nil {
			L.RaiseError("graphics.draw_circle failed: %v", err)
			return 0
		}

		return 0
	}))

	graphics.RawSetString("draw_ellipse", l.NewFunction(func(L *lua.LState) int {
		x := float32(L.CheckNumber(1))
		y := float32(L.CheckNumber(2))
		rx := float32(L.CheckNumber(3))
		ry := float32(L.CheckNumber(4))

		if rx <= 0 || ry <= 0 {
			L.RaiseError("graphics.draw_ellipse: rx and ry must be > 0")
			return 0
		}

		optsArg := lua.LValue(lua.LNil)
		if L.GetTop() >= 5 {
			optsArg = L.Get(5)
		}

		opts := parseShapeOptions(L, "graphics.draw_ellipse", optsArg, true)
		segments := parseSegmentsOption(L, "graphics.draw_ellipse", optsArg, defaultEllipseSegments, 3)
		points := buildEllipsePoints(x, y, rx, ry, segments)

		err := renderWithColor(renderer, opts.color, func() error {
			if opts.filled {
				vertices := polygonVertices(points, opts.color)
				indices := polygonIndices(len(points))
				return renderer.RenderGeometry(nil, vertices, indices)
			}

			return drawPolygonOutline(renderer, points, true)
		})
		if err != nil {
			L.RaiseError("graphics.draw_ellipse failed: %v", err)
			return 0
		}

		return 0
	}))

	graphics.RawSetString("draw_arc", l.NewFunction(func(L *lua.LState) int {
		x := float32(L.CheckNumber(1))
		y := float32(L.CheckNumber(2))
		r := float32(L.CheckNumber(3))
		startAngle := float64(L.CheckNumber(4))
		endAngle := float64(L.CheckNumber(5))

		if r <= 0 {
			L.RaiseError("graphics.draw_arc: radius must be > 0")
			return 0
		}

		optsArg := lua.LValue(lua.LNil)
		if L.GetTop() >= 6 {
			optsArg = L.Get(6)
		}

		opts := parseShapeOptions(L, "graphics.draw_arc", optsArg, true)
		if opts.filled {
			L.RaiseError("graphics.draw_arc: filled arcs are not supported")
			return 0
		}

		segments := parseSegmentsOption(L, "graphics.draw_arc", optsArg, defaultArcSegments, 1)
		points := buildArcPoints(x, y, r, startAngle, endAngle, segments)

		err := renderWithColor(renderer, opts.color, func() error {
			return drawPolygonOutline(renderer, points, false)
		})
		if err != nil {
			L.RaiseError("graphics.draw_arc failed: %v", err)
			return 0
		}

		return 0
	}))
}

func loadTexture(renderer *sdl.Renderer, fp string) (*sdl.Texture, error) {
	ext := strings.ToLower(filepath.Ext(fp))
	if ext == ".bmp" {
		surface, err := sdl.LoadBMP(fp)
		if err != nil {
			return nil, err
		}
		defer surface.Destroy()

		return renderer.CreateTextureFromSurface(surface)
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoded, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := decoded.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		return nil, os.ErrInvalid
	}

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, decoded, bounds.Min, draw.Src)

	surface, err := sdl.CreateSurfaceFrom(bounds.Dx(), bounds.Dy(), sdl.PIXELFORMAT_RGBA32, rgba.Pix, rgba.Stride)
	if err != nil {
		return nil, err
	}
	defer surface.Destroy()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	return texture, nil
}
