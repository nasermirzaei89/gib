package gib

import (
	"fmt"
	"path/filepath"

	"github.com/Zyko0/go-sdl3/mixer"
	lua "github.com/yuin/gopher-lua"
)

const (
	audioAssetMetatableName    = "audio.asset"
	audioInstanceMetatableName = "audio.instance"
)

type audioAsset struct {
	audio         *mixer.Audio
	instances     map[*audioInstance]struct{}
	defaultVolume float32
	released      bool
}

type audioInstance struct {
	asset    *audioAsset
	track    *mixer.Track
	volume   float32
	paused   bool
	stopped  bool
	released bool
}

type audioState struct {
	mixer        *mixer.Mixer
	baseDir      string
	masterVolume float32
	assets       map[*audioAsset]struct{}
	instances    map[*audioInstance]struct{}
}

func newAudioState(mix *mixer.Mixer, baseDir string) *audioState {
	return &audioState{
		mixer:        mix,
		baseDir:      baseDir,
		masterVolume: 1,
		assets:       make(map[*audioAsset]struct{}),
		instances:    make(map[*audioInstance]struct{}),
	}
}

func (s *audioState) Close() {
	if s == nil {
		return
	}

	for inst := range s.instances {
		s.releaseInstance(inst, true)
	}

	for asset := range s.assets {
		s.releaseAsset(asset)
	}
}

func (s *audioState) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(s.baseDir, path)
}

func (s *audioState) isInstancePlaying(inst *audioInstance) bool {
	if inst == nil || inst.released || inst.stopped || inst.track == nil {
		return false
	}

	return inst.track.Playing()
}

func (s *audioState) refreshInstanceState(inst *audioInstance) {
	if inst == nil || inst.released || inst.stopped || inst.track == nil {
		return
	}

	if !inst.track.Playing() && !inst.track.Paused() {
		inst.paused = false
		inst.stopped = true
		inst.track.Destroy()
		inst.track = nil
	}
}

func (s *audioState) applyInstanceGain(inst *audioInstance) error {
	if inst == nil || inst.released || inst.stopped || inst.track == nil {
		return nil
	}

	gain := s.masterVolume * inst.volume
	if err := inst.track.SetGain(gain); err != nil {
		return err
	}

	return nil
}

func (s *audioState) releaseInstance(inst *audioInstance, invalidate bool) {
	if inst == nil {
		return
	}

	if inst.track != nil {
		_ = inst.track.Stop(0)
		inst.track.Destroy()
		inst.track = nil
	}

	inst.paused = false
	inst.stopped = true

	if invalidate {
		inst.released = true
	}

	if inst.asset != nil {
		delete(inst.asset.instances, inst)
	}

	if invalidate {
		delete(s.instances, inst)
	}
}

func (s *audioState) releaseAsset(asset *audioAsset) {
	if asset == nil || asset.released {
		return
	}

	for inst := range asset.instances {
		s.releaseInstance(inst, true)
	}

	if asset.audio != nil {
		asset.audio.Destroy()
		asset.audio = nil
	}

	asset.released = true
	delete(s.assets, asset)
}

func checkAudioAsset(L *lua.LState, index int) *audioAsset {
	ud := L.CheckUserData(index)
	asset, ok := ud.Value.(*audioAsset)
	if !ok || asset == nil {
		L.ArgError(index, "expected audio asset userdata")
		return nil
	}

	return asset
}

func checkAudioInstance(L *lua.LState, index int) *audioInstance {
	ud := L.CheckUserData(index)
	inst, ok := ud.Value.(*audioInstance)
	if !ok || inst == nil {
		L.ArgError(index, "expected audio playback userdata")
		return nil
	}

	return inst
}

func checkAudioVolume(L *lua.LState, index int, funcName string) float32 {
	vol := float32(L.CheckNumber(index))
	if vol < 0 || vol > 1 {
		L.RaiseError("%s: volume must be in range [0, 1]", funcName)
		return 0
	}

	return vol
}

