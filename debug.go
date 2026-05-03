package gameengine

import (
	"github.com/Zyko0/go-sdl3/sdl"
	lua "github.com/yuin/gopher-lua"
)

func InitDebug(l *lua.LState, renderer *sdl.Renderer) {
	debug := l.NewTable()
	l.SetGlobal("debug", debug)

	debug.RawSetString("print", l.NewFunction(func(L *lua.LState) int {
		x := L.ToNumber(1)
		y := L.ToNumber(2)
		msg := L.ToString(3)

		renderer.DebugText(float32(x), float32(y), msg)
		return 0
	}))
}
