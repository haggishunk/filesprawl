# Change: Consolidation TUI UI Redesign - Bubble Tea with Side-by-Side Panels

## Why

The current gotui-based UI is hard to extend and lacks clear visual hierarchy. The new design using Bubble Tea provides:

- Cleaner Elm-inspired architecture (Model/Update/View)
- Two side-by-side panels (hash list + file table) for better data exploration
- Improved keyboard interaction model with typed messages
- Extensible component architecture for future features

## What Changes

- **Migrate from gotui to Bubble Tea** - More composable TUI framework
- **Redesign layout** - Banner top + 2-column layout (hash list left, file table right)
- **Hash panel** - List of duplicate groups (default), selectable by hash
- **File panel** - Table display showing (remote, filename) with optional filepath column
- **Deletion safety** - Can't delete without keeping at least 1 copy per (remote, filename) pair
- **New component system** - Reusable Bubble Tea components for lists, tables, modals
- **Keyboard-driven** - All navigation and actions via keyboard shortcuts

## Impact

- Affected specs: `consolidation-browsing` (UI presentation details)
- Affected code: `internal/consolidation/ui/*` (all TUI code)
- Breaking change: Old gotui widgets/screens replaced with Bubble Tea model
- Data model unchanged: Session, flags, database persistence unaffected
