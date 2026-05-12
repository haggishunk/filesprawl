# Duplicates CLI

## Purpose
Command-line interface for discovering and analyzing duplicate files within and across indexed remote storage locations.

## ADDED Requirements

### Requirement: Duplicates Command Interface
The system SHALL provide a `duplicates` command that accepts mode and filtering options.

#### Scenario: Command help
- **WHEN** running `filesprawl duplicates --help` or `filesprawl duplicates` without arguments
- **THEN** display usage information showing available modes and options

#### Scenario: Required mode selection
- **WHEN** running `duplicates` command
- **THEN** require either `--remote <name>` or `--across-remotes` flag to specify the search scope

#### Scenario: Mutually exclusive modes
- **WHEN** both `--remote` and `--across-remotes` are provided
- **THEN** return an error indicating these options are mutually exclusive

### Requirement: Within Remote Duplicates
The system SHALL support finding duplicates within a specific remote.

#### Scenario: Find duplicates in single remote
- **WHEN** running `filesprawl duplicates --remote <name>`
- **THEN** display duplicate groups where files share hashes within that remote only

#### Scenario: Remote name validation
- **WHEN** running `filesprawl duplicates --remote ""`
- **THEN** return an error that remote name cannot be empty

#### Scenario: Remote normalization
- **WHEN** running `filesprawl duplicates --remote media` or `--remote media:`
- **THEN** both forms are accepted and normalized consistently

### Requirement: Across Remotes Duplicates
The system SHALL support finding duplicates across all indexed remotes.

#### Scenario: Find duplicates across all remotes
- **WHEN** running `filesprawl duplicates --across-remotes`
- **THEN** display duplicate groups where files exist on multiple distinct remotes

#### Scenario: Cross-remote semantics
- **WHEN** displaying cross-remote duplicates
- **THEN** only include hash groups that span at least two distinct remote identities

### Requirement: Filtering Options
The system SHALL support optional filters to refine duplicate detection results.

#### Scenario: Filter by hash type
- **WHEN** running `filesprawl duplicates --across-remotes --hash-type sha256`
- **THEN** only include duplicates identified by SHA256 hashes

#### Scenario: Valid hash types
- **WHEN** specifying `--hash-type` option
- **THEN** accept: `md5`, `sha1`, `sha256`, `dropbox`

#### Scenario: Filter by minimum file size
- **WHEN** running `filesprawl duplicates --across-remotes --min-size 1048576`
- **THEN** only include duplicate files larger than or equal to 1MB (1,048,576 bytes)

#### Scenario: Limit result count
- **WHEN** running `filesprawl duplicates --across-remotes --limit 10`
- **THEN** return at most 10 duplicate groups

#### Scenario: Result offset for pagination
- **WHEN** running `filesprawl duplicates --across-remotes --offset 20`
- **THEN** skip the first 20 duplicate groups and return subsequent results

#### Scenario: Combined filters
- **WHEN** running `filesprawl duplicates --remote media --hash-type dropbox --min-size 10000000 --limit 5`
- **THEN** apply all filters: remote=media, hash type=dropbox, minimum size=10MB, limit to 5 groups

### Requirement: Output Format
The system SHALL use human-readable output format for duplicate reports.

#### Scenario: Use existing FormatReport
- **WHEN** displaying duplicate results
- **THEN** use the `analysis.FormatReport()` function for consistent formatting

#### Scenario: Group header format
- **WHEN** displaying each duplicate group
- **THEN** show: group number, hash prefix, hash type, and copy count

#### Scenario: File detail format
- **WHEN** displaying files within a group
- **THEN** show: file name, path, size in bytes, and remote identifier (hostname/remote-name)

#### Scenario: No duplicates found
- **WHEN** no duplicates match the criteria
- **THEN** display "No duplicates found."

#### Scenario: Empty result after filtering
- **WHEN** filters eliminate all duplicate groups
- **THEN** display "No duplicates found."

### Requirement: Database Connection
The system SHALL require database connectivity for duplicate detection.

#### Scenario: DATABASE_URL required
- **WHEN** running `duplicates` command without DATABASE_URL environment variable
- **THEN** return error "DATABASE_URL is required"

#### Scenario: Database connection failure
- **WHEN** database connection fails
- **THEN** return descriptive error message with connection details

### Requirement: Error Handling
The system SHALL provide clear error messages for invalid inputs and failures.

#### Scenario: Invalid hash type
- **WHEN** running `filesprawl duplicates --across-remotes --hash-type invalid`
- **THEN** return error listing valid hash types

#### Scenario: Invalid numeric values
- **WHEN** providing non-numeric values for `--min-size`, `--limit`, or `--offset`
- **THEN** return error indicating expected numeric value

#### Scenario: Negative numeric values
- **WHEN** providing negative values for `--min-size`, `--limit`, or `--offset`
- **THEN** return error indicating values must be non-negative

#### Scenario: Repository query failure
- **WHEN** repository query fails
- **THEN** return error with context about which operation failed

### Requirement: Exit Codes
The system SHALL return appropriate exit codes.

#### Scenario: Successful execution
- **WHEN** command completes successfully
- **THEN** return exit code 0

#### Scenario: User error (invalid arguments)
- **WHEN** command fails due to invalid arguments or usage
- **THEN** return exit code 1

#### Scenario: System error (database, repository)
- **WHEN** command fails due to system errors
- **THEN** return exit code 1

