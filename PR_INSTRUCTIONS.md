# Pull Request Instructions

## Branch Created and Pushed ✅

**Branch Name:** `feat/local-scanning-locality-interface`

**Commit:** ca21b69
```
feat: implement local file scanning and locality interface
```

**Pushed to:** `origin/feat/local-scanning-locality-interface`

## Create Pull Request

Visit the GitHub UI to create the PR:

https://github.com/haggishunk/filesprawl/pull/new/feat/local-scanning-locality-interface

### PR Details

**Title:**
```
feat: implement local file scanning and locality interface
```

**Description:**
See the full PR body template below or use the auto-generated GitHub template.

## Changes Summary

### Local File Scanning
- CLI: `filesprawl scan-local /path`
- Direct filesystem scanning (no rclone)
- MD5, SHA1, SHA256 hashing
- Integrated with duplicate detection

### Locality Interface
- CLI: `filesprawl locality resolve-remote|resolve-local|show-config`
- Bidirectional path mapping
- Pattern matching with variables
- Configuration with caching

## Testing Status
- ✅ 26 unit tests passing
- ✅ Build successful
- ✅ OpenSpec validation passing
- ✅ No breaking changes

## Files Changed
- 8 new files
- 3 modified files
- 1,141 insertions (+)
- 49 deletions (-)

## Next Steps
1. Create PR on GitHub using the link above
2. Add reviewers if needed
3. Wait for CI checks to pass
4. Merge to main when ready
