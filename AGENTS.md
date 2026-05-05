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

## Commit Message Guidelines

When creating commits, follow these conventions:

### Commit Title

- Use natural language that completes the phrase: **"[This commit will] ..."**
- Start with a capital letter
- Use imperative mood (e.g., "Add", "Fix", "Update", "Remove")
- Keep it concise (50-72 characters recommended)
- Do NOT use conventional commit prefixes like `feat:`, `fix:`, `test:`, `chore:`, etc.

**Good examples:**

- "Add comprehensive repository tests for SQLite implementations"
- "Refactor database layer to use direct sql.DB instances"
- "Update email field naming to emailAddress across codebase"

**Bad examples:**

- "test: add tests" (uses prefix)
- "added some stuff" (not imperative, too vague)
- "Fix bug" (too vague, no context)

### Commit Description

- Include a blank line between title and description
- Use bullet points to describe what changed
- Be specific about the changes made
- Group related changes together
- Focus on WHAT changed and WHY, not HOW

**Example:**

```txt
Add comprehensive repository tests for SQLite implementations

- Add OTP repository tests covering CRUD operations, validation, and edge cases
- Add session repository tests for user sessions and expiration handling
- Add user repository tests for all user operations and status verification
- Tests cover multiple scenarios including resends, multiple sessions, and cleanup operations
```
