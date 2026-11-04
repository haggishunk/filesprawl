# Phase 2 Implementation Details

## BaseListWidget with gotui.List

```go
type BaseListWidget struct {
    list       *widgets.List        // Real gotui List widget
    title      string
    rows       []string
    x, y, width, height int         // Rectangle dimensions
}

// GetWidget returns the actual gotui List for rendering
func (w *BaseListWidget) GetWidget() interface{} {
    if w == nil || w.list == nil {
        return nil
    }
    return w.list
}

// SetRows syncs rows with gotui List
func (w *BaseListWidget) SetRows(rows []string) error {
    w.rows = rows
    if w.list != nil {
        w.list.Rows = rows  // Direct gotui property
    }
    return nil
}

// Access selected index from gotui List
func (w *BaseListWidget) GetSelectedIndex() int {
    if w == nil || w.list == nil {
        return 0
    }
    return w.list.SelectedRow
}
```

## DetailsWidget for Group Information

```go
type DetailsWidget struct {
    paragraph *widgets.Paragraph
    analyzer  *browsing.PresenceAnalyzer
}

func (w *DetailsWidget) ShowGroup(group *DuplicateGroup) error {
    detail := w.formatGroupDetails(group)
    w.paragraph.Text = detail
    return nil
}

// Output format:
// Identifier: abc123def456...
// Total Copies: 5
// Flag Status: [KEEP]
//
// VERSIONS (3):
//   1. 🟢 file1.txt (Hash: abc123de)
//      → remote1 (remote)
//      → local_drive (local)
```

## Hash Rendering

```go
func (w *HashCentricWidget) renderHashCentricList() {
    for idx, group := range w.duplicates {
        indicator := w.analyzer.AnalyzeDuplicateGroup(group)

        flagStr := ""
        if group.Flag != nil {
            flagStr = fmt.Sprintf(" [%s]", group.Flag.FlagType)
        }

        // Format: "🔴 abc12345... (5 copies) [KEEP]"
        row := fmt.Sprintf("%s %s... (%d copies)%s",
            indicator,
            w.truncateHash(group.Identifier, 8),
            group.Count,
            flagStr,
        )
        rows = append(rows, row)
    }
    w.list.SetRows(rows)
}
```

## Filename Rendering

```go
func (w *FilenameCentricWidget) renderFilenameCentricList() {
    for idx, group := range w.duplicates {
        indicator := w.analyzer.AnalyzeDuplicateGroup(group)

        // Format: "🟢 document.pdf (3 versions) [ARCHIVE]"
        row := fmt.Sprintf("%s %s (%d %s)%s",
            indicator,
            group.Identifier,
            len(group.Versions),
            versionWord,  // "version" or "versions"
            flagStr,
        )
        rows = append(rows, row)
    }
}
```

## Nil Safety Pattern

```go
func (w *BaseListWidget) ScrollDown() {
    if w == nil || w.list == nil {
        return
    }
    if w.list.SelectedRow < len(w.rows)-1 {
        w.list.SelectedRow++
    }
}
```
