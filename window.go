package gameengine

import (
	"github.com/Zyko0/go-sdl3/sdl"
	lua "github.com/yuin/gopher-lua"
)

func checkInt32Arg(L *lua.LState, idx int, name string) int32 {
	n := L.CheckNumber(idx)
	if lua.LNumber(int64(n)) != n {
		L.RaiseError("window.%s must be an integer", name)
		return 0
	}

	return int32(int64(n))
}

func InitWindow(l *lua.LState, window *sdl.Window, requestQuit *bool) {
	windowTable := l.NewTable()
	l.SetGlobal("window", windowTable)

	windowTable.RawSetString("get_size", l.NewFunction(func(L *lua.LState) int {
		w, h, err := window.Size()
		if err != nil {
			L.RaiseError("window.get_size failed: %v", err)
			return 0
		}

		L.Push(lua.LNumber(w))
		L.Push(lua.LNumber(h))

		return 2
	}))

	windowTable.RawSetString("set_size", l.NewFunction(func(L *lua.LState) int {
		w := checkInt32Arg(L, 1, "set_size(w, h): w")
		h := checkInt32Arg(L, 2, "set_size(w, h): h")

		if w <= 0 || h <= 0 {
			L.RaiseError("window.set_size(w, h) expects w > 0 and h > 0")
			return 0
		}

		err := window.SetSize(w, h)
		if err != nil {
			L.RaiseError("window.set_size failed: %v", err)
			return 0
		}

		return 0
	}))

	windowTable.RawSetString("set_title", l.NewFunction(func(L *lua.LState) int {
		title := L.CheckString(1)

		err := window.SetTitle(title)
		if err != nil {
			L.RaiseError("window.set_title failed: %v", err)
			return 0
		}

		return 0
	}))

	windowTable.RawSetString("set_fullscreen", l.NewFunction(func(L *lua.LState) int {
		fullscreen := L.CheckBool(1)

		err := window.SetFullscreen(fullscreen)
		if err != nil {
			L.RaiseError("window.set_fullscreen failed: %v", err)
			return 0
		}

		return 0
	}))

	windowTable.RawSetString("is_fullscreen", l.NewFunction(func(L *lua.LState) int {
		isFullscreen := (window.Flags() & sdl.WINDOW_FULLSCREEN) != 0
		L.Push(lua.LBool(isFullscreen))

		return 1
	}))

	windowTable.RawSetString("close", l.NewFunction(func(L *lua.LState) int {
		if requestQuit != nil {
			*requestQuit = true
		}

		return 0
	}))
}
