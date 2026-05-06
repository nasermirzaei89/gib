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
	tps        = 60.0
	fixedDelta = 1.0 / tps
)

func Run() error {
	L := lua.NewState()
	defer L.Close()

	gameTable := L.NewTable()
	L.SetGlobal("game", gameTable)

	var (
		scriptPath   string
		baseDir      string
		scriptLoaded bool
	)

	if len(os.Args) <= 1 {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		baseDir = cwd
		candidate := filepath.Join(cwd, "main.lua")
		info, err := os.Stat(candidate)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to stat %q: %w", candidate, err)
			}
		} else {
			if info.IsDir() {
				return fmt.Errorf("expected %q to be a file", candidate)
			}

			scriptPath = candidate
		}
	} else {
		arg := os.Args[1]
		absArg, err := filepath.Abs(arg)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %q: %w", arg, err)
		}

		info, err := os.Stat(absArg)
		if err != nil {
			return fmt.Errorf("failed to resolve path %q: %w", arg, err)
		}

		if info.IsDir() {
			baseDir = absArg
			candidate := filepath.Join(absArg, "main.lua")
			candidateInfo, err := os.Stat(candidate)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("directory %q does not contain main.lua (%s)", arg, candidate)
				}
				return fmt.Errorf("failed to stat %q: %w", candidate, err)
			}

			if candidateInfo.IsDir() {
				return fmt.Errorf("expected %q to be a file", candidate)
			}

			scriptPath = candidate
		} else {
			if !strings.EqualFold(filepath.Ext(absArg), ".lua") {
				return fmt.Errorf("only .lua files are supported, got %q", arg)
			}

			baseDir = filepath.Dir(absArg)
			scriptPath = absArg
		}
	}

	if scriptPath != "" {
		err := L.DoFile(scriptPath)
		if err != nil {
			return fmt.Errorf("failed to load Lua file %q: %w", scriptPath, err)
		}
		scriptLoaded = true
	}

	defer binsdl.Load().Unload() // sdl.LoadLibrary(sdl.Path())
	defer binmix.Load().Unload()
	defer sdl.Quit()

	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
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

	window, renderer, err := sdl.CreateWindowAndRenderer("Game", 800, 600, 0)
	if err != nil {
		return fmt.Errorf("failed to create window and renderer: %w", err)
	}

	err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		return fmt.Errorf("failed to set renderer draw blend mode: %w", err)
	}

	var game *lua.LTable
	if scriptLoaded {
		ret := L.Get(-1)
		if tbl, ok := ret.(*lua.LTable); ok {
			game = tbl
		}
	}

	if game == nil {
		g := L.GetGlobal("game")
		if tbl, ok := g.(*lua.LTable); ok {
			game = tbl
		} else {
			return fmt.Errorf("no game table found")
		}
	}

	loadFn := game.RawGetString("load")
	if loadFn.Type() != lua.LTNil && loadFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.load() found")
	}

	fixedUpdateFn := game.RawGetString("fixed_update")
	if fixedUpdateFn.Type() != lua.LTNil && fixedUpdateFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.fixed_update(dt) found")
	}

	updateFn := game.RawGetString("update")
	if updateFn.Type() != lua.LTNil && updateFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.update(dt) found")
	}

	renderFn := game.RawGetString("render")
	if renderFn.Type() != lua.LTNil && renderFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.render() found")
	}

	eventFn := game.RawGetString("event")
	if eventFn.Type() != lua.LTNil && eventFn.Type() != lua.LTFunction {
		return fmt.Errorf("no game.event() found")
	}

	InitDebug(L, renderer)
	InitLog(L)
	inputState := newInputState()
	InitInput(L, inputState)
	InitGraphics(L, renderer, baseDir)
	audioState := InitAudio(L, audioMixer, baseDir)
	defer audioState.Close()
	quitRequested := false
	InitWindow(L, window, &quitRequested)

	defer window.Destroy()
	defer renderer.Destroy()

	renderer.SetDrawColor(255, 255, 255, 255)

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
