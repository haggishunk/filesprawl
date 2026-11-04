package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/components"
)

// ConsolidationModel implements tea.Model for the consolidation TUI
type ConsolidationModel struct {
	// Context
	ctx context.Context

	// Business logic managers
	sessionManager *browsing.SessionManager
	flagManager    *browsing.PersistentFlagManager

	// UI state
	selectedIndex       int    // Currently selected duplicate group
	isHashView          bool   // true = hash view, false = filename view
	modalType           string // Current modal being shown ("" = none)
	focusedPane         string // "list", "hash_table", or "filename_table" - which pane has focus
	width               int    // Terminal width
	height              int    // Terminal height
	showFilenamePaneVal bool   // Whether to show the lower right filename pane

	// Settings
	version    string
	reviewMode bool
	remotes    []string

	// Data
	duplicateGroups []browsing.DuplicateGroup

	// Components
	listPane          *components.ListPane
	hashTablePane     *components.TablePane         // Upper right: hash-grouped files
	filenameTablePane *components.FilenameTablePane // Lower right: filename-grouped files
	banner            *components.Banner
	statusBar         *components.StatusBar
	safetyChecker     *components.SafetyChecker
	helpModal         *components.HelpModal
	filenameGroups    []browsing.DuplicateGroup // Filename groups for display

	// UI feedback
	lastError string // Last error message to display
}

// NewConsolidationModel creates a new Bubble Tea model
func NewConsolidationModel(
	ctx context.Context,
	sessionManager *browsing.SessionManager,
	flagManager *browsing.PersistentFlagManager,
	duplicateGroups []browsing.DuplicateGroup,
	version string,
	reviewMode bool,
	remotes []string,
) *ConsolidationModel {
	model := &ConsolidationModel{
		ctx:                 ctx,
		sessionManager:      sessionManager,
		flagManager:         flagManager,
		duplicateGroups:     duplicateGroups,
		selectedIndex:       0,
		isHashView:          true,
		focusedPane:         "list", // Start with list pane focused
		modalType:           "",
		width:               80,
		height:              24,
		version:             version,
		reviewMode:          reviewMode,
		remotes:             remotes,
		showFilenamePaneVal: true,                        // Show filename pane by default
		filenameGroups:      []browsing.DuplicateGroup{}, // Will be loaded later
	}

	// Initialize components
	model.listPane = components.NewListPane(duplicateGroups, 40, 20)
	model.hashTablePane = components.NewTablePane(60, 20, flagManager)
	model.filenameTablePane = components.NewFilenameTablePane(60, 20, flagManager)
	model.banner = components.NewBanner(version, len(duplicateGroups), reviewMode, remotes, 80)
	model.statusBar = components.NewStatusBar(80)
	model.safetyChecker = components.NewSafetyChecker(flagManager)
	model.helpModal = components.NewHelpModal(80, 24)

	// Set initial group in hash table pane
	if len(duplicateGroups) > 0 {
		model.hashTablePane.SetGroup(&duplicateGroups[0])
	}

	return model
}

// Init initializes the model
func (m *ConsolidationModel) Init() tea.Cmd {
	// Load initial data
	return nil
}

// Update handles messages
func (m *ConsolidationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update component dimensions (40% left, 60% right split vertically)
		leftWidth := (m.width * 40) / 100
		rightWidth := m.width - leftWidth - 2 // Account for borders
		panelHeight := m.height - 4

		if m.listPane != nil {
			m.listPane.SetDimensions(leftWidth, panelHeight)
		}

		// Right side is split vertically: hash pane on top, filename pane on bottom
		if m.showFilenamePaneVal && panelHeight > 6 {
			// Split vertically
			hashPaneH := panelHeight / 2
			filenamePaneH := panelHeight - hashPaneH - 1 // -1 for divider
			if m.hashTablePane != nil {
				m.hashTablePane.SetDimensions(rightWidth, hashPaneH)
			}
			if m.filenameTablePane != nil {
				m.filenameTablePane.SetDimensions(rightWidth, filenamePaneH)
			}
		} else {
			// Only show hash pane, full height
			if m.hashTablePane != nil {
				m.hashTablePane.SetDimensions(rightWidth, panelHeight)
			}
		}

		if m.banner != nil {
			m.banner.SetDimensions(m.width)
		}
		if m.statusBar != nil {
			m.statusBar.SetDimensions(m.width)
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case ScrollMsg:
		m.handleScroll(msg.Direction)
		return m, nil

	case ViewSwitchMsg:
		m.isHashView = !m.isHashView
		return m, nil

	case ActionMsg:
		m.handleAction(msg.Action)
		return m, nil

	case ModalDismissMsg:
		m.modalType = ""
		return m, nil

	case QuitMsg:
		return m, tea.Quit

	default:
		return m, nil
	}
}

