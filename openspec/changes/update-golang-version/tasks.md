# Implementation Tasks

## 1. Toolchain Update
- [ ] 1.1 Update `go` directive in `go.mod` to `go 1.26`
- [ ] 1.2 Run `go get -u ./...` to upgrade direct and indirect dependencies
- [ ] 1.3 Run `go mod tidy` to prune unused dependencies and sync `go.sum`
- [ ] 1.4 Confirm build succeeds: `go build ./...`

## 2. Code Modernisation
- [ ] 2.1 Remove the hand-written `min(a, b int) int` helper from
        `internal/analysis/analysis.go` and replace call site with the builtin `min`
- [ ] 2.2 Run `go fix ./...` and review any suggestions from the Go 1.26 modernisers;
        apply safe fixes, discard noise
- [ ] 2.3 Audit all `for`-range loops for implicit variable capture (Go 1.22 semantics);
        confirm no behavioural changes are introduced

## 3. Verification
- [ ] 3.1 Run full test suite: `go test ./...`
- [ ] 3.2 Run `go vet ./...` and resolve any new warnings
- [ ] 3.3 Confirm `go build ./...` produces a clean binary

## 4. Documentation
- [ ] 4.1 Update the Go version reference in `openspec/project.md`
        from `1.21.3` to `1.26`

