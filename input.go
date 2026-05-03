package gameengine

import (
	"fmt"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
	lua "github.com/yuin/gopher-lua"
)

type inputState struct {
	down     map[sdl.Scancode]bool
	pressed  map[sdl.Scancode]bool
	released map[sdl.Scancode]bool
	aliases  map[string]sdl.Scancode
}

func newInputState() *inputState {
	return &inputState{
		down:     make(map[sdl.Scancode]bool),
		pressed:  make(map[sdl.Scancode]bool),
		released: make(map[sdl.Scancode]bool),
		aliases: map[string]sdl.Scancode{
			"left":        sdl.SCANCODE_LEFT,
			"right":       sdl.SCANCODE_RIGHT,
			"up":          sdl.SCANCODE_UP,
			"down":        sdl.SCANCODE_DOWN,
			"space":       sdl.SCANCODE_SPACE,
			"enter":       sdl.SCANCODE_RETURN,
			"return":      sdl.SCANCODE_RETURN,
			"esc":         sdl.SCANCODE_ESCAPE,
			"escape":      sdl.SCANCODE_ESCAPE,
			"tab":         sdl.SCANCODE_TAB,
			"backspace":   sdl.SCANCODE_BACKSPACE,
			"shift":       sdl.SCANCODE_LSHIFT,
			"left_shift":  sdl.SCANCODE_LSHIFT,
			"right_shift": sdl.SCANCODE_RSHIFT,
			"lshift":      sdl.SCANCODE_LSHIFT,
			"rshift":      sdl.SCANCODE_RSHIFT,
			"ctrl":        sdl.SCANCODE_LCTRL,
			"control":     sdl.SCANCODE_LCTRL,
			"left_ctrl":   sdl.SCANCODE_LCTRL,
			"right_ctrl":  sdl.SCANCODE_RCTRL,
			"lctrl":       sdl.SCANCODE_LCTRL,
			"rctrl":       sdl.SCANCODE_RCTRL,
			"alt":         sdl.SCANCODE_LALT,
			"left_alt":    sdl.SCANCODE_LALT,
			"right_alt":   sdl.SCANCODE_RALT,
			"lalt":        sdl.SCANCODE_LALT,
			"ralt":        sdl.SCANCODE_RALT,
		},
	}
}

func normalizeKeyName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (s *inputState) resolveScancode(name string) (sdl.Scancode, error) {
	normalized := normalizeKeyName(name)
	if normalized == "" {
		return sdl.SCANCODE_UNKNOWN, fmt.Errorf("key name cannot be empty")
	}

	if scancode, ok := s.aliases[normalized]; ok {
		return scancode, nil
	}

	candidates := []string{
		strings.TrimSpace(name),
		normalized,
		strings.ReplaceAll(normalized, "_", " "),
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}

		scancode := sdl.GetScancodeFromName(candidate)
		if scancode != sdl.SCANCODE_UNKNOWN {
			return scancode, nil
		}
	}

	return sdl.SCANCODE_UNKNOWN, fmt.Errorf("unknown key name %q", name)
}

func (s *inputState) beginFrame() {
	clear(s.pressed)
	clear(s.released)
}

func (s *inputState) clearAll() {
	clear(s.down)
	clear(s.pressed)
	clear(s.released)
}

func (s *inputState) setKeyDown(scancode sdl.Scancode) {
	if scancode == sdl.SCANCODE_UNKNOWN {
		return
	}

	wasDown := s.down[scancode]
	s.down[scancode] = true
	if !wasDown {
		s.pressed[scancode] = true
	}
}

func (s *inputState) setKeyUp(scancode sdl.Scancode) {
	if scancode == sdl.SCANCODE_UNKNOWN {
		return
	}

	wasDown := s.down[scancode]
	s.down[scancode] = false
	if wasDown {
		s.released[scancode] = true
	}
}

func (s *inputState) isKeyDown(name string) (bool, error) {
	scancode, err := s.resolveScancode(name)
	if err != nil {
		return false, err
	}
	return s.down[scancode], nil
}

func (s *inputState) isKeyPressed(name string) (bool, error) {
	scancode, err := s.resolveScancode(name)
	if err != nil {
		return false, err
	}
	return s.pressed[scancode], nil
}

func (s *inputState) isKeyReleased(name string) (bool, error) {
	scancode, err := s.resolveScancode(name)
	if err != nil {
		return false, err
	}
	return s.released[scancode], nil
}

func InitInput(l *lua.LState, state *inputState) {
	inputTable := l.NewTable()
	l.SetGlobal("input", inputTable)

	inputTable.RawSetString("is_key_down", l.NewFunction(func(l *lua.LState) int {
		keyName := l.CheckString(1)

		isDown, err := state.isKeyDown(keyName)
		if err != nil {
			l.RaiseError("%s", err)
			return 0
		}

		l.Push(lua.LBool(isDown))
		return 1
	}))

	inputTable.RawSetString("is_key_pressed", l.NewFunction(func(l *lua.LState) int {
		keyName := l.CheckString(1)

		isPressed, err := state.isKeyPressed(keyName)
		if err != nil {
			l.RaiseError("%s", err)
			return 0
		}

		l.Push(lua.LBool(isPressed))
		return 1
	}))

	inputTable.RawSetString("is_key_released", l.NewFunction(func(l *lua.LState) int {
		keyName := l.CheckString(1)

		isReleased, err := state.isKeyReleased(keyName)
		if err != nil {
			l.RaiseError("%s", err)
			return 0
		}

		l.Push(lua.LBool(isReleased))
		return 1
	}))
}
