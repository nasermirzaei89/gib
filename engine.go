package gib

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/Zyko0/go-sdl3/bin/binmix"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/mixer"
	"github.com/Zyko0/go-sdl3/sdl"
	lua "github.com/yuin/gopher-lua"
)

const (
	defaultTPS = 60.0
)

type startupWindowConfig struct {
	title      string
	width      int
	height     int
	resizable  bool
	fullscreen bool
}

type startupGraphicsConfig struct {
	autoClear  bool
	clearColor sdl.Color
}

type startupConfig struct {
	window   startupWindowConfig
	graphics startupGraphicsConfig
	tps      float64
}

func defaultStartupConfig() startupConfig {
	return startupConfig{
		window: startupWindowConfig{
			title:      "Game",
			width:      800,
			height:     600,
			resizable:  false,
			fullscreen: false,
		},
		graphics: startupGraphicsConfig{
			autoClear:  true,
			clearColor: sdl.Color{R: 0, G: 0, B: 0, A: 255},
		},
		tps: defaultTPS,
	}
}

func colorChannelToLuaNumber(channel uint8) lua.LNumber {
	return lua.LNumber(float64(channel) / 255.0)
}

func startupConfigToLuaTable(L *lua.LState, cfg startupConfig) *lua.LTable {
	window := L.NewTable()
	window.RawSetString("title", lua.LString(cfg.window.title))
	window.RawSetString("width", lua.LNumber(cfg.window.width))
	window.RawSetString("height", lua.LNumber(cfg.window.height))
	window.RawSetString("resizable", lua.LBool(cfg.window.resizable))
	window.RawSetString("fullscreen", lua.LBool(cfg.window.fullscreen))

	clearColor := L.NewTable()
	clearColor.Append(colorChannelToLuaNumber(cfg.graphics.clearColor.R))
	clearColor.Append(colorChannelToLuaNumber(cfg.graphics.clearColor.G))
	clearColor.Append(colorChannelToLuaNumber(cfg.graphics.clearColor.B))
	clearColor.Append(colorChannelToLuaNumber(cfg.graphics.clearColor.A))

	graphics := L.NewTable()
	graphics.RawSetString("auto_clear", lua.LBool(cfg.graphics.autoClear))
	graphics.RawSetString("clear_color", clearColor)

	conf := L.NewTable()
	conf.RawSetString("window", window)
	conf.RawSetString("graphics", graphics)
	conf.RawSetString("tps", lua.LNumber(cfg.tps))

	return conf
}

func parseStartupConfig(conf *lua.LTable) (startupConfig, error) {
	cfg := defaultStartupConfig()

	tpsVal := conf.RawGetString("tps")
	tpsNum, ok := tpsVal.(lua.LNumber)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.tps must be number")
	}
	cfg.tps = float64(tpsNum)
	if cfg.tps <= 0 {
		return cfg, fmt.Errorf("game.config(conf): conf.tps must be > 0")
	}

	windowVal := conf.RawGetString("window")
	windowTable, ok := windowVal.(*lua.LTable)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.window must be table")
	}

	titleVal := windowTable.RawGetString("title")
	title, ok := titleVal.(lua.LString)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.window.title must be string")
	}
	cfg.window.title = string(title)

	widthVal := windowTable.RawGetString("width")
	widthNum, ok := widthVal.(lua.LNumber)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.window.width must be number")
	}
	if lua.LNumber(int64(widthNum)) != widthNum {
		return cfg, fmt.Errorf("game.config(conf): conf.window.width must be integer")
	}
	cfg.window.width = int(int64(widthNum))
	if cfg.window.width <= 0 {
		return cfg, fmt.Errorf("game.config(conf): conf.window.width must be > 0")
	}

	heightVal := windowTable.RawGetString("height")
	heightNum, ok := heightVal.(lua.LNumber)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.window.height must be number")
	}
	if lua.LNumber(int64(heightNum)) != heightNum {
		return cfg, fmt.Errorf("game.config(conf): conf.window.height must be integer")
	}
	cfg.window.height = int(int64(heightNum))
	if cfg.window.height <= 0 {
		return cfg, fmt.Errorf("game.config(conf): conf.window.height must be > 0")
	}

	resizableVal := windowTable.RawGetString("resizable")
	resizable, ok := resizableVal.(lua.LBool)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.window.resizable must be boolean")
	}
	cfg.window.resizable = bool(resizable)

	fullscreenVal := windowTable.RawGetString("fullscreen")
	fullscreen, ok := fullscreenVal.(lua.LBool)
	if !ok {
		return cfg, fmt.Errorf("game.config(conf): conf.window.fullscreen must be boolean")
	}
	cfg.window.fullscreen = bool(fullscreen)

	graphicsVal := conf.RawGetString("graphics")
	if graphicsVal.Type() != lua.LTNil {
		graphicsTable, ok := graphicsVal.(*lua.LTable)
		if !ok {
			return cfg, fmt.Errorf("game.config(conf): conf.graphics must be table")
		}

		autoClearVal := graphicsTable.RawGetString("auto_clear")
		if autoClearVal.Type() != lua.LTNil {
			autoClear, ok := autoClearVal.(lua.LBool)
			if !ok {
				return cfg, fmt.Errorf("game.config(conf): conf.graphics.auto_clear must be boolean")
			}
			cfg.graphics.autoClear = bool(autoClear)
		}

		clearColorVal := graphicsTable.RawGetString("clear_color")
		if clearColorVal.Type() != lua.LTNil {
			clearColorTable, ok := clearColorVal.(*lua.LTable)
			if !ok {
				return cfg, fmt.Errorf("game.config(conf): conf.graphics.clear_color must be table {r, g, b, a}")
			}

			if clearColorTable.Len() != 4 {
				return cfg, fmt.Errorf("game.config(conf): conf.graphics.clear_color must contain exactly 4 values (r, g, b, a)")
			}

			components := [4]uint8{}
			for i := 1; i <= 4; i++ {
				componentVal := clearColorTable.RawGetInt(i)
				componentNum, ok := componentVal.(lua.LNumber)
				if !ok {
					return cfg, fmt.Errorf("game.config(conf): conf.graphics.clear_color value at index %d must be number", i)
				}

				component := float32(componentNum)
				if component < 0 || component > 1 {
					return cfg, fmt.Errorf("game.config(conf): conf.graphics.clear_color values must be in range [0, 1]")
				}

				components[i-1] = normalizeColorChannel(component)
			}

			cfg.graphics.clearColor = sdl.Color{
				R: components[0],
				G: components[1],
				B: components[2],
				A: components[3],
			}
		}
	}

	return cfg, nil
}

