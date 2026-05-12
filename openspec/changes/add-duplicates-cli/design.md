# Design: Add Duplicates CLI Command

## Overview
This change adds a `duplicates` command to the filesprawl CLI, exposing the existing duplicate detection backend (`internal/analysis` package) to end users. The command provides two primary modes: finding duplicates within a single remote, and finding duplicates across all remotes.

## Command Syntax

### Basic Usage
```bash
# Find duplicates within a specific remote
filesprawl duplicates --remote <remote-name>

# Find duplicates across all remotes
filesprawl duplicates --across-remotes
```

### With Filters
```bash
# Find large duplicates (>10MB) across remotes
filesprawl duplicates --across-remotes --min-size 10485760

# Find duplicates by specific hash type
filesprawl duplicates --remote media --hash-type sha256

# Paginate results
filesprawl duplicates --across-remotes --limit 20 --offset 0

# Combined filters
filesprawl duplicates --across-remotes --hash-type dropbox --min-size 5000000 --limit 10
```

## Architecture

### Data Flow
```
User Command
    ↓
parseDuplicatesArgs() → duplicatesArgs
    ↓
runDuplicates()
    ↓
openRepository() → ObjectRepository
    ↓
NewDuplicateDetector(repo) → DuplicateDetector
    ↓
detector.FindWithinRemote() OR detector.FindAcrossRemotes()
    ↓
analysis.FormatReport(groups)
    ↓
Write to stdout
```

### Key Types

```go
type duplicatesArgs struct {
    Remote        string  // Remote name for within-remote mode
    AcrossRemotes bool    // Flag for across-remotes mode
    HashType      string  // Optional: md5, sha1, sha256, dropbox
    MinSize       int64   // Optional: minimum file size in bytes
    Limit         int     // Optional: max duplicate groups to return
    Offset        int     // Optional: skip N groups for pagination
}
```

## Example Outputs

### Within Remote (No Duplicates)
```
$ filesprawl duplicates --remote media
No duplicates found.
```

### Within Remote (With Duplicates)
```
$ filesprawl duplicates --remote dbox-haggis-at-gmail --limit 3
Found 1097 duplicate group(s):

Group 1 — 1cfced9b90dd... (dropbox), 2 copies
  • Achi_Koji_Emily.mov  [McDonald Talent Show 2023/Achi_Koji_Emily.mov]  257698816 bytes  remote: archlinux/dbox-haggis-at-gmail
  • Achi_Koji_Emily.mov  [McDonald Talent Show 2023 my copy/McDonald Talent Show 2023/Achi_Koji_Emily.mov]  257698816 bytes  remote: archlinux/dbox-haggis-at-gmail

Group 2 — aa423158e0e2... (dropbox), 2 copies
  • Odin Talent Show 2023.mov  [Odin Talent Show 2023.mov]  173015040 bytes  remote: archlinux/dbox-haggis-at-gmail
  • Odin Talent Show 2023.mov  [McDonald Talent Show 2023 my copy/Odin Talent Show 2023.mov]  173015040 bytes  remote: archlinux/dbox-haggis-at-gmail

Group 3 — d7d55ec7667e... (dropbox), 2 copies
  • IMG_5311.MOV  [IMG_5311.MOV]  80740352 bytes  remote: archlinux/dbox-haggis-at-gmail
  • IMG_5311.MOV  [Cara Mattera Halloween 10.22/IMG_5311.MOV]  80740352 bytes  remote: archlinux/dbox-haggis-at-gmail
```

### Across Remotes
```
$ filesprawl duplicates --across-remotes --min-size 1000000
Found 15 duplicate group(s):

Group 1 — 9b7a2bf6b0dd... (sha256), 3 copies
  • backup.tar.gz  [/backups/2024/backup.tar.gz]  45000000 bytes  remote: server1/s3-archive
  • backup.tar.gz  [/archive/backup.tar.gz]  45000000 bytes  remote: server2/dropbox-backup
  • backup.tar.gz  [/local/backup.tar.gz]  45000000 bytes  remote: server1/local-nas

...
```

