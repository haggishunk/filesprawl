# Design: Consolidation TUI Redesign - Bubble Tea with Side-by-Side Panels

## Context

Current gotui implementation is monolithic and hard to extend. Bubble Tea provides better composability. Users need clear two-panel workflow: browse duplicates (left) → inspect files (right) → make decisions → execute.

## Goals

1. **Clear Visual Hierarchy** - Banner, 2 panels, status bar all distinct
2. **Keyboard-Driven Workflow** - Complete control without mouse
3. **Data-Rich Display** - Show enough info to make smart deletion decisions
4. **Deletion Safety** - Impossible to delete all copies of a (remote, filename) pair

## Layout (4 Regions)

```
┌─ BANNER ────────────────────────────────────────────────────────┐
│ filesprawl v0.1 | 42 groups | REVIEW MODE | Remotes: s3,dbox   │
├─────────────────────┬──────────────────────────────────────────┤
│                     │                                          │
│ HASH LIST (left)    │  FILE TABLE (right)                    │
│                     │                                          │
│ 📎 Hash 1abc..      │  remote | filename | size | status     │
│    (5 copies)       │  ────────────────────────────────────  │
│ 📎 Hash 2def..      │  s3     | photo.jpg | 2.5M | ✓ 3 locs │
│    (3 copies)       │  dbox   | photo.jpg | 2.5M | ✓ 2 locs │
│ 📎 Hash 3ghi..      │  local  | photo.jpg | 2.5M | 🔴 1 loc │
│    (2 copies)       │                                          │
│                     │  [v]iew wide  [K]eep  [a]rchive       │
├─────────────────────┴──────────────────────────────────────────┤
│ j/k: nav | v: view | K/a/d/s: actions | q: quit | ?: help     │
└──────────────────────────────────────────────────────────────────┘
```

## Components

### 1. Banner (top)

- Version, group count, mode (REVIEW/EXECUTE)
- Remotes list
- Brief keystroke reference

### 2. Hash List Panel (left 40%)

- Duplicate groups listed by hash
- Show: indicator + hash prefix + copy count + location count
- Selection highlights current group
- Keyboard: j/k to scroll, Enter to select

### 3. File Table Panel (right 60%)

- All files in selected hash group
- Columns: remote | filename | size | status
- Optional column: filepath (via 'v' toggle for "wide" view)
- Status: ✓=safe, 🟡=caution, 🔴=critical (single copy)
- Row highlighting shows flags (KEEP/ARCHIVE/DELETE)
- Keyboard: j/k to scroll, space to flag

### 4. Status Bar (bottom)

- Current actions available + deletion safety warning
- Flag summary (X flagged for archive, etc.)

## Deletion Safety Logic

Cannot mark a file for deletion if:

- It's the only copy of (remote, filename) pair
- Example: If hash group has:
  ```
  remote=s3, filename=photo.jpg, size=2.5M  ← Can't delete both
  remote=dbox, filename=photo.jpg, size=2.5M ← Can't delete both
  remote=local, filename=photo.jpg, size=2.5M ← Can delete if other 2 kept
  ```
- Rule: For each (remote, filename) combination, at least 1 copy must not be flagged DELETE

## Keyboard Shortcuts

| Key | Action       | Notes                                 |
| --- | ------------ | ------------------------------------- |
| j/k | Navigate     | Move selection up/down                |
| v   | Toggle wide  | Show/hide filepath column             |
| K   | Flag KEEP    | Keep this copy, suggest delete others |
| a   | Flag ARCHIVE | Archive this copy to backup           |
| d   | Flag DELETE  | Mark for deletion (safety checked)    |
| s   | Skip/Clear   | Clear all flags on current            |
| ?   | Help         | Show keyboard reference               |
| q   | Quit         | Review summary, save flags            |

## Data Flow

1. **Selection** (left panel):
   - User presses j/k to select hash group
   - ConsolidationModel.selectedIndex updated
   - Message: `SelectionMsg{Index: n}`

2. **File Table Update** (right panel):
   - Update message triggers render
   - Table populated from session.DuplicateGroups[selectedIndex].Versions
   - Each row includes flag status from flagManager

3. **Flagging**:
   - User presses K/a/d on file row
   - ActionMsg dispatched with action + fileIndex
   - FlagManager.PersistKeep/Archive/Delete called
   - Table re-renders showing flag status

4. **Deletion Safety Check**:
   - When flagging DELETE, check (remote, filename) pair count
   - If would result in all copies deleted → reject with error message
   - Show count of remaining unflagged copies per remote/filename