func resolveRunTarget(target string) (string, string, error) {
	if target == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("failed to get current directory: %w", err)
		}

		candidate := filepath.Join(cwd, "main.lua")
		info, err := os.Stat(candidate)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return "", "", fmt.Errorf("failed to stat %q: %w", candidate, err)
			}

			return "", cwd, nil
		}

		if info.IsDir() {
			return "", "", fmt.Errorf("expected %q to be a file", candidate)
		}

		return candidate, cwd, nil
	}

	absArg, err := filepath.Abs(target)
	if err != nil {
		return "", "", fmt.Errorf("failed to get absolute path for %q: %w", target, err)
	}

	info, err := os.Stat(absArg)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve path %q: %w", target, err)
	}

	if info.IsDir() {
		candidate := filepath.Join(absArg, "main.lua")
		candidateInfo, err := os.Stat(candidate)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return "", "", fmt.Errorf("directory %q does not contain main.lua (%s)", target, candidate)
			}
			return "", "", fmt.Errorf("failed to stat %q: %w", candidate, err)
		}

		if candidateInfo.IsDir() {
			return "", "", fmt.Errorf("expected %q to be a file", candidate)
		}

		return candidate, absArg, nil
	}

	if !strings.EqualFold(filepath.Ext(absArg), ".lua") {
		return "", "", fmt.Errorf("only .lua files are supported, got %q", target)
	}

	return absArg, filepath.Dir(absArg), nil
}

