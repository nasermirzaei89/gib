function game.load()
    sprite_sheet = graphics.load_image("../assets/eagle.png")
    sprite_width = 160
    sprite_height = 320
    sprite_index = 0
    sprite_index_f = 0
    image_speed = 5
    image_count = 6

    movement_speed = 250

    mirror_horizontal = false

    window_width, window_height = window.get_size()

    x = window_width / 2 - sprite_width / 2
    y = window_height / 2 - sprite_height / 2
end

function game.update(dt)
    sprite_index_f = (sprite_index_f + image_speed * dt) % image_count
    sprite_index = math.floor(sprite_index_f)

    local x_speed = 0
    local y_speed = 0

    if input.is_key_down("left") then
        x_speed = x_speed - movement_speed * dt
    end
    if input.is_key_down("right") then
        x_speed = x_speed + movement_speed * dt
    end
    if input.is_key_down("up") then
        y_speed = y_speed - movement_speed * dt
    end
    if input.is_key_down("down") then
        y_speed = y_speed + movement_speed * dt
    end

    x = x + x_speed
    y = y + y_speed

    if x_speed < 0 then
        mirror_horizontal = true
    elseif x_speed > 0 then
        mirror_horizontal = false
    end

    if x > window_width then
        x = -sprite_width
    end
    if x < -sprite_width then
        x = window_width
    end
    if y > window_height then
        y = -sprite_height
    end
    if y < -sprite_height then
        y = window_height
    end
end

function game.render()
    graphics.draw_image(sprite_sheet, x, y, {
        sx = sprite_index * sprite_width,
        sy = 0,
        sw = sprite_width,
        sh = sprite_height,
        scale_x = mirror_horizontal and -1 or 1,
    })
end
