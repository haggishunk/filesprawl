# Directory Deletion Safety Checker

## Overview

Before deleting a directory from your storage locations, you need to verify that everything in it is represented elsewhere with proper content. This script automates that safety check.

## Usage

```bash
# Check a specific directory pattern
psql $DATABASE_URL -v "dir_pattern='%travis-transfer/2019/03/10/%'" \
  -f db/scripts/check_directory_deletion_safety.sql
```

Replace the `dir_pattern` with the directory path you want to check. Use wildcards (`%`) for path matching.

## Interpreting Results

The script outputs a safety summary with three categories:

### ✅ ZERO-BYTE (safe)
- Empty or corrupted files (0 bytes)
- Safe to delete unconditionally
- Usually incomplete transfers or system metadata

### ✅ Safe (exists elsewhere)
- Files with copies in other directories
- All copies are non-zero bytes and verified
- Safe to delete from target directory

### ⚠️ RISKY (only copy!)
- Files that **only exist** in the target directory
- **DO NOT DELETE** until you've verified copies elsewhere or backed them up
- These are your only instance of the file

## Example: travis-transfer/2019/03/10/

```
Checking directory pattern: %travis-transfer/2019/03/10/%

         status          | file_count | total_bytes 
-------------------------+------------+-------------
 Safe (exists elsewhere) |         28 |   384544420
 ZERO-BYTE (safe)        |        252 |           0
```

**Verdict**: ✅ **SAFE TO DELETE**
- 252 zero-byte broken files
- 28 non-zero files with verified copies elsewhere
- No files at risk

## Common Use Cases

### Cleaning up incomplete transfer directories
```bash
# Check xforants/light/travis-transfer/
psql $DATABASE_URL -v "dir_pattern='%travis-transfer%'" \
  -f db/scripts/check_directory_deletion_safety.sql
```

### Verifying before archiving to backup
```bash
# Check local old backups
psql $DATABASE_URL -v "dir_pattern='%old-backup%'" \
  -f db/scripts/check_directory_deletion_safety.sql
```

## How It Works

1. **Finds all files** in the target directory
2. **Computes their hashes** (MD5, SHA1, SHA256)
3. **Counts copies elsewhere** with same hash
4. **Categorizes as safe/risky** based on:
   - File size (0-byte = safe)
   - Number of non-target copies
   - Whether copies are verified

## Safety Guarantees

✅ Files marked "Safe" exist elsewhere in your system  
✅ Zero-byte files are assumed corrupted/incomplete  
✅ Only checks database, not actual filesystem  
⚠️ Always verify critical files manually  

## Integration with Consolidation Workflow

This is the first step in the interactive consolidation workflow:
1. **Check safety** (this script)
2. **Archive** the directory
3. **Delete** from original location
4. **Verify** archival and cleanup

See `CONSOLIDATION_WORKFLOW.md` for next steps.
