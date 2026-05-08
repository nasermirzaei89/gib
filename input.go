package gib

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

	mouseDown     map[uint8]bool
	mousePressed  map[uint8]bool
	mouseReleased map[uint8]bool
	mouseAliases  map[string]uint8
	mouseX        float32
	mouseY        float32
}

func newInputState() *inputState {
	_, mouseX, mouseY := sdl.GetMouseState()

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
		mouseDown:     make(map[uint8]bool),
		mousePressed:  make(map[uint8]bool),
		mouseReleased: make(map[uint8]bool),
		mouseAliases: map[string]uint8{
			"left":   uint8(sdl.BUTTON_LEFT),
			"middle": uint8(sdl.BUTTON_MIDDLE),
			"right":  uint8(sdl.BUTTON_RIGHT),
			"x1":     uint8(sdl.BUTTON_X1),
			"x2":     uint8(sdl.BUTTON_X2),
		},
		mouseX: mouseX,
		mouseY: mouseY,
	}
}

func normalizeKeyName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func canonicalKeyNameFromScancode(scancode sdl.Scancode) string {
	switch scancode {
	case sdl.SCANCODE_RETURN:
		return "enter"
	case sdl.SCANCODE_ESCAPE:
		return "escape"
	case sdl.SCANCODE_LSHIFT:
		return "left_shift"
	case sdl.SCANCODE_RSHIFT:
		return "right_shift"
	case sdl.SCANCODE_LCTRL:
		return "left_ctrl"
	case sdl.SCANCODE_RCTRL:
		return "right_ctrl"
	case sdl.SCANCODE_LALT:
		return "left_alt"
	case sdl.SCANCODE_RALT:
		return "right_alt"
	}

	name := normalizeKeyName(scancode.Name())
	name = strings.ReplaceAll(name, " ", "_")
	return name
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
	clear(s.mousePressed)
	clear(s.mouseReleased)
}

func (s *inputState) clearAll() {
	clear(s.down)
	clear(s.pressed)
	clear(s.released)
	clear(s.mouseDown)
	clear(s.mousePressed)
	clear(s.mouseReleased)
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

func (s *inputState) resolveMouseButton(name string) (uint8, error) {
	normalized := normalizeKeyName(name)
	if normalized == "" {
		return 0, fmt.Errorf("mouse button name cannot be empty")
	}

	button, ok := s.mouseAliases[normalized]
	if !ok {
		return 0, fmt.Errorf("unknown mouse button %q", name)
	}

	return button, nil
}

func (s *inputState) setMousePosition(x, y float32) {
	s.mouseX = x
	s.mouseY = y
}

func (s *inputState) setMouseButtonDown(button uint8) {
	if button == 0 {
		return
	}

	wasDown := s.mouseDown[button]
	s.mouseDown[button] = true
	if !wasDown {
		s.mousePressed[button] = true
	}
}

func (s *inputState) setMouseButtonUp(button uint8) {
	if button == 0 {
		return
	}

	wasDown := s.mouseDown[button]
	s.mouseDown[button] = false
	if wasDown {
		s.mouseReleased[button] = true
	}
}

func (s *inputState) getMousePosition() (int, int) {
	return int(s.mouseX), int(s.mouseY)
}

func (s *inputState) isMouseButtonDown(name string) (bool, error) {
	button, err := s.resolveMouseButton(name)
	if err != nil {
		return false, err
	}
	return s.mouseDown[button], nil
}

func (s *inputState) isMouseButtonPressed(name string) (bool, error) {
	button, err := s.resolveMouseButton(name)
	if err != nil {
		return false, err
	}
	return s.mousePressed[button], nil
}

func (s *inputState) isMouseButtonReleased(name string) (bool, error) {
	button, err := s.resolveMouseButton(name)
	if err != nil {
		return false, err
	}
	return s.mouseReleased[button], nil
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

	inputTable.RawSetString("get_mouse_position", l.NewFunction(func(l *lua.LState) int {
		x, y := state.getMousePosition()
		l.Push(lua.LNumber(x))
		l.Push(lua.LNumber(y))
		return 2
	}))

	inputTable.RawSetString("is_mouse_button_down", l.NewFunction(func(l *lua.LState) int {
		buttonName := l.CheckString(1)

		isDown, err := state.isMouseButtonDown(buttonName)
		if err != nil {
			l.RaiseError("%s", err)
			return 0
		}

		l.Push(lua.LBool(isDown))
		return 1
	}))

	inputTable.RawSetString("is_mouse_button_pressed", l.NewFunction(func(l *lua.LState) int {
		buttonName := l.CheckString(1)

		isPressed, err := state.isMouseButtonPressed(buttonName)
		if err != nil {
			l.RaiseError("%s", err)
			return 0
		}

		l.Push(lua.LBool(isPressed))
		return 1
	}))

	inputTable.RawSetString("is_mouse_button_released", l.NewFunction(func(l *lua.LState) int {
		buttonName := l.CheckString(1)

		isReleased, err := state.isMouseButtonReleased(buttonName)
		if err != nil {
			l.RaiseError("%s", err)
			return 0
		}

		l.Push(lua.LBool(isReleased))
		return 1
	}))
}