### With Specific Hash Type
```
$ filesprawl duplicates --remote media --hash-type md5
Found 234 duplicate group(s):

Group 1 — a1b2c3d4e5f6... (md5), 2 copies
  • photo.jpg  [/photos/2024/photo.jpg]  2048576 bytes  remote: laptop/media
  • photo.jpg  [/photos/backup/photo.jpg]  2048576 bytes  remote: laptop/media

...
```

## Decisions

### Decision: Mutually Exclusive Modes
**Rationale**: Users must explicitly choose between finding duplicates within one remote versus across all remotes. This prevents ambiguity and ensures clear user intent.

### Decision: Reuse Existing Backend
**Rationale**: The `internal/analysis` package already provides all necessary duplicate detection logic. The CLI layer is purely an interface to expose this functionality.

### Decision: Use FormatReport() Directly
**Rationale**: The existing `FormatReport()` function produces well-formatted, human-readable output suitable for CLI display without modification.

### Decision: No Default Mode
**Rationale**: Requiring explicit mode selection prevents accidental expensive operations (like scanning all remotes) and makes the user's intention clear.

### Decision: Filter Validation in Parse Phase
**Rationale**: Validating hash type and numeric values during argument parsing provides immediate feedback to users before any database operations occur.

### Decision: Pagination via Limit/Offset
**Rationale**: For users with large duplicate sets, limit and offset allow iterating through results in manageable chunks without loading everything at once.

## Edge Cases

### Empty Database
- **Scenario**: No files have been indexed yet
- **Behavior**: Return "No duplicates found."

### Single File Per Hash
- **Scenario**: All indexed files have unique hashes
- **Behavior**: Return "No duplicates found."

### All Duplicates Filtered Out
- **Scenario**: Filters eliminate all duplicate groups (e.g., --min-size too large)
- **Behavior**: Return "No duplicates found."

### Remote Name Not Found
- **Scenario**: Specified remote doesn't exist in database
- **Behavior**: May return "No duplicates found." or empty result set (repository behavior)

### Invalid Hash Type
- **Scenario**: User provides `--hash-type invalid`
- **Behavior**: Return error: "invalid hash type 'invalid', must be one of: md5, sha1, sha256, dropbox"

### Both Modes Specified
- **Scenario**: User provides both `--remote media` and `--across-remotes`
- **Behavior**: Return error: "--remote and --across-remotes are mutually exclusive"

### Neither Mode Specified
- **Scenario**: User runs `filesprawl duplicates` without mode flags
- **Behavior**: Return error: "must specify either --remote or --across-remotes"

## Testing Strategy

### Unit Tests
- Argument parsing with valid inputs
- Argument validation (hash types, numeric ranges)
- Mutually exclusive flag detection
- Remote name normalization

### Integration Tests
- Full command execution with mocked repository
- Output format verification
- Error message verification
- Exit code verification

### Manual Testing
- Run against actual indexed database with known duplicates
- Verify output matches expected duplicate sets
- Test all filter combinations
- Test pagination with limit/offset

## Future Enhancements (Out of Scope)

### JSON Output Format
Add `--json` flag for machine-readable output:
```bash
filesprawl duplicates --across-remotes --json
```

### Summary Statistics
Add `--summary` flag showing statistics without file lists:
```bash
filesprawl duplicates --across-remotes --summary
# Output: 
# Total duplicate groups: 1097
# Total duplicate files: 2426
# Total wasted space: 9.1 GB
```

### Interactive Mode
Add interactive selection for removing duplicates:
```bash
filesprawl duplicates --across-remotes --interactive
# Prompts user to choose which copies to keep/remove
```

### CSV Export
Add `--csv` flag for spreadsheet import:
```bash
filesprawl duplicates --across-remotes --csv > duplicates.csv
```

These enhancements would require separate change proposals and specs.

