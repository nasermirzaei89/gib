# AGENTS.md

## Available commands

- `make build`: Compiles the game engine and produces an executable in the `bin` directory.

## Lua

- **Variables & Functions:** Use `snake_case` (e.g., `local player_score = 0`). This is favored because it provides high readability and distinguishes user code from the standard library.
- **Booleans:** Prefix with `is_` or `has_` (e.g., `is_active`, `has_item`).
- **Classes/Factories:** Use `PascalCase` (e.g., `local PlayerAccount = {}`).
- **Constants:** Use `LOUD_SNAKE_CASE` (e.g., `local MAX_PLAYERS = 10`). Note that Lua doesn't have true constants, so this is purely a visual hint.
- **Internal/Private Members:** Prefix with a single underscore (e.g., `_private_var`).
  - **Warning:** Avoid starting names with an underscore followed by uppercase letters (like `_VERSION`), as these are reserved for Lua's internal use.
- **File Names:** Use `snake_case.lua` (e.g., `player_account.lua`).
