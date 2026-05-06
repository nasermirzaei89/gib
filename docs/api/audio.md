# audio API

## load(path)

```lua
local sound = audio.load("shoot.wav")
local music = audio.load("bgm.ogg")
```

Loads an audio asset and returns asset userdata.

Notes:

- Supports at least WAV and OGG through SDL_mixer.
- Relative paths are resolved from the script base directory.

## unload(asset)

```lua
audio.unload(sound)
```

Releases an audio asset.

Behavior:

- Stops and invalidates all playback instances created from that asset.
- Calling `audio.unload` repeatedly on the same asset is allowed.

## play(asset, opts?)

```lua
local hit = audio.play(sound)

local bgm = audio.play(music, {
  loop = true,
  volume = 1.0,
})
```

Starts playback for an asset and returns playback userdata.

Options:

- `loop` boolean (default `false`)
- `volume` number in `[0, 1]` (default is the asset volume)

Validation behavior:

- Unknown option keys raise errors.
- `volume` must be in `[0, 1]`.

## stop(target)
## pause(target)
## resume(target)

```lua
audio.stop(sound)      -- stop all active instances of this asset
audio.pause(sound)     -- pause all active instances of this asset
audio.resume(sound)    -- resume all paused instances of this asset

audio.stop(hit)        -- stop one playback instance
audio.pause(hit)       -- pause one playback instance
audio.resume(hit)      -- resume one playback instance
```

Control playback for either an asset target or playback-instance target.

Behavior:

- Redundant transitions for valid handles are no-ops.
- Invalidated/released handles raise errors.

## set_volume(target, volume)
## get_volume(target)

```lua
audio.set_volume(sound, 0.5)
local v1 = audio.get_volume(sound)

audio.set_volume(hit, 0.8)
local v2 = audio.get_volume(hit)
```

Sets/gets volume for asset or playback instance target.

Behavior:

- Asset target changes the default volume for future plays.
- Asset target also applies to currently active instances.
- Volume must be in `[0, 1]`.

## set_master_volume(volume)
## get_master_volume()

```lua
audio.set_master_volume(0.5)
local master = audio.get_master_volume()
```

Sets/gets global master volume in `[0, 1]`.

Behavior:

- Applied immediately to active playback and future playback.
- Final output gain is:

```txt
final_volume = master_volume * instance_volume
```

## is_playing(target)

```lua
if audio.is_playing(music) then
  -- at least one music instance is playing
end

if audio.is_playing(hit) then
  -- this instance is currently playing
end
```

Returns `true` when target is currently playing.

## stop_all()

```lua
audio.stop_all()
```

Stops and invalidates all playback instances.

## TBD

See [Roadmap](../roadmap.md#audio).
