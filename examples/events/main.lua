local event_handlers = {
    key_down = function(event)
        log.info(string.format("Event received. type=%s key=%s is_repeat=%s", event.type, event.key, tostring(event.is_repeat)))
    end,
    key_up = function(event)
        log.info(string.format("Event received. type=%s key=%s", event.type, event.key))
    end,
    mouse_motion = function(event)
        log.info(string.format("Event received. type=%s x=%.1f y=%.1f dx=%.1f dy=%.1f", event.type, event.x, event.y, event.dx, event.dy))
    end,
    mouse_button_down = function(event)
        log.info(string.format("Event received. type=%s button=%s clicks=%d x=%.1f y=%.1f", event.type, event.button, event.clicks, event.x, event.y))
    end,
    mouse_button_up = function(event)
        log.info(string.format("Event received. type=%s button=%s clicks=%d x=%.1f y=%.1f", event.type, event.button, event.clicks, event.x, event.y))
    end,
    mouse_wheel = function(event)
        log.info(string.format("Event received. type=%s x=%.1f y=%.1f mouse_x=%.1f mouse_y=%.1f", event.type, event.x, event.y, event.mouse_x, event.mouse_y))
    end,
}

function game.event(event)
    local handler = event_handlers[event.type]
    if handler ~= nil then
        handler(event)
        return
    end

    log.info("Event received. type=" .. event.type)
end

function game.render()
    debug.print(280, 300, "Watch the console for events!")
end
