function game.load()
    window_width, window_height = window.get_size()

    image = graphics.load_image("../assets/leo.png")
    image_width = 160
    image_height = 320
    x = window_width / 2 - image_width / 2
    y = window_height / 2 - image_height / 2
end

function game.render()
    graphics.draw_image(image, x, y)
end
