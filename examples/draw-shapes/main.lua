function game.load()
    window_width, window_height = window.get_size()
end

function game.render()
    graphics.draw_rect(40, 40, 180, 120)

    graphics.draw_rect(260, 40, 180, 120, {
        color = {0.2, 0.7, 1.0, 1.0},
        filled = true,
    })

    graphics.draw_line(40, 220, window_width - 40, 220, {
        color = {1.0, 0.3, 0.3, 1.0},
    })

    graphics.draw_line(40, 260, window_width - 40, window_height - 40, {
        color = {0.2, 1.0, 0.5, 1.0},
    })

    graphics.draw_polygon({
        {80, 340},
        {160, 280},
        {260, 340},
    }, {
        color = {1.0, 1.0, 1.0, 1.0},
    })

    graphics.draw_polygon({
        {300, 340},
        {360, 280},
        {460, 320},
        {560, 280},
    }, {
        color = {1.0, 0.6, 0.2, 1.0},
        closed = false,
    })

    graphics.draw_polygon({
        {620, 340},
        {700, 260},
        {780, 340},
        {740, 430},
        {660, 430},
    }, {
        color = {0.8, 0.9, 0.2, 1.0},
        filled = true,
    })

    graphics.draw_circle(120, 520, 36, {
        color = {0.3, 0.9, 1.0, 1.0},
    })

    graphics.draw_circle(240, 520, 36, {
        color = {0.1, 0.7, 1.0, 1.0},
        filled = true,
    })

    graphics.draw_ellipse(380, 520, 56, 30, {
        color = {1.0, 0.7, 0.2, 1.0},
    })

    graphics.draw_ellipse(520, 520, 56, 30, {
        color = {1.0, 0.5, 0.2, 1.0},
        filled = true,
    })

    graphics.draw_arc(680, 520, 46, 0.0, math.pi * 1.5, {
        color = {1.0, 0.3, 0.4, 1.0},
        segments = 32,
    })

    graphics.draw_arc(680, 520, 30, math.pi * 1.25, math.pi * 0.5, {
        color = {0.6, 0.9, 1.0, 1.0},
        segments = 20,
    })
end
