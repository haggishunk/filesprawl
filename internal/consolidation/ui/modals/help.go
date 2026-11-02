package modals

import "fmt"

// HelpContent returns formatted help text
func HelpContent() string {
	return fmt.Sprintf(`
╔════════════════════════════════════════════════════════════════╗
║                    FILESPRAWL HELP                            ║
╚════════════════════════════════════════════════════════════════╝

📋 KEYBOARD SHORTCUTS:

Navigation:
  j         Scroll down to next duplicate group
  k         Scroll up to previous duplicate group
  v         Switch between Hash view and Filename view

Actions:
  K         Mark entire group as KEEP (all copies)
  a         Mark group for ARCHIVE (backup to remote)
  d         Mark group for DELETE (remove duplicates)
  s         SKIP group (clear any flags)

Help & Exit:
  ?         Show this help screen
  q         Quit and review summary

═══════════════════════════════════════════════════════════════════

📊 DISPLAY GUIDE:

Colors & Indicators:
  🟢        Safe - Multiple backup locations exist
  🟡        Caution - Limited replication, backup recommended
  🔴        Critical - Single location, high risk
  🚩        Flagged - Decision already made

List View Shows:
  1. Safety indicator (🔴🟡🟢)
  2. Hash or filename identifier
  3. Number of copies
  4. Number of storage locations
  5. Any applied flags

Details Pane Shows:
  1. Each individual copy with safety indicator
  2. File path and size (human-readable)
  3. Storage locations for each copy
  4. Location summary table
  5. Action shortcuts

═══════════════════════════════════════════════════════════════════

🎯 WORKFLOW TIPS:

1. Review list: Scroll through duplicates with j/k
2. Check details: Automatically shown when selected
3. Assess risk: Look at location indicators (🔴🟡🟢)
4. Make decision: Use K/a/d/s to flag the group
5. Switch views: Use v to see same data differently
6. Review summary: Press q to exit and see decisions

═══════════════════════════════════════════════════════════════════

❓ UNDERSTANDING THE DATA:

What are duplicates?
  Files with identical content (same hash) across
  different locations in your storage systems.

Why does location matter?
  🟢 Multiple copies = safe to delete some
  🟡 Few copies = backup before deleting
  🔴 Single copy = CRITICAL - very risky to delete

What do the modes mean?
  REVIEW MODE: See changes without executing
  EXECUTE MODE: Changes are applied immediately

═══════════════════════════════════════════════════════════════════

💡 QUICK REFERENCE:

To mark a group (after selecting with j/k):
  [K] → Keep all copies (no changes)
  [a] → Archive extras to remote backup
  [d] → Delete all but one copy
  [s] → Skip for now (undo any flags)

Press SPACE or ENTER to continue...
`)
}

// ShowHelpModal displays help and returns when closed
func ShowHelpModal() error {
	content := HelpContent()
	fmt.Print(content)
	return nil
}
