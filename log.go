package gameengine

import (
	"log/slog"

	lua "github.com/yuin/gopher-lua"
)

func InitLog(l *lua.LState) {
	log := l.NewTable()
	l.SetGlobal("log", log)

	log.RawSetString("debug", l.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		slog.Debug(msg)
		return 0
	}))

	log.RawSetString("info", l.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		slog.Info(msg)
		return 0
	}))

	log.RawSetString("warn", l.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		slog.Warn(msg)
		return 0
	}))

	log.RawSetString("error", l.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		slog.Error(msg)
		return 0
	}))
}
