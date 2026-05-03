# Getting Started

## Requirements

- Go 1.26+
- SDL3 runtime available on your machine

## Run the Engine

From the project root:

```bash
go run cmd/gib/main.go
```

Behavior:

- No argument: engine looks for `./main.lua` in current directory.
- Argument is a `.lua` file: runs that file.
- Argument is a folder: runs `<folder>/main.lua`.

If no-arg mode does not find `main.lua`, the engine starts an empty window.

## Build Binary

```bash
make build
```

This creates `bin/gib`.

## Minimal Script

```lua
function game.render()
    debug.print(340, 300, "Hello")
end
```

Run it:

```bash
go run cmd/gib/main.go ./path/to/main.lua
```

## Next

- [Game Lifecycle](lifecycle.md)
- [API: graphics](api/graphics.md)
- [Examples](examples.md)
