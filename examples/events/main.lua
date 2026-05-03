function game.event(event)
    log.info("Event received: " .. event.type)
end

function game.render()
    debug.print(280, 300, "Watch the console for events!")
end
