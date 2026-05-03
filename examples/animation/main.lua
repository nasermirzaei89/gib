function game.load()
    sprite_sheet = graphics.load_image("../assets/eagle.png")
    sprite_width = 160
    sprite_height = 320
    sprite_index = 0
    sprite_index_f = 0
    image_speed = 0.15
    image_count = 6

    movement_speed = 5

    window_width, window_height = window.get_size()

    x = window_width / 2 - sprite_width / 2
    y = window_height / 2 - sprite_height / 2
end

function game.fixed_update()
    sprite_index_f = (sprite_index_f + image_speed) % image_count
    sprite_index = math.floor(sprite_index_f)

    x = x + movement_speed
    if x > window_width then
        x = -sprite_width
    end
end

function game.render()
    graphics.draw_image(sprite_sheet, x, y, {
        sx = sprite_index * sprite_width,
        sy = 0,
        sw = sprite_width,
        sh = sprite_height
    })
end
