package gameengine

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

func eventTypeString(eventType sdl.EventType) string {
	switch eventType {
	case sdl.EVENT_FIRST:
		return "first"
	case sdl.EVENT_QUIT:
		return "quit"
	case sdl.EVENT_TERMINATING:
		return "terminating"
	case sdl.EVENT_LOW_MEMORY:
		return "low_memory"
	case sdl.EVENT_WILL_ENTER_BACKGROUND:
		return "will_enter_background"
	case sdl.EVENT_DID_ENTER_BACKGROUND:
		return "did_enter_background"
	case sdl.EVENT_WILL_ENTER_FOREGROUND:
		return "will_enter_foreground"
	case sdl.EVENT_DID_ENTER_FOREGROUND:
		return "did_enter_foreground"
	case sdl.EVENT_LOCALE_CHANGED:
		return "locale_changed"
	case sdl.EVENT_SYSTEM_THEME_CHANGED:
		return "system_theme_changed"

	case sdl.EVENT_DISPLAY_ORIENTATION:
		return "display_orientation"
	case sdl.EVENT_DISPLAY_ADDED:
		return "display_added"
	case sdl.EVENT_DISPLAY_REMOVED:
		return "display_removed"
	case sdl.EVENT_DISPLAY_MOVED:
		return "display_moved"
	case sdl.EVENT_DISPLAY_DESKTOP_MODE_CHANGED:
		return "display_desktop_mode_changed"
	case sdl.EVENT_DISPLAY_CURRENT_MODE_CHANGED:
		return "display_current_mode_changed"
	case sdl.EVENT_DISPLAY_CONTENT_SCALE_CHANGED:
		return "display_content_scale_changed"
	case sdl.EVENT_DISPLAY_USABLE_BOUNDS_CHANGED:
		return "display_usable_bounds_changed"
	// SDL_EVENT_DISPLAY_FIRST = SDL_EVENT_DISPLAY_ORIENTATION
	// SDL_EVENT_DISPLAY_LAST = SDL_EVENT_DISPLAY_USABLE_BOUNDS_CHANGED

	case sdl.EVENT_WINDOW_SHOWN:
		return "window_shown"
	case sdl.EVENT_WINDOW_HIDDEN:
		return "window_hidden"
	case sdl.EVENT_WINDOW_EXPOSED:
		return "window_exposed"
	case sdl.EVENT_WINDOW_MOVED:
		return "window_moved"
	case sdl.EVENT_WINDOW_RESIZED:
		return "window_resized"
	case sdl.EVENT_WINDOW_PIXEL_SIZE_CHANGED:
		return "window_pixel_size_changed"
	case sdl.EVENT_WINDOW_METAL_VIEW_RESIZED:
		return "window_metal_view_resized"
	case sdl.EVENT_WINDOW_MINIMIZED:
		return "window_minimized"
	case sdl.EVENT_WINDOW_MAXIMIZED:
		return "window_maximized"
	case sdl.EVENT_WINDOW_RESTORED:
		return "window_restored"
	case sdl.EVENT_WINDOW_MOUSE_ENTER:
		return "window_mouse_enter"
	case sdl.EVENT_WINDOW_MOUSE_LEAVE:
		return "window_mouse_leave"
	case sdl.EVENT_WINDOW_FOCUS_GAINED:
		return "window_focus_gained"
	case sdl.EVENT_WINDOW_FOCUS_LOST:
		return "window_focus_lost"
	case sdl.EVENT_WINDOW_CLOSE_REQUESTED:
		return "window_close_requested"
	case sdl.EVENT_WINDOW_HIT_TEST:
		return "window_hit_test"
	case sdl.EVENT_WINDOW_ICCPROF_CHANGED:
		return "window_iccprof_changed"
	case sdl.EVENT_WINDOW_DISPLAY_CHANGED:
		return "window_display_changed"
	case sdl.EVENT_WINDOW_DISPLAY_SCALE_CHANGED:
		return "window_display_scale_changed"
	case sdl.EVENT_WINDOW_SAFE_AREA_CHANGED:
		return "window_safe_area_changed"
	case sdl.EVENT_WINDOW_OCCLUDED:
		return "window_occluded"
	case sdl.EVENT_WINDOW_ENTER_FULLSCREEN:
		return "window_enter_fullscreen"
	case sdl.EVENT_WINDOW_LEAVE_FULLSCREEN:
		return "window_leave_fullscreen"
	case sdl.EVENT_WINDOW_DESTROYED:
		return "window_destroyed"
	case sdl.EVENT_WINDOW_HDR_STATE_CHANGED:
		return "window_hdr_state_changed"
	// SDL_EVENT_WINDOW_FIRST = SDL_EVENT_WINDOW_SHOWN
	// SDL_EVENT_WINDOW_LAST = SDL_EVENT_WINDOW_HDR_STATE_CHANGED

	case sdl.EVENT_KEY_DOWN:
		return "key_down"
	case sdl.EVENT_KEY_UP:
		return "key_up"
	case sdl.EVENT_TEXT_EDITING:
		return "text_editing"
	case sdl.EVENT_TEXT_INPUT:
		return "text_input"
	case sdl.EVENT_KEYMAP_CHANGED:
		return "keymap_changed"
	case sdl.EVENT_KEYBOARD_ADDED:
		return "keyboard_added"
	case sdl.EVENT_KEYBOARD_REMOVED:
		return "keyboard_removed"
	case sdl.EVENT_TEXT_EDITING_CANDIDATES:
		return "text_editing_candidates"
	case sdl.EVENT_SCREEN_KEYBOARD_SHOWN:
		return "screen_keyboard_shown"
	case sdl.EVENT_SCREEN_KEYBOARD_HIDDEN:
		return "screen_keyboard_hidden"

	case sdl.EVENT_MOUSE_MOTION:
		return "mouse_motion"
	case sdl.EVENT_MOUSE_BUTTON_DOWN:
		return "mouse_button_down"
	case sdl.EVENT_MOUSE_BUTTON_UP:
		return "mouse_button_up"
	case sdl.EVENT_MOUSE_WHEEL:
		return "mouse_wheel"
	case sdl.EVENT_MOUSE_ADDED:
		return "mouse_added"
	case sdl.EVENT_MOUSE_REMOVED:
		return "mouse_removed"

	case sdl.EVENT_JOYSTICK_AXIS_MOTION:
		return "joystick_axis_motion"
	case sdl.EVENT_JOYSTICK_BALL_MOTION:
		return "joystick_ball_motion"
	case sdl.EVENT_JOYSTICK_HAT_MOTION:
		return "joystick_hat_motion"
	case sdl.EVENT_JOYSTICK_BUTTON_DOWN:
		return "joystick_button_down"
	case sdl.EVENT_JOYSTICK_BUTTON_UP:
		return "joystick_button_up"
	case sdl.EVENT_JOYSTICK_ADDED:
		return "joystick_added"
	case sdl.EVENT_JOYSTICK_REMOVED:
		return "joystick_removed"
	case sdl.EVENT_JOYSTICK_BATTERY_UPDATED:
		return "joystick_battery_updated"
	case sdl.EVENT_JOYSTICK_UPDATE_COMPLETE:
		return "joystick_update_complete"
	case sdl.EVENT_GAMEPAD_AXIS_MOTION:
		return "gamepad_axis_motion"
	case sdl.EVENT_GAMEPAD_BUTTON_DOWN:
		return "gamepad_button_down"
	case sdl.EVENT_GAMEPAD_BUTTON_UP:
		return "gamepad_button_up"
	case sdl.EVENT_GAMEPAD_ADDED:
		return "gamepad_added"
	case sdl.EVENT_GAMEPAD_REMOVED:
		return "gamepad_removed"
	case sdl.EVENT_GAMEPAD_REMAPPED:
		return "gamepad_remapped"
	case sdl.EVENT_GAMEPAD_TOUCHPAD_DOWN:
		return "gamepad_touchpad_down"
	case sdl.EVENT_GAMEPAD_TOUCHPAD_MOTION:
		return "gamepad_touchpad_motion"
	case sdl.EVENT_GAMEPAD_TOUCHPAD_UP:
		return "gamepad_touchpad_up"
	case sdl.EVENT_GAMEPAD_SENSOR_UPDATE:
		return "gamepad_sensor_update"
	case sdl.EVENT_GAMEPAD_UPDATE_COMPLETE:
		return "gamepad_update_complete"
	case sdl.EVENT_GAMEPAD_STEAM_HANDLE_UPDATED:
		return "gamepad_steam_handle_updated"

	case sdl.EVENT_FINGER_DOWN:
		return "finger_down"
	case sdl.EVENT_FINGER_UP:
		return "finger_up"
	case sdl.EVENT_FINGER_MOTION:
		return "finger_motion"
	case sdl.EVENT_FINGER_CANCELED:
		return "finger_canceled"
	case sdl.EVENT_PINCH_BEGIN:
		return "pinch_begin"
	case sdl.EVENT_PINCH_UPDATE:
		return "pinch_update"
	case sdl.EVENT_PINCH_END:
		return "pinch_end"

	case sdl.EVENT_CLIPBOARD_UPDATE:
		return "clipboard_update"

	case sdl.EVENT_DROP_FILE:
		return "drop_file"
	case sdl.EVENT_DROP_TEXT:
		return "drop_text"
	case sdl.EVENT_DROP_BEGIN:
		return "drop_begin"
	case sdl.EVENT_DROP_COMPLETE:
		return "drop_complete"
	case sdl.EVENT_DROP_POSITION:
		return "drop_position"

	case sdl.EVENT_AUDIO_DEVICE_ADDED:
		return "audio_device_added"
	case sdl.EVENT_AUDIO_DEVICE_REMOVED:
		return "audio_device_removed"
	case sdl.EVENT_AUDIO_DEVICE_FORMAT_CHANGED:
		return "audio_device_format_changed"
	case sdl.EVENT_SENSOR_UPDATE:
		return "sensor_update"

	case sdl.EVENT_PEN_PROXIMITY_IN:
		return "pen_proximity_in"
	case sdl.EVENT_PEN_PROXIMITY_OUT:
		return "pen_proximity_out"
	case sdl.EVENT_PEN_DOWN:
		return "pen_down"
	case sdl.EVENT_PEN_UP:
		return "pen_up"
	case sdl.EVENT_PEN_BUTTON_DOWN:
		return "pen_button_down"
	case sdl.EVENT_PEN_BUTTON_UP:
		return "pen_button_up"
	case sdl.EVENT_PEN_MOTION:
		return "pen_motion"
	case sdl.EVENT_PEN_AXIS:
		return "pen_axis"

	case sdl.EVENT_CAMERA_DEVICE_ADDED:
		return "camera_device_added"
	case sdl.EVENT_CAMERA_DEVICE_REMOVED:
		return "camera_device_removed"
	case sdl.EVENT_CAMERA_DEVICE_APPROVED:
		return "camera_device_approved"
	case sdl.EVENT_CAMERA_DEVICE_DENIED:
		return "camera_device_denied"

	case sdl.EVENT_RENDER_TARGETS_RESET:
		return "render_targets_reset"
	case sdl.EVENT_RENDER_DEVICE_RESET:
		return "render_device_reset"
	case sdl.EVENT_RENDER_DEVICE_LOST:
		return "render_device_lost"

	case sdl.EVENT_PRIVATE0:
		return "private0"
	case sdl.EVENT_PRIVATE1:
		return "private1"
	case sdl.EVENT_PRIVATE2:
		return "private2"
	case sdl.EVENT_PRIVATE3:
		return "private3"

	case sdl.EVENT_POLL_SENTINEL:
		return "poll_sentinel"

	case sdl.EVENT_USER:
		return "user"

	case sdl.EVENT_LAST:
		return "last"

	case sdl.EVENT_ENUM_PADDING:
		return "enum_padding"
	default:
		return "unknown"
	}
}
