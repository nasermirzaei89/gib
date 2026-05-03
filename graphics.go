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
	released bool
}

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

// InitGraphics registers the graphics module and its functions to the Lua state.
// Graphics module is responsible rendering.
func InitGraphics(l *lua.LState, renderer *sdl.Renderer, baseDir string) {
	mt := l.NewTypeMetatable(imageMetatableName)
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
		ud.Value = &Image{texture: texture}
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