func parsePlayOptions(L *lua.LState, index int) (bool, float32) {
	loop := false
	volume := float32(1)

	if L.GetTop() < index {
		return loop, volume
	}

	arg := L.Get(index)
	if arg.Type() == lua.LTNil {
		return loop, volume
	}

	tbl, ok := arg.(*lua.LTable)
	if !ok {
		L.RaiseError("audio.play: options must be a table")
		return loop, volume
	}

	tbl.ForEach(func(key lua.LValue, _ lua.LValue) {
		keyString, ok := key.(lua.LString)
		if !ok {
			L.RaiseError("audio.play: option keys must be strings")
			return
		}

		name := string(keyString)
		if name != "loop" && name != "volume" {
			L.RaiseError("audio.play: unknown option key %q", name)
		}
	})

	loopVal := tbl.RawGetString("loop")
	if loopVal.Type() != lua.LTNil {
		loopBool, ok := loopVal.(lua.LBool)
		if !ok {
			L.RaiseError("audio.play: options key %q must be a boolean", "loop")
			return loop, volume
		}
		loop = bool(loopBool)
	}

	volVal := tbl.RawGetString("volume")
	if volVal.Type() != lua.LTNil {
		volNum, ok := volVal.(lua.LNumber)
		if !ok {
			L.RaiseError("audio.play: options key %q must be a number", "volume")
			return loop, volume
		}

		volume = float32(volNum)
		if volume < 0 || volume > 1 {
			L.RaiseError("audio.play: options key %q must be in range [0, 1]", "volume")
			return loop, volume
		}
	}

	return loop, volume
}

func checkAudioTarget(L *lua.LState, index int) (*audioAsset, *audioInstance) {
	ud := L.CheckUserData(index)
	switch target := ud.Value.(type) {
	case *audioAsset:
		if target == nil {
			L.ArgError(index, "expected audio target userdata")
			return nil, nil
		}
		return target, nil
	case *audioInstance:
		if target == nil {
			L.ArgError(index, "expected audio target userdata")
			return nil, nil
		}
		return nil, target
	default:
		L.ArgError(index, "expected audio target userdata")
		return nil, nil
	}
}

func ensureAssetAlive(L *lua.LState, asset *audioAsset, funcName string) {
	if asset.released || asset.audio == nil {
		L.RaiseError("%s: audio asset has been unloaded", funcName)
	}
}

func ensureInstanceAlive(L *lua.LState, inst *audioInstance, funcName string) {
	if inst.released {
		L.RaiseError("%s: audio playback has been invalidated", funcName)
	}
}

func controlInstance(state *audioState, inst *audioInstance, op string) error {
	state.refreshInstanceState(inst)

	switch op {
	case "stop":
		if inst.stopped {
			return nil
		}
		if inst.track != nil {
			if err := inst.track.Stop(0); err != nil {
				return err
			}
			inst.track.Destroy()
			inst.track = nil
		}
		inst.paused = false
		inst.stopped = true
		return nil
	case "pause":
		if inst.stopped || inst.paused {
			return nil
		}
		if inst.track != nil && inst.track.Playing() {
			if err := inst.track.Pause(); err != nil {
				return err
			}
		}
		inst.paused = true
		return nil
	case "resume":
		if inst.stopped || !inst.paused {
			return nil
		}
		if inst.track != nil {
			if err := inst.track.Resume(); err != nil {
				return err
			}
		}
		inst.paused = false
		return nil
	default:
		return fmt.Errorf("unsupported audio control operation %q", op)
	}
}

func controlAsset(state *audioState, asset *audioAsset, op string) error {
	for inst := range asset.instances {
		if inst.released {
			continue
		}

		if err := controlInstance(state, inst, op); err != nil {
			return err
		}
	}

	return nil
}

func isTargetPlaying(state *audioState, asset *audioAsset, inst *audioInstance) bool {
	if asset != nil {
		for candidate := range asset.instances {
			if candidate.released {
				continue
			}
			state.refreshInstanceState(candidate)
			if state.isInstancePlaying(candidate) {
				return true
			}
		}
		return false
	}

	state.refreshInstanceState(inst)
	return state.isInstancePlaying(inst)
}

