# Change: Migrate to Cobra CLI Framework

## Why
The current CLI implementation uses manual flag parsing with the standard library `flag` package, which results in:
- No built-in help generation (`--help` must be manually implemented)
- No subcommand structure (everything is in a flat switch statement)
- Manual validation and error messages for each command
- Inconsistent help formatting across commands
- No support for advanced features like command aliases, persistent flags, or command suggestions
- Poor user experience compared to modern CLI tools

Migrating to Cobra (https://github.com/spf13/cobra) would provide:
- Automatic help generation with consistent formatting
- Nested subcommand support with proper structure
- Built-in flag inheritance (persistent flags)
- Shell completion generation (bash, zsh, fish, PowerShell)
- Command aliases and suggestions for typos
- Industry-standard CLI patterns (used by kubectl, docker, gh, hugo, etc.)
- Better maintainability and extensibility

## What Changes
- Add `github.com/spf13/cobra` dependency to go.mod
- Restructure main.go to use Cobra's command pattern
- Create root command with version, help, and completion support
- Convert existing commands (list-remotes, index, duplicates) to Cobra commands
- Add proper help text, examples, and flag descriptions
- Implement shell completion generation
- Maintain backward compatibility with existing command syntax

## Impact
- Affected specs:
  - Modified capability `duplicates-cli` (improved help and validation)
  - New capability `cobra-cli` (framework integration)
- Affected code:
  - `main.go` - Complete restructure to use Cobra commands
  - Existing command handlers remain largely unchanged (business logic stays same)
  - Flag definitions move to Cobra's flag system
  - Better error messages and user experience
- Users benefit from:
  - Professional help output with `filesprawl --help`, `filesprawl duplicates --help`, etc.
  - Shell completion for commands and flags
  - Better error messages with suggestions
  - Consistent CLI experience matching other modern tools
- Development benefits:
  - Easier to add new commands
  - Less boilerplate code
  - Better testing support
  - Standard patterns for CLI applications

