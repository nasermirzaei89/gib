local image
local image_width = 0
local image_height = 0
local t = 0

function game.load()
    image = graphics.load_image("../assets/leo.png")
    image_width, image_height = image:get_size()
    graphics.set_clear_color({0.3, 0.3, 0.3, 1.0})
end

function game.update(dt)
    t = t + dt
end

function game.render()
    local cx = 220
    local cy = 170
    local image_x = cx - image_width / 2
    local image_y = cy - image_height / 2

    local mirror = 1
    if math.sin(t * 1.2) < 0 then
        mirror = -1
    end

    graphics.draw_image(image, image_x, image_y, {
        rotation = t,
        origin = {image_width / 2, image_height / 2},
        scale = {mirror * 1.2, 1.2},
    })

    local rw = 180
    local rh = 100
    graphics.draw_rect(420, 110, rw, rh, {
        color = {0.2, 0.8, 1.0, 1.0},
        filled = false,
        rotation = -t * 0.8,
        origin = {rw / 2, rh / 2},
        scale = {1.0 + 0.25 * math.sin(t * 2.0), 1.0},
    })

    graphics.draw_ellipse(640, 160, 90, 45, {
        color = {1.0, 0.7, 0.2, 1.0},
        rotation = t * 0.7,
        origin = {0, 0},
        scale = {1.0, 1.0},
    })

    graphics.draw_arc(640, 360, 80, 0.0, math.pi * 1.5, {
        color = {1.0, 0.3, 0.4, 1.0},
        segments = 40,
        rotation = t,
        origin = {0, 0},
        scale = {1.0, 1.0},
    })

    debug.print(20, 20, "Transforms: rotation/origin/scale via opts")
end
