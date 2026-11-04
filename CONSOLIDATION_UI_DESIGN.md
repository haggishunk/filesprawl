# Consolidation TUI Design: Bubble Tea Side-by-Side Panels

## Overview

The new consolidation TUI uses **Bubble Tea** with a **2-column layout**:

- **Left (40%)**: Hash list - unique duplicate groups, selectable
- **Right (60%)**: File table - all files in selected group, with remote/filename/size/status

## Layout

```
┌─────────────────────────────────────────────────────────────────┐
│ filesprawl v0.1 | Groups: 42 | REVIEW | Remotes: s3,dbox      │
├─────────────────────┬──────────────────────────────────────────┤
│ HASH LIST           │ FILE TABLE                               │
│ (j/k scroll, v)     │ (j/k scroll, space flag, v wide)        │
│                     │                                          │
│ 📎 abc12...         │ remote | filename | size | status        │
│    5 copies, 3 loc  │ ────────────────────────────────────── │
│ 📎 def34...         │ s3     | photo.jpg | 2.5M | ✓ 3 locs  │
│    3 copies, 2 loc  │ dbox   | photo.jpg | 2.5M | ✓ 2 locs  │
│ 📎 ghi56...         │ local  | photo.jpg | 2.5M | 🔴 1 loc  │
│    2 copies, 1 loc  │                                          │
├─────────────────────┴──────────────────────────────────────────┤
│ j/k: nav | v: wide | K/a/d/s: actions | q: quit | ?: help     │
└─────────────────────────────────────────────────────────────────┘
```

## Key Features

### Hash List Panel (Left)

- Lists unique duplicate groups by content hash
- Shows: 🟢/🟡/🔴 indicator + hash + copy count + location count
- Selection with j/k keyboard navigation
- Default view on startup

### File Table Panel (Right)

- Displays all files in selected hash group
- Columns: remote | filename | size | status
- Status: ✓=safe, 🟡=caution, 🔴=critical
- Optional filepath column (toggle with 'v')
- Shows flag status (✓ KEEP, ▶ ARCHIVE, ✗ DELETE)

## Deletion Safety

**Core Rule**: Cannot delete ALL copies of a (remote, filename) pair.

Example:

```
Hash: abc123
├─ (s3, photo.jpg, 2.5M)       ← DELETE ✓ allowed (copies on dbox + local)
├─ (dbox, photo.jpg, 2.5M)     ← DELETE ✓ allowed (copies on s3 + local)
└─ (local, photo.jpg, 2.5M)    ← DELETE ✓ allowed (copies on s3 + dbox)

But: If s3 + dbox both deleted, cannot delete local's only copy
```

## Keyboard Shortcuts

| Key | Panel | Action                               |
| --- | ----- | ------------------------------------ |
| j/k | Both  | Scroll/Navigate                      |
| v   | Right | Toggle filepath column (wide view)   |
| K   | Right | Mark file as KEEP                    |
| a   | Right | Mark file as ARCHIVE                 |
| d   | Right | Mark file as DELETE (safety checked) |
| s   | Right | Skip/clear flags                     |
| ?   | Any   | Show help modal                      |
| q   | Any   | Quit and review summary              |

## Data Flow

1. **Model State**: `selectedIndex`, `isWideView`, `modalType`, dimensions
2. **Messages**: Typed (not string-based)
   - `ScrollMsg{Direction}`
   - `SelectionMsg{Index}`
   - `ViewSwitchMsg{}`
   - `ActionMsg{Action, FileIndex}`
3. **Update**: Pure function processes message, updates model, returns re-render cmd
4. **View**: Renders banner + panels + status based on model state

## Safety Logic

Before allowing DELETE flag:

1. Get file's (remote, filename) pair
2. Count all files with same (remote, filename) in group
3. Count how many are NOT flagged DELETE
4. If count <= 1: REJECT with error message
5. If count > 1: ALLOW and update flag

## Component Structure

```
ConsolidationModel (main)
├── Banner (read-only info)
├── ListPane (left - selectable list)
├── TablePane (right - flaggable table)
└── StatusBar (read-only, shows commands + summary)
```

Each component:

- Pure render function
- Keyboard event handler
- State updates through messages

## Next Steps (Phase 2+)

1. Implement ListPane component
2. Implement TablePane component
3. Wire keyboard handlers
4. Build safety checker
5. Test with real data
