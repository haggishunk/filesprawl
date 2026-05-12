# Implementation Tasks

## 1. Command Interface
- [ ] 1.1 Add `duplicates` case to the command switch in `runCommand()`
- [ ] 1.2 Define `duplicatesArgs` struct with fields for remote, acrossRemotes, hashType, minSize, limit, offset
- [ ] 1.3 Create `parseDuplicatesArgs()` function to parse CLI flags
- [ ] 1.4 Validate mutually exclusive `--remote` and `--across-remotes` flags
- [ ] 1.5 Add usage help text for `duplicates` command
- [ ] 1.6 Update main usage string to include `duplicates` command

## 2. Command Handler
- [ ] 2.1 Create `runDuplicates()` function that accepts context, args, and stdout
- [ ] 2.2 Call `openRepository()` to get database connection
- [ ] 2.3 Create `analysis.NewDuplicateDetector()` with repository
- [ ] 2.4 Build `analysis.FilterOptions` from parsed arguments
- [ ] 2.5 Call appropriate detector method based on mode (within remote vs across remotes)
- [ ] 2.6 Use `analysis.FormatReport()` to format results
- [ ] 2.7 Write formatted report to stdout

## 3. Argument Validation
- [ ] 3.1 Validate hash type is one of: md5, sha1, sha256, dropbox (or empty for all)
- [ ] 3.2 Validate min-size is non-negative integer
- [ ] 3.3 Validate limit is non-negative integer
- [ ] 3.4 Validate offset is non-negative integer
- [ ] 3.5 Normalize remote name using existing `normalizeRemoteName()` function
- [ ] 3.6 Return clear error messages for invalid inputs

## 4. Error Handling
- [ ] 4.1 Handle missing DATABASE_URL with clear error message
- [ ] 4.2 Handle database connection failures with context
- [ ] 4.3 Handle repository query errors with wrapped context
- [ ] 4.4 Handle missing mode selection (neither --remote nor --across-remotes)
- [ ] 4.5 Handle conflicting mode selection (both --remote and --across-remotes)

## 5. Testing
- [ ] 5.1 Write unit tests for `parseDuplicatesArgs()` with valid inputs
- [ ] 5.2 Write unit tests for argument validation edge cases
- [ ] 5.3 Write unit tests for mutually exclusive flag detection
- [ ] 5.4 Write integration test for within-remote duplicates command
- [ ] 5.5 Write integration test for across-remotes duplicates command
- [ ] 5.6 Write integration test with filtering options
- [ ] 5.7 Write integration test for "no duplicates found" case
- [ ] 5.8 Test error paths (missing DATABASE_URL, invalid remote, etc.)

## 6. Documentation
- [ ] 6.1 Update README.md with duplicates command examples
- [ ] 6.2 Document all available flags and their purpose
- [ ] 6.3 Provide example outputs showing typical duplicate reports
- [ ] 6.4 Document common use cases (finding large duplicates, specific hash types)

## 7. Integration
- [ ] 7.1 Build and test locally with actual database
- [ ] 7.2 Run against live indexed data to verify output format
- [ ] 7.3 Test with various filter combinations
- [ ] 7.4 Verify exit codes in success and error cases
- [ ] 7.5 Run `go test ./...` to ensure no regressions
- [ ] 7.6 Run `openspec validate add-duplicates-cli --strict`

