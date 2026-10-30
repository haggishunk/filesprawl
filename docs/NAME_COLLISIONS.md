# Name Collision Detection for Safe Consolidation

## Overview

When consolidating files by name (e.g., "delete duplicates.txt"), you risk data loss if different files happen to share the same name. This detector identifies these collisions and flags risks.

## Common Scenarios

### ✅ Safe: Multi-Format Collections
```
01 Intro.mp3  (from album A, FLAC version)
01 Intro.mp3  (from album A, MP3 version)
01 Intro.mp3  (from album B, completely different song)
```
If all 3 exist in your system → Safe to consolidate by hash, but NOT by filename!

### ⚠️ Risky: Corrupted Downloads
```
photo.jpg  (complete, 2.5 MB, hash: abc123)
photo.jpg  (incomplete, 0 bytes, hash: d41d8cd98f00b204e9800998ecf8427e)
photo.jpg  (good copy elsewhere, 2.5 MB, hash: abc123)
```
Safe only if you keep the complete versions!

### 🚨 CRITICAL: Lost Variants
```
song.mp3  (version A, hash: abc123)
song.mp3  (version B, hash: xyz789) - MISSING!
```
Only 1 variant preserved → Data loss!

## Usage

```bash
# Detect all filename collisions
psql $DATABASE_URL -f db/scripts/detect_name_collisions.sql
```

## Interpreting Results

### ✅ OK: All variants preserved
```
AlbumArtSmall.jpg | 224 unique versions | 224 preserved | status: ✅ OK
```
Safe to consolidate **by hash**, but be careful with **filename-based** operations.

### ⚠️ WARNING: Only 1 version preserved
```
cover.jpg | 5 unique variants | 1 preserved | status: ⚠️ WARNING
```
4 different photos named "cover.jpg" are lost → **CRITICAL data loss!**

### 🚨 CRITICAL: Some variants missing
```
photo.jpg | 3 unique variants | 2 preserved | status: 🚨 CRITICAL
```
1 photo version is completely missing from all locations!

## Safe Consolidation Rules

### Rule 1: Hash-based operations are safe
```sql
-- Safe: Delete by hash, keeps all variants
DELETE WHERE hash = 'abc123' AND size > 0 AND other_copies_exist
```

### Rule 2: Never delete by filename alone
```sql
-- UNSAFE: Might delete only variant of "photo.jpg"!
DELETE WHERE object_name = 'photo.jpg'
```

### Rule 3: Use path hints for disambiguation
```sql
-- Safer: Delete from specific directory if all variants exist elsewhere
DELETE WHERE object_name = 'photo.jpg' 
  AND path LIKE '%old-backup%'
  AND all_variants_exist_in_other_dirs
```

## Integration with Deletion Safety

Combine with `DELETION_SAFETY.md`:

1. **Run name collision detector** ← Current step
   - Find files with same name but different content
   - Verify all variants are preserved

2. **Run deletion safety checker**
   - Verify all files in target directory exist elsewhere
   - Check for 0-byte copies

3. **Archive and delete with confidence**
   - Safe to remove since all variants preserved

## Example: Your Data

Your music library has **hundreds of multi-variant files**:
- `AlbumArtSmall.jpg`: 224 different versions (all preserved ✅)
- `Folder.jpg`: 224 different versions (all preserved ✅)
- CIMG photos: 13-15 variants each (all preserved ✅)

**Verdict**: Safe to consolidate by hash, but use hash-based logic, not filename-based!

## Pro Tip: Partial Path Matching

When consolidating, use partial path matching for confidence:

```sql
-- Example: Keep music from preferred source
DELETE FROM files 
WHERE object_name = 'song.mp3'
  AND path LIKE '%low-quality-rip%'  -- specific source
  AND hash IN (
    SELECT hash FROM files 
    WHERE object_name = 'song.mp3' 
      AND path LIKE '%official-release%'  -- same song exists in better source
  )
```

This combines filename, path, and hash for safe consolidation.
