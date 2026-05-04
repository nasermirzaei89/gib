local x = 0
local y = 0

function game.update(_dt)
    x, y = input.get_mouse_position()
end

function game.render()
    debug.print(20, 20, string.format("mouse: (%d, %d)", x, y))

    local left_down = input.is_mouse_button_down("left")
    local left_pressed = input.is_mouse_button_pressed("left")
    local left_released = input.is_mouse_button_released("left")

    local middle_down = input.is_mouse_button_down("middle")
    local middle_pressed = input.is_mouse_button_pressed("middle")
    local middle_released = input.is_mouse_button_released("middle")

    local right_down = input.is_mouse_button_down("right")
    local right_pressed = input.is_mouse_button_pressed("right")
    local right_released = input.is_mouse_button_released("right")

    debug.print(20, 44, "left down: " .. tostring(left_down))
    debug.print(20, 68, "left pressed: " .. tostring(left_pressed))
    debug.print(20, 92, "left released: " .. tostring(left_released))

    debug.print(20, 132, "middle down: " .. tostring(middle_down))
    debug.print(20, 156, "middle pressed: " .. tostring(middle_pressed))
    debug.print(20, 180, "middle released: " .. tostring(middle_released))

    debug.print(20, 220, "right down: " .. tostring(right_down))
    debug.print(20, 244, "right pressed: " .. tostring(right_pressed))
    debug.print(20, 268, "right released: " .. tostring(right_released))

    debug.print(20, 308, "Try: left, middle, right, x1, x2")
end
