function game.event(event)
    if event.type == "key_down" or event.type == "key_up" then
        log.info(string.format("%s key=%s scancode=%d is_repeat=%s", event.type, event.key, event.scancode, tostring(event.is_repeat)))
        return
    end

    if event.type == "mouse_motion" then
        log.info(string.format("Mouse motion x=%.1f y=%.1f dx=%.1f dy=%.1f", event.x, event.y, event.dx, event.dy))
        return
    end

    if event.type == "mouse_button_down" or event.type == "mouse_button_up" then
        log.info(string.format("%s button=%d clicks=%d x=%.1f y=%.1f", event.type, event.button, event.clicks, event.x, event.y))
        return
    end

    if event.type == "mouse_wheel" then
        log.info(string.format("Mouse wheel x=%.1f y=%.1f mouse_x=%.1f mouse_y=%.1f", event.x, event.y, event.mouse_x, event.mouse_y))
        return
    end

    log.info("Event received: " .. event.type)
end

function game.render()
    debug.print(280, 300, "Watch the console for events!")
end