// View renders the model
func (m *ConsolidationModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Waiting for terminal...\n"
	}

	var result strings.Builder

	// If help modal is open, show it instead of normal view
	if m.modalType == "help" {
		if m.helpModal != nil {
			m.helpModal.Width = m.width
			m.helpModal.Height = m.height
			return m.helpModal.Render()
		}
	}

	// Render banner
	if m.banner != nil {
		result.WriteString(m.banner.Render())
	}

	// Update focus indicators on panes
	if m.listPane != nil {
		m.listPane.SetFocus(m.focusedPane == "list")
	}
	if m.hashTablePane != nil {
		m.hashTablePane.SetFocus(m.focusedPane == "hash_table")
	}
	if m.filenameTablePane != nil {
		m.filenameTablePane.SetFocus(m.focusedPane == "filename_table")
	}

	// Render list pane
	var listView string
	if m.listPane != nil {
		listView = m.listPane.Render()
	}

	// Render right side panes (vertically split or single)
	var rightPanelsView string
	if m.showFilenamePaneVal {
		// Both hash and filename panes visible - render vertically split
		var hashTableView, filenameTableView string
		if m.hashTablePane != nil {
			hashTableView = m.hashTablePane.Render()
		}
		if m.filenameTablePane != nil {
			filenameTableView = m.filenameTablePane.Render()
		}
		rightPanelsView = m.combineRightPanels(hashTableView, filenameTableView)
	} else {
		// Only hash pane visible
		if m.hashTablePane != nil {
			rightPanelsView = m.hashTablePane.Render()
		}
	}

	// Combine left and right panes horizontally
	panelsView := m.combinePanels(listView, rightPanelsView)
	result.WriteString(panelsView)

	// Render error message if present
	if m.lastError != "" {
		fmt.Fprintf(&result, "⚠️ ERROR: %s\n", m.lastError)
	}

	// Render status bar
	if m.statusBar != nil {
		m.statusBar.SetFocusedPane(m.focusedPane)
		result.WriteString(m.statusBar.Render())
	}

	return result.String()
}

// combinePanels combines two pane views side-by-side (left + right)
func (m *ConsolidationModel) combinePanels(leftPane, rightPane string) string {
	leftLines := strings.Split(leftPane, "\n")
	rightLines := strings.Split(rightPane, "\n")

	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	var result strings.Builder
	for i := 0; i < maxLines; i++ {
		leftLine := ""
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		rightLine := ""
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		result.WriteString(leftLine)
		result.WriteString(rightLine)
		result.WriteString("\n")
	}

	return result.String()
}

// combineRightPanels combines hash and filename panes vertically (top + bottom)
func (m *ConsolidationModel) combineRightPanels(hashPane, filenamePane string) string {
	// Simply stack vertically with a divider line between them
	// The panes themselves handle their rendering, we just concatenate
	var result strings.Builder
	result.WriteString(hashPane)
	// Add a divider line - simple horizontal line
	result.WriteString("─────────────────────────────────────────────────────────────────\n")
	result.WriteString(filenamePane)
	return result.String()
}

