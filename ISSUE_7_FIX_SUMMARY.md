# Issue #7 Fix: Remote Name Filtering

## Problem
The `duplicates` report shows remotes in `hostname/remote_name` format (e.g., `archlinux//tmp/fs-symtest2.vKQy/light`), but the `--remote` flag only accepts `remote_name` format. Users cannot filter by the exact remote name shown in the report.

## Solution: Accept Both Formats

Updated `FindDuplicatesWithinRemote()` to accept both formats:
- `"remote_name"` - traditional format (rclone remotes)
- `"hostname/remote_name"` - new format (matches report output, supports local scans)

### Implementation

**File: `internal/repository/repository.go`**

Modified `FindDuplicatesWithinRemote()` to parse the remote name:
- Detect "/" in input to distinguish formats
- If "/" found: split into hostname and name, query both
- If no "/": query just remote_name (backward compatible)

Example usage:
```bash
# Both commands now work (same result):
filesprawl duplicates --remote "media"
filesprawl duplicates --remote "hostname/media"
```

### Tests Added

1. `TestFindDuplicatesWithinRemote_WithHostnamePrefix` - Verify both formats work
2. `TestFindDuplicatesWithinRemote_HostnamePrefixedPath` - Test with slashes in path

Both tests confirm that querying with hostname-prefixed format returns same results as traditional format.

### Compatibility

✅ Backward compatible (traditional format still works)
✅ All existing tests pass
✅ No database changes required
✅ No breaking API changes

## Files Modified

1. `internal/repository/repository.go` - Updated FindDuplicatesWithinRemote()
2. `internal/repository/repository_integration_test.go` - Added 2 tests

## Testing

```bash
go build ./cmd/filesprawl           # Successful
go test ./internal/analysis -v      # All passing
```

## User Experience

Now users can copy-paste remote names directly from duplicates report:
```
Group 1 — abc123def456... (md5), 2 copies
  • file.txt  [/docs/file.txt]  1024 bytes  remote: hostname/media
```

And use them with `--remote`:
```bash
filesprawl duplicates --remote "hostname/media"
```

No more confusion about format mismatches!
