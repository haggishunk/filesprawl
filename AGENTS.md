<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

## Repository-specific instructions

- When developing or testing `filesprawl` command implementations that talk to rclone RC, prefer exercising them with `rclone rc <command>` against the active RC server so behavior is verified against the same command surface Filesprawl uses.
- In particular, use `rclone rc` calls such as `rclone rc config/listremotes` and `rclone rc operations/list fs=<remote> remote=<path>` to confirm request/response expectations before or while updating Filesprawl CLI behavior.