// handleKey processes keyboard input
func (m *ConsolidationModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If help modal is open, close it on any key (except q for quit)
	if m.modalType == "help" {
		if msg.String() != "q" {
			m.modalType = ""
			return m, nil
		}
		// q will fall through and quit below
	}

	switch msg.String() {
	case "tab":
		// Switch focus between panes (cycle through: list -> hash_table -> filename_table -> list)
		switch m.focusedPane {
		case "list":
			m.focusedPane = "hash_table"
		case "hash_table":
			if m.showFilenamePaneVal {
				m.focusedPane = "filename_table"
			} else {
				m.focusedPane = "list"
			}
		case "filename_table":
			m.focusedPane = "list"
		default:
			m.focusedPane = "list"
		}
	case "shift+tab":
		// Switch focus between panes (reverse cycle)
		switch m.focusedPane {
		case "list":
			if m.showFilenamePaneVal {
				m.focusedPane = "filename_table"
			} else {
				m.focusedPane = "hash_table"
			}
		case "hash_table":
			m.focusedPane = "list"
		case "filename_table":
			m.focusedPane = "hash_table"
		default:
			m.focusedPane = "list"
		}
	case "j":
		// Navigate down in focused pane
		switch m.focusedPane {
		case "list":
			if m.listPane != nil {
				m.listPane.SelectionDown()
				m.selectedIndex = m.listPane.GetSelectedIndex()
				m.updateHashTablePane()
			}
		case "hash_table":
			if m.hashTablePane != nil {
				m.hashTablePane.SelectionDown()
			}
		case "filename_table":
			if m.filenameTablePane != nil {
				m.filenameTablePane.SelectionDown()
			}
		}
	case "k":
		// Navigate up in focused pane
		switch m.focusedPane {
		case "list":
			if m.listPane != nil {
				m.listPane.SelectionUp()
				m.selectedIndex = m.listPane.GetSelectedIndex()
				m.updateHashTablePane()
			}
		case "hash_table":
			if m.hashTablePane != nil {
				m.hashTablePane.SelectionUp()
			}
		case "filename_table":
			if m.filenameTablePane != nil {
				m.filenameTablePane.SelectionUp()
			}
		}
	case "w":
		// Toggle filename pane visibility
		m.ToggleFilenamePaneVisibility()
		// Adjust focus if needed: if filename pane is hidden and focused, move to hash_table
		if !m.showFilenamePaneVal && m.focusedPane == "filename_table" {
			m.focusedPane = "hash_table"
		}
	case "v":
		// Toggle wide view in focused table pane
		if m.focusedPane == "hash_table" && m.hashTablePane != nil {
			m.hashTablePane.ToggleWideView()
		} else if m.focusedPane == "filename_table" && m.filenameTablePane != nil {
			m.filenameTablePane.ToggleWideView()
		}
	case "p":
		// Open file for preview (only in table panes)
		var selectedFile *browsing.FileVersion
		if m.focusedPane == "hash_table" && m.hashTablePane != nil {
			selectedFile = m.hashTablePane.GetSelectedFile()
		} else if m.focusedPane == "filename_table" && m.filenameTablePane != nil {
			selectedFile = m.filenameTablePane.GetSelectedFile()
		}
		if selectedFile != nil && selectedFile.FileInfo != nil {
			// Construct full filepath for preview
			fullPath := m.constructFilePath(selectedFile)
			if fullPath != "" {
				// TODO: Open file with system viewer (e.g., less, vim, etc.)
				m.lastError = fmt.Sprintf("📁 Preview: %s", fullPath)
			} else {
				m.lastError = "❌ Cannot construct file path"
			}
		}
	case "K":
		// Flag as KEEP (only in table panes)
		if m.focusedPane == "hash_table" || m.focusedPane == "filename_table" {
			m.handleAction("keep")
		}
	case "a":
		// Flag as ARCHIVE (only in table panes)
		if m.focusedPane == "hash_table" || m.focusedPane == "filename_table" {
			m.handleAction("archive")
		}
	case "d":
		// Flag as DELETE (only in table panes)
		if m.focusedPane == "hash_table" || m.focusedPane == "filename_table" {
			m.handleAction("delete")
		}
	case "s":
		// Flag as SKIP (only in table panes)
		if m.focusedPane == "hash_table" || m.focusedPane == "filename_table" {
			m.handleAction("skip")
		}
	case "?":
		m.showHelp()
	case "q":
		return m, tea.Quit
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

// updateHashTablePane updates the hash table pane with the currently selected group
func (m *ConsolidationModel) updateHashTablePane() {
	if m.listPane != nil && m.hashTablePane != nil {
		selectedGroup := m.listPane.GetSelectedGroup()
		if selectedGroup != nil {
			m.hashTablePane.SetGroup(selectedGroup)
		}
	}
}

// handleScroll moves selection up/down
func (m *ConsolidationModel) handleScroll(direction string) {
	if m.listPane != nil {
		switch direction {
		case "down":
			m.listPane.SelectionDown()
		case "up":
			m.listPane.SelectionUp()
		}
		m.selectedIndex = m.listPane.GetSelectedIndex()
	}
}

// handleAction processes user actions (K/a/d/s flags)
func (m *ConsolidationModel) handleAction(action string) {
	if m.flagManager == nil {
		return
	}

	// Get the currently selected file from the focused pane
	var version *browsing.FileVersion
	if m.focusedPane == "hash_table" && m.hashTablePane != nil {
		version = m.hashTablePane.GetSelectedFile()
	} else if m.focusedPane == "filename_table" && m.filenameTablePane != nil {
		version = m.filenameTablePane.GetSelectedFile()
	}

	if version == nil {
		m.lastError = "No file selected"
		return
	}

	switch action {
	case "keep":
		m.lastError = "" // Clear previous error
		if err := m.flagManager.PersistKeep(m.ctx, version.FileID); err != nil {
			m.lastError = fmt.Sprintf("❌ Keep failed: %v", err)
		} else {
			m.lastError = "✓ Flagged as KEEP"
		}
	case "archive":
		// TODO: Show modal to get archive destination
		m.lastError = ""         // Clear previous error
		destination := "archive" // Placeholder
		reason := "User flagged via TUI"
		if err := m.flagManager.PersistArchive(m.ctx, version.FileID, destination, &reason); err != nil {
			m.lastError = fmt.Sprintf("❌ Archive failed: %v", err)
		} else {
			m.lastError = "✓ Flagged as ARCHIVE"
		}
	case "delete":
		m.lastError = "" // Clear previous error
		// Check deletion safety
		selectedGroup := m.listPane.GetSelectedGroup()
		if selectedGroup == nil {
			m.lastError = "❌ No group selected"
			return
		}

		safe, reason := m.safetyChecker.CheckDeletionSafety(version, selectedGroup)
		if !safe {
			m.lastError = fmt.Sprintf("❌ %s", reason)
			return
		}

		// Safe to delete
		deleteReason := "User flagged via TUI"
		if err := m.flagManager.PersistDelete(m.ctx, version.FileID, &deleteReason); err != nil {
			m.lastError = fmt.Sprintf("❌ Delete failed: %v", err)
		} else {
			m.lastError = "✓ Flagged as DELETE"
		}
	case "skip":
		m.lastError = "" // Clear previous error
		if err := m.flagManager.PersistSkip(m.ctx, version.FileID); err != nil {
			m.lastError = fmt.Sprintf("❌ Skip failed: %v", err)
		} else {
			m.lastError = "✓ Cleared flags"
		}
	}
}

// constructFilePath builds the full filepath for opening/previewing a file.
//
// LOCAL FILES (ScanRoot set):
//
//	Display: hostname | filename | /full/canonical/path
//	Accessor: /full/canonical/path (pass directly to system viewer)
//
// REMOTE FILES (ScanRoot null):
//
// constructFilePath returns the absolute path for a file.
// For local files: full filesystem path (e.g., "/mnt/data/photos/IMG_0587.MOV")
// For rclone files: remote:path format (e.g., "goblet:photos/IMG_0587.MOV")
// The absolute path is stored directly in Meta.AbsolutePath.
func (m *ConsolidationModel) constructFilePath(version *browsing.FileVersion) string {
	if version == nil || version.FileInfo == nil {
		return ""
	}

	// AbsolutePath contains the full, ready-to-use path
	if version.FileInfo.AbsolutePath != "" {
		return version.FileInfo.AbsolutePath
	}

	// Fallback if AbsolutePath not set
	return version.FileInfo.Path
}

// SetFilenameGroups loads filename-based duplicate groups for display
func (m *ConsolidationModel) SetFilenameGroups(groups []browsing.DuplicateGroup) {
	m.filenameGroups = groups
	if len(groups) > 0 && m.filenameTablePane != nil {
		m.filenameTablePane.SetGroup(&groups[0])
	}
}

// ToggleFilenamePaneVisibility toggles whether the filename pane is shown
func (m *ConsolidationModel) ToggleFilenamePaneVisibility() {
	m.showFilenamePaneVal = !m.showFilenamePaneVal
	// Focus cycling updates handled in navigation logic
}

// IsFilenamePaneVisible returns whether the filename pane should be displayed
func (m *ConsolidationModel) IsFilenamePaneVisible() bool {
	return m.showFilenamePaneVal
}

// showHelp displays help
func (m *ConsolidationModel) showHelp() {
	m.modalType = "help"
}