// InitAudio registers the audio module and its functions to the Lua state.
func InitAudio(l *lua.LState, mix *mixer.Mixer, baseDir string) *audioState {
	state := newAudioState(mix, baseDir)

	assetMetatable := l.NewTypeMetatable(audioAssetMetatableName)
	assetMethods := l.NewTable()
	assetMetatable.RawSetString("__index", assetMethods)
	assetMetatable.RawSetString("__gc", l.NewFunction(func(L *lua.LState) int {
		asset := checkAudioAsset(L, 1)
		state.releaseAsset(asset)
		return 0
	}))

	instanceMetatable := l.NewTypeMetatable(audioInstanceMetatableName)
	instanceMethods := l.NewTable()
	instanceMetatable.RawSetString("__index", instanceMethods)
	instanceMetatable.RawSetString("__gc", l.NewFunction(func(L *lua.LState) int {
		inst := checkAudioInstance(L, 1)
		state.releaseInstance(inst, true)
		return 0
	}))

	audioTable := l.NewTable()
	l.SetGlobal("audio", audioTable)

	audioTable.RawSetString("load", l.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		resolved := state.resolvePath(path)

		audio, err := state.mixer.LoadAudio(resolved, true)
		if err != nil {
			L.RaiseError("audio.load(%q) failed: %v", resolved, err)
			return 0
		}

		asset := &audioAsset{
			audio:         audio,
			instances:     make(map[*audioInstance]struct{}),
			defaultVolume: 1,
		}
		state.assets[asset] = struct{}{}

		ud := L.NewUserData()
		ud.Value = asset
		L.SetMetatable(ud, L.GetTypeMetatable(audioAssetMetatableName))
		L.Push(ud)
		return 1
	}))

	audioTable.RawSetString("play", l.NewFunction(func(L *lua.LState) int {
		asset := checkAudioAsset(L, 1)
		ensureAssetAlive(L, asset, "audio.play")

		loop, volume := parsePlayOptions(L, 2)
		if L.GetTop() < 2 || L.Get(2).Type() == lua.LTNil {
			volume = asset.defaultVolume
		}

		track, err := state.mixer.CreateTrack()
		if err != nil {
			L.RaiseError("audio.play: failed to create playback track: %v", err)
			return 0
		}

		if err := track.SetAudio(asset.audio); err != nil {
			track.Destroy()
			L.RaiseError("audio.play: failed to set track audio: %v", err)
			return 0
		}

		if loop {
			if err := track.SetLoops(-1); err != nil {
				track.Destroy()
				L.RaiseError("audio.play: failed to set loop mode: %v", err)
				return 0
			}
		}

		inst := &audioInstance{
			asset:  asset,
			track:  track,
			volume: volume,
		}
		state.instances[inst] = struct{}{}
		asset.instances[inst] = struct{}{}

		if err := state.applyInstanceGain(inst); err != nil {
			state.releaseInstance(inst, true)
			L.RaiseError("audio.play: failed to apply volume: %v", err)
			return 0
		}

		if err := track.Play(0); err != nil {
			state.releaseInstance(inst, true)
			L.RaiseError("audio.play: failed to start playback: %v", err)
			return 0
		}

		ud := L.NewUserData()
		ud.Value = inst
		L.SetMetatable(ud, L.GetTypeMetatable(audioInstanceMetatableName))
		L.Push(ud)
		return 1
	}))

	audioTable.RawSetString("stop", l.NewFunction(func(L *lua.LState) int {
		asset, inst := checkAudioTarget(L, 1)
		if asset != nil {
			ensureAssetAlive(L, asset, "audio.stop")
			if err := controlAsset(state, asset, "stop"); err != nil {
				L.RaiseError("audio.stop: failed to stop playback: %v", err)
			}
			return 0
		}

		ensureInstanceAlive(L, inst, "audio.stop")
		if err := controlInstance(state, inst, "stop"); err != nil {
			L.RaiseError("audio.stop: failed to stop playback: %v", err)
		}
		return 0
	}))

	audioTable.RawSetString("pause", l.NewFunction(func(L *lua.LState) int {
		asset, inst := checkAudioTarget(L, 1)
		if asset != nil {
			ensureAssetAlive(L, asset, "audio.pause")
			if err := controlAsset(state, asset, "pause"); err != nil {
				L.RaiseError("audio.pause: failed to pause playback: %v", err)
			}
			return 0
		}

		ensureInstanceAlive(L, inst, "audio.pause")
		if err := controlInstance(state, inst, "pause"); err != nil {
			L.RaiseError("audio.pause: failed to pause playback: %v", err)
		}
		return 0
	}))

	audioTable.RawSetString("resume", l.NewFunction(func(L *lua.LState) int {
		asset, inst := checkAudioTarget(L, 1)
		if asset != nil {
			ensureAssetAlive(L, asset, "audio.resume")
			if err := controlAsset(state, asset, "resume"); err != nil {
				L.RaiseError("audio.resume: failed to resume playback: %v", err)
			}
			return 0
		}

		ensureInstanceAlive(L, inst, "audio.resume")
		if err := controlInstance(state, inst, "resume"); err != nil {
			L.RaiseError("audio.resume: failed to resume playback: %v", err)
		}
		return 0
	}))

	audioTable.RawSetString("set_volume", l.NewFunction(func(L *lua.LState) int {
		asset, inst := checkAudioTarget(L, 1)
		volume := checkAudioVolume(L, 2, "audio.set_volume")

		if asset != nil {
			ensureAssetAlive(L, asset, "audio.set_volume")
			asset.defaultVolume = volume

			for candidate := range asset.instances {
				if candidate.released {
					continue
				}
				candidate.volume = volume
				if err := state.applyInstanceGain(candidate); err != nil {
					L.RaiseError("audio.set_volume: failed to set asset volume: %v", err)
				}
			}
			return 0
		}

		ensureInstanceAlive(L, inst, "audio.set_volume")
		inst.volume = volume
		if err := state.applyInstanceGain(inst); err != nil {
			L.RaiseError("audio.set_volume: failed to set playback volume: %v", err)
		}
		return 0
	}))

	audioTable.RawSetString("get_volume", l.NewFunction(func(L *lua.LState) int {
		asset, inst := checkAudioTarget(L, 1)
		if asset != nil {
			ensureAssetAlive(L, asset, "audio.get_volume")
			L.Push(lua.LNumber(asset.defaultVolume))
			return 1
		}

		ensureInstanceAlive(L, inst, "audio.get_volume")
		L.Push(lua.LNumber(inst.volume))
		return 1
	}))

	audioTable.RawSetString("set_master_volume", l.NewFunction(func(L *lua.LState) int {
		volume := checkAudioVolume(L, 1, "audio.set_master_volume")
		state.masterVolume = volume

		for inst := range state.instances {
			if inst.released {
				continue
			}
			if err := state.applyInstanceGain(inst); err != nil {
				L.RaiseError("audio.set_master_volume: failed to apply master volume: %v", err)
			}
		}

		return 0
	}))

	audioTable.RawSetString("get_master_volume", l.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(state.masterVolume))
		return 1
	}))

	audioTable.RawSetString("is_playing", l.NewFunction(func(L *lua.LState) int {
		asset, inst := checkAudioTarget(L, 1)
		if asset != nil {
			ensureAssetAlive(L, asset, "audio.is_playing")
			L.Push(lua.LBool(isTargetPlaying(state, asset, nil)))
			return 1
		}

		ensureInstanceAlive(L, inst, "audio.is_playing")
		L.Push(lua.LBool(isTargetPlaying(state, nil, inst)))
		return 1
	}))

	audioTable.RawSetString("unload", l.NewFunction(func(L *lua.LState) int {
		asset := checkAudioAsset(L, 1)
		state.releaseAsset(asset)
		return 0
	}))

	audioTable.RawSetString("stop_all", l.NewFunction(func(L *lua.LState) int {
		for inst := range state.instances {
			state.releaseInstance(inst, true)
		}
		return 0
	}))

	return state
}
