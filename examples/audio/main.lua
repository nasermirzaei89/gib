local MUSIC_PATH = "../assets/455516__ispeakwaves__the-plan-upbeat-loop-no-voice-edit-mono-track.ogg"
local GUNSHOT_PATH = "../assets/91572__steveygos93__layeredgunshot.wav"
local VOLUME_STEP = 0.1

local music_asset
local gunshot_asset
local music_playback
local is_music_paused = false

local function clamp(v, min_v, max_v)
    if v < min_v then
        return min_v
    end

    if v > max_v then
        return max_v
    end

    return v
end

local function start_music()
    music_playback = audio.play(music_asset, { loop = true })
    is_music_paused = false
end

function game.load()
    music_asset = audio.load(MUSIC_PATH)
    gunshot_asset = audio.load(GUNSHOT_PATH)

    audio.set_master_volume(1.0)
    start_music()
end

function game.update(dt)
    if input.is_key_pressed("p") then
        if is_music_paused then
            audio.resume(music_playback)
            is_music_paused = false
        elseif music_playback == nil or not audio.is_playing(music_playback) then
            start_music()
        else
            audio.pause(music_playback)
            is_music_paused = true
        end
    end

    if input.is_key_pressed("s") and music_playback ~= nil then
        audio.stop(music_playback)
        is_music_paused = false
    end

    if input.is_key_pressed("space") then
        audio.play(gunshot_asset)
    end

    if input.is_key_pressed("up") then
        local master_volume = audio.get_master_volume()
        audio.set_master_volume(clamp(master_volume + VOLUME_STEP, 0.0, 1.0))
    end

    if input.is_key_pressed("down") then
        local master_volume = audio.get_master_volume()
        audio.set_master_volume(clamp(master_volume - VOLUME_STEP, 0.0, 1.0))
    end
end

function game.render()
    local music_state = "stopped"
    if is_music_paused then
        music_state = "paused"
    elseif music_playback ~= nil and audio.is_playing(music_playback) then
        music_state = "playing"
    end

    local master_volume = audio.get_master_volume()

    debug.print(20, 20, "Audio example")
    debug.print(20, 48, "P: pause/resume music")
    debug.print(20, 76, "S: stop music")
    debug.print(20, 104, "Space: play gunshot")
    debug.print(20, 132, "Up/Down: master volume")
    debug.print(20, 176, string.format("Music state: %s", music_state))
    debug.print(20, 204, string.format("Master volume: %.1f", master_volume))
end
