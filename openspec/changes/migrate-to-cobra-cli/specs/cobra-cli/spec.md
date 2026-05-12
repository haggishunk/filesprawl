# Cobra CLI Framework

## Purpose
Professional command-line interface framework providing automatic help generation, subcommand structure, and modern CLI patterns.

## ADDED Requirements

### Requirement: Root Command
The system SHALL provide a root command with version, help, and metadata.

#### Scenario: Help output
- **WHEN** running `filesprawl --help` or `filesprawl -h`
- **THEN** display formatted help with available commands, flags, and examples

#### Scenario: Version output
- **WHEN** running `filesprawl --version` or `filesprawl version`
- **THEN** display the application version

#### Scenario: No command
- **WHEN** running `filesprawl` without arguments
- **THEN** display help output automatically

### Requirement: Subcommand Structure
The system SHALL organize functionality into well-defined subcommands.

#### Scenario: List remotes command
- **WHEN** running `filesprawl list-remotes --help`
- **THEN** show command-specific help with description and examples

#### Scenario: Index command
- **WHEN** running `filesprawl index --help`
- **THEN** show index command help with required and optional flags

#### Scenario: Duplicates command
- **WHEN** running `filesprawl duplicates --help`
- **THEN** show duplicates command help with all filtering options

### Requirement: Automatic Help Generation
The system SHALL automatically generate help text for all commands and flags.

#### Scenario: Command help includes description
- **WHEN** viewing command help
- **THEN** display the command's purpose and usage

#### Scenario: Command help includes examples
- **WHEN** viewing command help
- **THEN** display practical usage examples

#### Scenario: Flag help includes descriptions
- **WHEN** viewing command help
- **THEN** list all flags with their types, defaults, and descriptions

#### Scenario: Flag help shows required flags
- **WHEN** viewing command help
- **THEN** clearly mark which flags are required

### Requirement: Flag Management
The system SHALL use Cobra's flag system for consistent flag handling.

#### Scenario: Persistent flags
- **WHEN** defining flags like `--database-url` on root command
- **THEN** make them available to all subcommands

#### Scenario: Local flags
- **WHEN** defining command-specific flags
- **THEN** only show them in that command's help

#### Scenario: Flag validation
- **WHEN** invalid flag values are provided
- **THEN** return clear error messages with expected format

#### Scenario: Flag shortcuts
- **WHEN** defining flags
- **THEN** support both long form (`--remote`) and short form (`-r`) where appropriate

### Requirement: Shell Completion
The system SHALL support shell completion generation.

#### Scenario: Bash completion
- **WHEN** running `filesprawl completion bash`
- **THEN** output bash completion script

#### Scenario: Zsh completion
- **WHEN** running `filesprawl completion zsh`
- **THEN** output zsh completion script

#### Scenario: Fish completion
- **WHEN** running `filesprawl completion fish`
- **THEN** output fish completion script

#### Scenario: PowerShell completion
- **WHEN** running `filesprawl completion powershell`
- **THEN** output PowerShell completion script

### Requirement: Error Handling
The system SHALL provide helpful error messages using Cobra's error handling.

#### Scenario: Unknown command
- **WHEN** running `filesprawl unknowncommand`
- **THEN** show error with command suggestions if similar commands exist

#### Scenario: Missing required flag
- **WHEN** running command without required flags
- **THEN** show clear error indicating which flag is required

#### Scenario: Invalid flag value
- **WHEN** providing invalid flag value
- **THEN** show error with expected value format

#### Scenario: Mutually exclusive flags
- **WHEN** providing conflicting flags
- **THEN** show error explaining the conflict

### Requirement: Command Examples
The system SHALL provide usage examples for each command.

#### Scenario: List remotes examples
- **WHEN** viewing `filesprawl list-remotes --help`
- **THEN** show example: `filesprawl list-remotes`

#### Scenario: Index examples
- **WHEN** viewing `filesprawl index --help`
- **THEN** show examples for basic and path-specific indexing

#### Scenario: Duplicates examples
- **WHEN** viewing `filesprawl duplicates --help`
- **THEN** show examples for within-remote, across-remotes, and filtered queries

### Requirement: Backward Compatibility
The system SHALL maintain existing command syntax.

#### Scenario: Existing commands work unchanged
- **WHEN** running `filesprawl list-remotes`
- **THEN** behavior is identical to current implementation

#### Scenario: Existing flags work unchanged
- **WHEN** running `filesprawl index --remote media --path code/flux`
- **THEN** behavior is identical to current implementation

#### Scenario: Existing duplicates command works
- **WHEN** running `filesprawl duplicates --remote media --min-size 1000000`
- **THEN** behavior is identical to current implementation

### Requirement: Exit Codes
The system SHALL return appropriate exit codes via Cobra's error handling.

#### Scenario: Successful execution
- **WHEN** command completes successfully
- **THEN** return exit code 0

#### Scenario: User error
- **WHEN** command fails due to user input
- **THEN** return exit code 1

#### Scenario: System error
- **WHEN** command fails due to system error
- **THEN** return exit code 1

