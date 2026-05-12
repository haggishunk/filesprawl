# Change: Add Index Command

## Why
The application currently hardcodes the scan target in `main.go`, which makes indexing useful only as a developer demo. Users need to choose which locally configured rclone remote and optional path to scan from the host running Filesprawl. Persisted scan results also need to record the host computer and the local rclone remote name/config context that the active `rcd` session used so later analysis can clearly attribute indexed files to the correct remote identity.

## What Changes
- Add a user-facing `list-remotes` CLI command that returns the local rclone remote names available through the active `rcd` session
- Add a user-facing `index` CLI command that accepts a required local rclone remote name and an optional remote path
- Treat the selected remote name as a name resolved by the local rclone `rcd` session that Filesprawl is connected to
- Persist remote identity for indexed scans using the host computer identity together with the local rclone remote name
- Persist object-to-remote associations so indexed files can be queried and analyzed by remote origin
- Replace the hardcoded scan target in `main.go` with CLI-driven index execution

## Impact
- Affected specs:
  - New capability `index-command`
  - Modified capability `remote-configuration`
  - Modified capability `object-repository`
  - Modified capability `rclone-integration`
- Affected code:
  - `main.go` for CLI parsing and command dispatch
  - `internal/operation` for scan entrypoint wiring
  - `internal/remote` for persisted remote identity data
  - `internal/repository` for remote read/write and object-to-remote junction persistence
  - `internal/rclone` for `config/listremotes` support and any additional rc metadata lookup needed for persisted remote attributes