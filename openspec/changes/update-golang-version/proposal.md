# Change: Upgrade Go toolchain to 1.26

## Why
The project currently targets Go 1.21.3 (released October 2023), which has been
out of the two-release support window since Go 1.23 shipped. Go 1.26 (released
February 2026) is the current stable release and brings performance improvements
from the now-default Green Tea garbage collector, reduced cgo overhead (~30%),
improved stack-allocation of slices, and security fixes accumulated across five
major releases. Staying on 1.21.3 also means missing language conveniences
(built-in `min`/`max`, loop-variable-per-iteration semantics, range-over-int,
iterators) that simplify and harden the codebase.

## What Changes

- Update `go` directive in `go.mod` from `1.21.3` to `1.26`
- Run `go get -u ./...` and `go mod tidy` to update all direct and indirect
  dependencies to versions compatible with Go 1.26 (pgx/v5, rclone, mapstructure)
- Remove the hand-written `min()` helper in `internal/analysis/analysis.go` and
  use the built-in `min` (available since Go 1.21, now idiomatic to rely on)
- Verify all `for`-range loops are safe under Go 1.22's per-iteration variable
  semantics (audit shows no unsafe captures; no functional changes needed)
- Update `project.md` to reflect the new minimum Go version requirement
- Run `go fix ./...` to apply any standard moderniser suggestions from the
  rewritten `go fix` tool shipped with Go 1.26

## Impact

- Affected specs: `database-persistence` (Database Interface Abstraction
  requirement updated to include the `Query` method formalised during
  add-duplicate-detection)
- Affected code:
  - `go.mod` and `go.sum` (toolchain directive + dependency versions)
  - `internal/analysis/analysis.go` (remove custom `min`)
  - `openspec/project.md` (Go version reference)
- No breaking changes to public APIs or database schema
- All existing tests must continue to pass after the upgrade

