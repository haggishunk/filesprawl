# Implementation Tasks

## 1. CLI Command Surface
- [x] 1.1 Define the `list-remotes` command interface for enumerating remotes from the active local `rcd` session
- [x] 1.2 Define the `index` command interface with a required remote argument and optional path argument
- [x] 1.3 Replace the hardcoded scan target in `main.go` with CLI-based command dispatch
- [x] 1.4 Validate CLI inputs and return helpful usage errors when the remote argument is missing

## 2. rclone RC Integration
- [x] 2.1 Add support for rclone RC `config/listremotes`
- [x] 2.2 Parse and return the configured remote names from the active local `rcd` session
- [x] 2.3 Add unit coverage for remote-list parsing and error handling

## 3. Remote Identity Persistence
- [x] 3.1 Extend the remote domain/repository model to read or write persisted remote identities for indexed scans
- [x] 3.2 Define the persisted remote identity around the local host and the local rclone remote name used by the active `rcd` session
- [x] 3.3 Persist `object_remote_junction` records during scan result persistence
- [x] 3.4 Add unit and integration coverage for remote identity persistence and object-to-remote associations

## 4. Scan Execution Integration
- [x] 4.1 Pass the CLI-selected remote and path into the existing scan pipeline
- [x] 4.2 Ensure persisted scan results remain attributable to the selected host/remote identity across repeated scans
- [x] 4.3 Document local `rcd` assumptions and provide example command invocations

## 5. Validation
- [x] 5.1 Run `go test ./...`
- [x] 5.2 Run `openspec validate add-index-command --strict`