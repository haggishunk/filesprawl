## Context
Filesprawl already has a scanner and repository pipeline, but the executable still uses a hardcoded remote and path. The project relies on an rclone `rcd` session and, for this change, assumes that session is local to the host running Filesprawl. Under that assumption, the remote name supplied at the CLI is a local rclone configuration name understood by that host's `rcd` session.

## Goals / Non-Goals
- Goals:
  - Provide a CLI command to discover remote names exposed by the active local `rcd` session
  - Provide a minimal CLI entrypoint for indexing a selected remote
  - Preserve the identity of the host and local rclone remote name used for a scan
  - Link persisted objects back to the remote identity that produced them
- Non-Goals:
  - Support remote `rcd` sessions running on another host
  - Add remote enumeration or interactive selection UX
  - Redesign the scan algorithm itself

## Decisions
- Decision: Add a dedicated `list-remotes` CLI command backed by rclone RC `config/listremotes`.
  - Rationale: Users need a reliable way to discover which local remote names the active `rcd` session recognizes before choosing one for indexing.
- Decision: Add a dedicated `index` CLI command with a required remote argument and an optional path argument.
  - Rationale: This is the smallest user-facing interface that removes the hardcoded scan target while pairing naturally with `list-remotes`.
- Decision: Treat the remote CLI argument as a local rclone remote name resolved by the active local `rcd` session.
  - Rationale: This matches the current architecture, where Filesprawl already forwards remote names to rclone RC operations.
- Decision: Use the tuple of host identity and local rclone remote name as the primary persisted remote identity for indexed scans.
  - Rationale: The same remote name may exist on multiple hosts with different local rclone configurations, so persistence must retain host context.
- Decision: Persist remote associations during scan result persistence rather than as a later reconciliation step.
  - Rationale: Downstream duplicate and reporting queries already expect remote associations to exist.

## Risks / Trade-offs
- Remote metadata beyond hostname and remote name may require additional rclone RC inspection.
  - Mitigation: Keep host + remote name as the authoritative identity for this change; treat richer metadata as additive.
- Introducing CLI parsing in `main.go` may make the executable structure more complex.
  - Mitigation: Limit the initial command surface to a single `index` command and keep parsing logic small.

## Migration Plan
1. Add the `index` command and route it into the existing scanner.
2. Add the `list-remotes` command using rclone RC `config/listremotes`.
3. Add repository support for persisted remote identities and object-to-remote links.
4. Update documentation and examples to use the new commands instead of the hardcoded sample flow.

## Open Questions
- Should a future change expose the local hostname as an overrideable CLI option, or should Filesprawl always derive it from the running host?
- Should a future change query the `rcd` session for richer remote configuration details beyond the local remote name?