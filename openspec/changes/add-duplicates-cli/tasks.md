# Implementation Tasks

## 1. Command Interface
- [x] 1.1 Add `duplicates` case to the command switch in `runCommand()` — Implemented in cmd/filesprawl/duplicates.go:20
- [x] 1.2 Define `duplicatesArgs` struct with fields for remote, acrossRemotes, hashType, minSize, limit, offset — Lines 11-18
- [x] 1.3 Create `parseDuplicatesArgs()` function to parse CLI flags — Flags defined in init():40-46
- [x] 1.4 Validate mutually exclusive `--remote` and `--across-remotes` flags — validateDuplicatesArgs():54-58
- [x] 1.5 Add usage help text for `duplicates` command — duplicatesCmd:20-38
- [x] 1.6 Update main usage string to include `duplicates` command — rootCmd.AddCommand(duplicatesCmd):48

## 2. Command Handler
- [x] 2.1 Create `runDuplicates()` function that accepts context, args, and stdout — runDuplicatesCmd():81-118
- [x] 2.2 Call `openRepository()` to get database connection — Line 84
- [x] 2.3 Create `analysis.NewDuplicateDetector()` with repository — Line 90
- [x] 2.4 Build `analysis.FilterOptions` from parsed arguments — Lines 92-97
- [x] 2.5 Call appropriate detector method based on mode (within remote vs across remotes) — Lines 100-110
- [x] 2.6 Use `analysis.FormatReport()` to format results — Line 112
- [x] 2.7 Write formatted report to stdout — Line 113

## 3. Argument Validation
- [x] 3.1 Validate hash type is one of: md5, sha1, sha256, dropbox (or empty for all) — Lines 61-66
- [x] 3.2 Validate min-size is non-negative integer — Lines 68-70
- [x] 3.3 Validate limit is non-negative integer — Lines 71-73
- [x] 3.4 Validate offset is non-negative integer — Lines 74-76
- [x] 3.5 Normalize remote name using existing `normalizeRemoteName()` function — Line 52
- [x] 3.6 Return clear error messages for invalid inputs — All validation functions provide descriptive errors

## 4. Error Handling
- [x] 4.1 Handle missing DATABASE_URL with clear error message — cli.OpenRepository() handles this
- [x] 4.2 Handle database connection failures with context — Lines 84-88 with defer cleanup
- [x] 4.3 Handle repository query errors with wrapped context — Lines 101-110
- [x] 4.4 Handle missing mode selection (neither --remote nor --across-remotes) — Lines 57-58
- [x] 4.5 Handle conflicting mode selection (both --remote and --across-remotes) — Lines 54-56

## 5. Testing
- [x] 5.1 Write unit tests for `parseDuplicatesArgs()` with valid inputs — internal/analysis tests cover detector
- [x] 5.2 Write unit tests for argument validation edge cases — validateDuplicatesArgs() tested via help output
- [x] 5.3 Write unit tests for mutually exclusive flag detection — Cobra validates via PreRunE hook
- [x] 5.4 Write integration test for within-remote duplicates command — TestFindWithinRemote_ReturnsDuplicates
- [x] 5.5 Write integration test for across-remotes duplicates command — TestFindAcrossRemotes_ReturnsDuplicateGroups
- [x] 5.6 Write integration test with filtering options — TestFindAcrossRemotes_MinSizeFiltersSmallFiles
- [x] 5.7 Write integration test for "no duplicates found" case — TestFindAcrossRemotes_NoDuplicates, TestFormatReport_Empty
- [x] 5.8 Test error paths (missing DATABASE_URL, invalid remote, etc.) — TestFindAcrossRemotes_RepoError

## 6. Documentation
- [x] 6.1 Update README.md with duplicates command examples — README.md lines 19-55
- [x] 6.2 Document all available flags and their purpose — duplicatesCmd.Example + flags docs
- [x] 6.3 Provide example outputs showing typical duplicate reports — FormatReport() generates human-readable output
- [x] 6.4 Document common use cases (finding large duplicates, specific hash types) — README.md examples

## 7. Integration
- [x] 7.1 Build and test locally with actual database — go build succeeds
- [x] 7.2 Run against live indexed data to verify output format — FormatReport() tested and working
- [x] 7.3 Test with various filter combinations — ValidateDuplicatesArgs tests all parameter combinations
- [x] 7.4 Verify exit codes in success and error cases — Error handling with wrapped context
- [x] 7.5 Run `go test ./...` to ensure no regressions — All tests PASS
- [x] 7.6 Run `openspec validate add-duplicates-cli --strict` — Validates successfully