func RunGame(target string) error {
	scriptPath, baseDir, err := resolveRunTarget(target)
	if err != nil {
		return err
	}

	L := lua.NewState()
	defer L.Close()

	gameTable := L.NewTable()
	L.SetGlobal("game", gameTable)
	conf := startupConfigToLuaTable(L, defaultStartupConfig())

	if scriptPath != "" {
		err := L.DoFile(scriptPath)
		if err != nil {
			return fmt.Errorf("failed to load Lua file %q: %w", scriptPath, err)
		}
	}

	globalGame := L.GetGlobal("game")
	tbl, ok := globalGame.(*lua.LTable)
	if !ok {
		return fmt.Errorf("global game must be a table")
	}
	gameTable = tbl

	configFn := gameTable.RawGetString("config")
	if configFn.Type() != lua.LTNil && configFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.config(conf) found")
	}
	if configFn.Type() == lua.LTFunction {
		err := L.CallByParam(lua.P{
			Fn:      configFn,
			NRet:    0,
			Protect: true,
		}, conf)
		if err != nil {
			return fmt.Errorf("failed to call game.config(conf): %w", err)
		}
	}

	startupCfg, err := parseStartupConfig(conf)
	if err != nil {
		return err
	}

	defer binsdl.Load().Unload() // sdl.LoadLibrary(sdl.Path())
	defer binmix.Load().Unload()
	defer sdl.Quit()

	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		return fmt.Errorf("failed to initialize SDL: %w", err)
	}

	err = mixer.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize SDL_mixer: %w", err)
	}
	defer mixer.Quit()

	audioMixer, err := mixer.CreateMixerDevice(sdl.AUDIO_DEVICE_DEFAULT_PLAYBACK, nil)
	if err != nil {
		return fmt.Errorf("failed to create audio mixer: %w", err)
	}
	defer audioMixer.Destroy()

	var windowFlags sdl.WindowFlags
	if startupCfg.window.resizable {
		windowFlags |= sdl.WINDOW_RESIZABLE
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(
		startupCfg.window.title,
		startupCfg.window.width,
		startupCfg.window.height,
		windowFlags,
	)
	if err != nil {
		return fmt.Errorf("failed to create window and renderer: %w", err)
	}

	if startupCfg.window.fullscreen {
		err = window.SetFullscreen(true)
		if err != nil {
			return fmt.Errorf("failed to set fullscreen on startup: %w", err)
		}
	}

	err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return fmt.Errorf("failed to set renderer draw blend mode: %w", err)
	}

	loadFn := gameTable.RawGetString("load")
	if loadFn.Type() != lua.LTNil && loadFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.load() found")
	}

	fixedUpdateFn := gameTable.RawGetString("fixed_update")
	if fixedUpdateFn.Type() != lua.LTNil && fixedUpdateFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.fixed_update(dt) found")
	}

	updateFn := gameTable.RawGetString("update")
	if updateFn.Type() != lua.LTNil && updateFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.update(dt) found")
	}

	renderFn := gameTable.RawGetString("render")
	if renderFn.Type() != lua.LTNil && renderFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.render() found")
	}

	eventFn := gameTable.RawGetString("event")
	if eventFn.Type() != lua.LTNil && eventFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.event(e) found")
	}

	InitDebug(L, renderer)
	InitLog(L)
	inputState := newInputState()
	InitInput(L, inputState)
	graphicsState := &graphicsRuntimeState{
		autoClear:  startupCfg.graphics.autoClear,
		clearColor: startupCfg.graphics.clearColor,
	}
	InitGraphics(L, renderer, baseDir, graphicsState)
	audioState := InitAudio(L, audioMixer, baseDir)
	defer audioState.Close()
	quitRequested := false
	InitWindow(L, window, &quitRequested)

	defer window.Destroy()
	defer renderer.Destroy()

	renderer.SetDrawColor(255, 255, 255, 255)

	fixedDelta := 1.0 / startupCfg.tps
	last := sdl.TicksNS()
	var accumulator float64

	if loadFn.Type() != lua.LTNil {
		err := L.CallByParam(lua.P{
			Fn:      loadFn,
			NRet:    0,
			Protect: true,
		})
		if err != nil {
			return fmt.Errorf("failed to call game.load(): %w", err)
		}
	}

	err = sdl.RunLoop(func() error {
		if quitRequested {
			return sdl.EndLoop
		}

		// Keep edge states alive for the full frame so both fixed_update and update can read them.
		inputState.beginFrame()

		var event sdl.Event

		for sdl.PollEvent(&event) {
			if quitRequested {
				return sdl.EndLoop
			}

			switch event.Type {
			case sdl.EVENT_KEY_DOWN:
				if keyboardEvent := event.KeyboardEvent(); keyboardEvent != nil {
					inputState.setKeyDown(keyboardEvent.Scancode)
				}
			case sdl.EVENT_KEY_UP:
				if keyboardEvent := event.KeyboardEvent(); keyboardEvent != nil {
					inputState.setKeyUp(keyboardEvent.Scancode)
				}
			case sdl.EVENT_MOUSE_MOTION:
				if mouseMotionEvent := event.MouseMotionEvent(); mouseMotionEvent != nil {
					inputState.setMousePosition(mouseMotionEvent.X, mouseMotionEvent.Y)
				}
			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				if mouseButtonEvent := event.MouseButtonEvent(); mouseButtonEvent != nil {
					inputState.setMousePosition(mouseButtonEvent.X, mouseButtonEvent.Y)
					inputState.setMouseButtonDown(mouseButtonEvent.Button)
				}
			case sdl.EVENT_MOUSE_BUTTON_UP:
				if mouseButtonEvent := event.MouseButtonEvent(); mouseButtonEvent != nil {
					inputState.setMousePosition(mouseButtonEvent.X, mouseButtonEvent.Y)
					inputState.setMouseButtonUp(mouseButtonEvent.Button)
				}
			case sdl.EVENT_WINDOW_FOCUS_LOST:
				inputState.clearAll()
			}

			if eventFn.Type() != lua.LTNil {
				eventTable := L.NewTable()
				eventTable.RawSetString("type", lua.LString(eventTypeString(event.Type)))

				switch event.Type {
				case sdl.EVENT_KEY_DOWN:
					if keyboardEvent := event.KeyboardEvent(); keyboardEvent != nil {
						eventTable.RawSetString("key", lua.LString(canonicalKeyNameFromScancode(keyboardEvent.Scancode)))
						eventTable.RawSetString("is_repeat", lua.LBool(keyboardEvent.Repeat))
					}
				case sdl.EVENT_KEY_UP:
					if keyboardEvent := event.KeyboardEvent(); keyboardEvent != nil {
						eventTable.RawSetString("key", lua.LString(canonicalKeyNameFromScancode(keyboardEvent.Scancode)))
					}
				case sdl.EVENT_MOUSE_MOTION:
					if mouseMotionEvent := event.MouseMotionEvent(); mouseMotionEvent != nil {
						eventTable.RawSetString("x", lua.LNumber(mouseMotionEvent.X))
						eventTable.RawSetString("y", lua.LNumber(mouseMotionEvent.Y))
						eventTable.RawSetString("dx", lua.LNumber(mouseMotionEvent.Xrel))
						eventTable.RawSetString("dy", lua.LNumber(mouseMotionEvent.Yrel))
					}
				case sdl.EVENT_MOUSE_BUTTON_DOWN, sdl.EVENT_MOUSE_BUTTON_UP:
					if mouseButtonEvent := event.MouseButtonEvent(); mouseButtonEvent != nil {
						eventTable.RawSetString("button", lua.LString(canonicalMouseButtonName(mouseButtonEvent.Button)))
						eventTable.RawSetString("clicks", lua.LNumber(mouseButtonEvent.Clicks))
						eventTable.RawSetString("x", lua.LNumber(mouseButtonEvent.X))
						eventTable.RawSetString("y", lua.LNumber(mouseButtonEvent.Y))
					}
				case sdl.EVENT_MOUSE_WHEEL:
					if mouseWheelEvent := event.MouseWheelEvent(); mouseWheelEvent != nil {
						eventTable.RawSetString("x", lua.LNumber(mouseWheelEvent.X))
						eventTable.RawSetString("y", lua.LNumber(mouseWheelEvent.Y))
						eventTable.RawSetString("mouse_x", lua.LNumber(mouseWheelEvent.MouseX))
						eventTable.RawSetString("mouse_y", lua.LNumber(mouseWheelEvent.MouseY))
					}
				}

				err := L.CallByParam(lua.P{
					Fn:      eventFn,
					NRet:    0,
					Protect: true,
				}, eventTable)
				if err != nil {
					slog.Error("Failed to call game.event(event)", "error", err)
					return err
				}
			}

			switch event.Type {
			case sdl.EVENT_QUIT:
				return sdl.EndLoop
			case sdl.EVENT_KEY_DOWN:
				fallthrough
			case sdl.EVENT_KEY_UP:
				fallthrough
			case sdl.EVENT_MOUSE_MOTION:
				fallthrough
			case sdl.EVENT_MOUSE_BUTTON_DOWN:
				fallthrough
			case sdl.EVENT_MOUSE_BUTTON_UP:
				fallthrough
			case sdl.EVENT_WINDOW_FOCUS_LOST:
				// Handled above for input-state updates.
			default:
				slog.Debug("Unhandled event", "type", event.Type)
			}
		}

		current := sdl.TicksNS()
		dt := float64(current-last) / 1e9
		last = current
		accumulator += dt

		if fixedUpdateFn.Type() != lua.LTNil {
			for accumulator >= fixedDelta {
				err := L.CallByParam(lua.P{
					Fn:      fixedUpdateFn,
					NRet:    0,
					Protect: true,
				}, lua.LNumber(fixedDelta))
				if err != nil {
					slog.Error("Failed to call game.fixed_update(fixed_dt)", "error", err)
					return err
				}
				accumulator -= fixedDelta
			}
		}

		if updateFn.Type() != lua.LTNil {
			err := L.CallByParam(lua.P{
				Fn:      updateFn,
				NRet:    0,
				Protect: true,
			}, lua.LNumber(dt))
			if err != nil {
				slog.Error("Failed to call game.update(dt)", "error", err)
				return err
			}
		}

		if graphicsState.autoClear {
			err := clearRenderer(renderer, graphicsState.clearColor)
			if err != nil {
				slog.Error("Failed to clear renderer", "error", err)
				return err
			}
		}

		if renderFn.Type() != lua.LTNil {
			err := L.CallByParam(lua.P{
				Fn:      renderFn,
				NRet:    0,
				Protect: true,
			})
			if err != nil {
				slog.Error("Failed to call game.render()", "error", err)
				return err
			}
		}

		renderer.Present()

		return nil
	})
	if err != nil {
		return fmt.Errorf("error in main loop: %w", err)
	}

	return nil
}
