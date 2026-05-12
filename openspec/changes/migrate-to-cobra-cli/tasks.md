# Implementation Tasks

## 1. Add Cobra Dependency
- [x] 1.1 Run `go get github.com/spf13/cobra@latest`
- [x] 1.2 Verify cobra dependency in go.mod
- [x] 1.3 Run `go mod tidy` to clean up dependencies

## 2. Create Root Command
- [x] 2.1 Create `cmd/root.go` with root command definition
- [x] 2.2 Add application description and version information
- [ ] 2.3 Add persistent flags (e.g., --database-url if needed)
- [ ] 2.4 Configure help and usage templates if custom formatting needed
- [x] 2.5 Initialize root command in main.go

## 3. Create List Remotes Command
- [x] 3.1 Create `cmd/list_remotes.go` with listRemotesCmd
- [x] 3.2 Add command description and examples
- [x] 3.3 Move runListRemotes logic to command's Run function
- [x] 3.4 Register command with root via init()
- [x] 3.5 Test `filesprawl list-remotes --help`

## 4. Create Index Command
- [x] 4.1 Create `cmd/index.go` with indexCmd
- [x] 4.2 Add command description and examples
- [x] 4.3 Define --remote flag (required, string)
- [x] 4.4 Define --path flag (optional, string)
- [x] 4.5 Move runIndex logic to command's Run function
- [x] 4.6 Move parseIndexArgs validation to PreRunE
- [x] 4.7 Register command with root via init()
- [x] 4.8 Test `filesprawl index --help`

## 5. Create Duplicates Command
- [x] 5.1 Create `cmd/duplicates.go` with duplicatesCmd
- [x] 5.2 Add command description and examples
- [x] 5.3 Define --remote flag (optional, string)
- [x] 5.4 Define --across-remotes flag (optional, bool)
- [x] 5.5 Define --hash-type flag (optional, string with validation)
- [x] 5.6 Define --min-size flag (optional, int64)
- [x] 5.7 Define --limit flag (optional, int)
- [x] 5.8 Define --offset flag (optional, int)
- [x] 5.9 Move runDuplicates logic to command's Run function
- [x] 5.10 Move parseDuplicatesArgs validation to PreRunE
- [x] 5.11 Implement custom flag validation for mutually exclusive modes
- [x] 5.12 Register command with root via init()
- [x] 5.13 Test `filesprawl duplicates --help`

## 6. Add Shell Completion
- [x] 6.1 Add completion command via cobra's built-in completion
- [ ] 6.2 Test `filesprawl completion bash`
- [ ] 6.3 Test `filesprawl completion zsh`
- [ ] 6.4 Test `filesprawl completion fish`
- [ ] 6.5 Test `filesprawl completion powershell`
- [ ] 6.6 Add completion installation instructions to README

## 7. Refactor main.go
- [x] 7.1 Simplify main.go to only initialize and execute root command
- [x] 7.2 Remove manual flag parsing code
- [x] 7.3 Remove manual command switch statement
- [x] 7.4 Keep openRepository and other utility functions
- [x] 7.5 Consider moving utility functions to internal package

## 8. Add Version Command
- [x] 8.1 Define version constant or read from build info
- [x] 8.2 Add version command or --version flag to root
- [x] 8.3 Test `filesprawl --version` and `filesprawl version`

## 9. Testing
- [x] 9.1 Run all existing tests to ensure no regressions
- [x] 9.2 Test all commands with --help flag
- [ ] 9.3 Test backward compatibility with existing command syntax
- [ ] 9.4 Test error messages for invalid inputs
- [ ] 9.5 Test mutually exclusive flag validation
- [ ] 9.6 Test required flag validation
- [ ] 9.7 Test shell completion generation
- [ ] 9.8 Manually test actual command execution (list-remotes, index, duplicates)

## 10. Documentation
- [ ] 10.1 Update README.md with new help command examples
- [ ] 10.2 Document shell completion installation
- [ ] 10.3 Update command examples to show --help usage
- [ ] 10.4 Add version information to README
- [ ] 10.5 Document any changes to command syntax (if any)

## 11. Cleanup
- [x] 11.1 Remove unused parseIndexArgs function (if fully moved to Cobra)
- [x] 11.2 Remove unused parseDuplicatesArgs function (if fully moved to Cobra)
- [x] 11.3 Remove manual error message strings (use Cobra's built-in)
- [x] 11.4 Run `go mod tidy`
- [ ] 11.5 Run `go fmt ./...`
- [x] 11.6 Run `go vet ./...`

## 12. Integration
- [x] 12.1 Build and test locally: `go build -o build/filesprawl`
- [ ] 12.2 Test with actual database and rclone setup
- [ ] 12.3 Verify all commands work as expected
- [ ] 12.4 Test error scenarios
- [x] 12.5 Run `go test ./...`
- [ ] 12.6 Run `openspec validate migrate-to-cobra-cli --strict`

