# Data Loss Prevention for File Consolidation

## Three Critical Safety Checks

Before deleting or consolidating any files, run these three checks in order:

### 1. Detect Name Collisions
**Script**: `db/scripts/detect_name_collisions.sql`

**Problem**: Different files with same filename (e.g., two photos named "photo.jpg")

**Check**: Verify all variants are preserved before consolidating by filename
```bash
psql $DATABASE_URL -f db/scripts/detect_name_collisions.sql
```

**Result Interpretation**:
- ✅ OK: All unique variants preserved → Safe to consolidate
- ⚠️ WARNING: Only 1 variant preserved → **DATA LOSS RISK!**
- 🚨 CRITICAL: Variants completely missing → **UNRECOVERABLE DATA LOSS!**

---

### 2. Check Deletion Safety
**Script**: `db/scripts/check_directory_deletion_safety.sql`

**Problem**: Deleting files that don't exist elsewhere

**Check**: Verify all files in a directory exist elsewhere before deletion
```bash
psql $DATABASE_URL -v "dir_pattern='%target/dir/%'" \
  -f db/scripts/check_directory_deletion_safety.sql
```

**Result Interpretation**:
- ✅ ZERO-BYTE (safe): Corrupted/incomplete files
- ✅ Safe (exists elsewhere): Verified copies with non-zero bytes
- ⚠️ RISKY (only copy!): **DO NOT DELETE!**

---

### 3. Verify by Hash and Size
**Problem**: Same hash but 0-byte files (incomplete downloads)

**Check**: Before consolidating, ensure:
```sql
-- For each hash:
-- 1. Zero-byte copies → safe to delete
-- 2. Non-zero copies → must exist elsewhere before deleting
SELECT 
  hash_value,
  COUNT(DISTINCT CASE WHEN size > 0 THEN file_id END) as nonzero_copies,
  COUNT(DISTINCT CASE WHEN size = 0 THEN file_id END) as zero_copies
FROM files
GROUP BY hash_value
HAVING COUNT(*) > 1;
```

---

## Safe Consolidation Workflow

### Step 1: Analyze
```bash
# Check for name collisions
psql $DATABASE_URL -f db/scripts/detect_name_collisions.sql

# Identify target directory
cd /path/to/consolidate
```

### Step 2: Verify Safety
```bash
# Check if all files exist elsewhere
psql $DATABASE_URL -v "dir_pattern='%/path/to/consolidate/%'" \
  -f db/scripts/check_directory_deletion_safety.sql
```

### Step 3: Decision
If any:
- ⚠️ WARNING collisions → Don't consolidate by filename
- 🚨 CRITICAL variants → DON'T DELETE!
- ⚠️ RISKY files → DON'T DELETE!

Otherwise:
- ✅ Safe to proceed with consolidation

### Step 4: Archive
```bash
# Archive the directory (don't delete yet)
tar czf archive-$(date +%Y%m%d).tar.gz /path/to/consolidate
```

### Step 5: Verify Archive
```bash
# Verify archive integrity
tar tzf archive-$(date +%Y%m%d).tar.gz | wc -l
```

### Step 6: Delete
```bash
# Only delete after verified archive
rm -rf /path/to/consolidate
```

---

## Your Data: Analysis Summary

### Filename Collisions
**Result**: ✅ **ALL SAFE**
- AlbumArtSmall.jpg: 224 unique versions (all preserved)
- Folder.jpg: 224 unique versions (all preserved)
- Music files: 2-3 variants each (all preserved)
- Photos: 13-15 variants each (all preserved)

### Zero-Byte Files
**Result**: ✅ **ALL SAFE**
- travis-transfer/2019/03/10: 252 zero-byte (safe to delete)
- 28 non-zero files with copies elsewhere (safe to delete)

**Verdict**: ✅ Your data is protected! Safe consolidation possible.

---

## Pro Tips

1. **Always archive before deleting**
   - Preserves ability to recover
   - Verify archive before deleting originals

2. **Use hash-based consolidation**
   - Hash: unique identifier for file content
   - Name: can have collisions (different files)
   - Size: helps distinguish versions

3. **Keep metadata**
   - Store consolidation log with timestamps
   - Track which files were deleted
   - Store archive location

4. **Partial path matching**
   - Use source directory as secondary identifier
   - "delete photo.jpg if copy exists in official-backup"
   - Combines name, path, and hash for safety

---

## References

See also:
- `DELETION_SAFETY.md` — Directory deletion safety checker
- `NAME_COLLISIONS.md` — Filename collision detection
- `duplicates` command — Find actual duplicates by hash
