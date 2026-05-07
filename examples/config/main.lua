function game.config(conf)
    conf.window.width = 1280
    conf.window.height = 720
    conf.window.title = "Config Example"
    conf.window.resizable = true
    conf.window.fullscreen = false

    conf.tps = 60
end

function game.render()
    debug.print(430, 340, "Config callback applied.")
end
