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
	selectedIndex int    // Currently selected duplicate group
	isHashView    bool   // true = hash view, false = filename view
	modalType     string // Current modal being shown ("" = none)
	focusedPane   string // "list" or "table" - which pane has focus
	width         int    // Terminal width
	height        int    // Terminal height

	// Settings
	version    string
	reviewMode bool
	remotes    []string

	// Data
	duplicateGroups []browsing.DuplicateGroup

	// Components
	listPane      *components.ListPane
	tablePane     *components.TablePane
	banner        *components.Banner
	statusBar     *components.StatusBar
	safetyChecker *components.SafetyChecker

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
		ctx:             ctx,
		sessionManager:  sessionManager,
		flagManager:     flagManager,
		duplicateGroups: duplicateGroups,
		selectedIndex:   0,
		isHashView:      true,
		focusedPane:     "list", // Start with list pane focused
		modalType:       "",
		width:           80,
		height:          24,
		version:         version,
		reviewMode:      reviewMode,
		remotes:         remotes,
	}

	// Initialize components
	model.listPane = components.NewListPane(duplicateGroups, 40, 20)
	model.tablePane = components.NewTablePane(60, 20, flagManager)
	model.banner = components.NewBanner(version, len(duplicateGroups), reviewMode, remotes, 80)
	model.statusBar = components.NewStatusBar(80)
	model.safetyChecker = components.NewSafetyChecker(flagManager)

	// Set initial group in table
	if len(duplicateGroups) > 0 {
		model.tablePane.SetGroup(&duplicateGroups[0])
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
		// Update component dimensions (40% left, 60% right)
		leftWidth := (m.width * 40) / 100
		rightWidth := m.width - leftWidth - 2 // Account for borders
		panelHeight := m.height - 4
		if m.listPane != nil {
			m.listPane.SetDimensions(leftWidth, panelHeight)
		}
		if m.tablePane != nil {
			m.tablePane.SetDimensions(rightWidth, panelHeight)
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

	// Render banner
	if m.banner != nil {
		result.WriteString(m.banner.Render())
	}

	// Update focus indicators on panes
	if m.listPane != nil {
		m.listPane.SetFocus(m.focusedPane == "list")
	}
	if m.tablePane != nil {
		m.tablePane.SetFocus(m.focusedPane == "table")
	}

	// Render panes side-by-side
	var listView, tableView string
	if m.listPane != nil {
		listView = m.listPane.Render()
	}
	if m.tablePane != nil {
		tableView = m.tablePane.Render()
	}

	// Combine panes horizontally
	panelsView := m.combinePanels(listView, tableView)
	result.WriteString(panelsView)

	// Render error message if present
	if m.lastError != "" {
		fmt.Fprintf(&result, "⚠️ ERROR: %s\n", m.lastError)
	}

	// Render status bar
	if m.statusBar != nil {
		result.WriteString(m.statusBar.Render())
	}

	return result.String()
}

// combinePanels combines two pane views side-by-side
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
		// Switch focus between panes
		if m.focusedPane == "list" {
			m.focusedPane = "table"
		} else {
			m.focusedPane = "list"
		}
	case "shift+tab":
		// Switch focus between panes (reverse)
		if m.focusedPane == "list" {
			m.focusedPane = "table"
		} else {
			m.focusedPane = "list"
		}
	case "j":
		// Navigate down in focused pane
		if m.focusedPane == "list" {
			if m.listPane != nil {
				m.listPane.SelectionDown()
				m.selectedIndex = m.listPane.GetSelectedIndex()
				m.updateTablePane()
			}
		} else {
			// Navigate in table
			if m.tablePane != nil {
				m.tablePane.SelectionDown()
			}
		}
	case "k":
		// Navigate up in focused pane
		if m.focusedPane == "list" {
			if m.listPane != nil {
				m.listPane.SelectionUp()
				m.selectedIndex = m.listPane.GetSelectedIndex()
				m.updateTablePane()
			}
		} else {
			// Navigate in table
			if m.tablePane != nil {
				m.tablePane.SelectionUp()
			}
		}
	case "v":
		// Toggle wide view in table pane (only works when table is focused)
		if m.focusedPane == "table" && m.tablePane != nil {
			m.tablePane.ToggleWideView()
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

// updateTablePane updates the table pane with the currently selected group
func (m *ConsolidationModel) updateTablePane() {
	if m.listPane != nil && m.tablePane != nil {
		selectedGroup := m.listPane.GetSelectedGroup()
		if selectedGroup != nil {
			m.tablePane.SetGroup(selectedGroup)
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
	if m.tablePane == nil || m.flagManager == nil {
		return
	}

	// Get the currently selected file
	version := m.tablePane.GetSelectedVersion()
	if version == nil {
		return
	}

	switch action {
	case "keep":
		if err := m.flagManager.PersistKeep(m.ctx, version.FileID); err != nil {
			fmt.Printf("Error flagging as KEEP: %v\n", err)
		}
	case "archive":
		// TODO: Show modal to get archive destination
		destination := "archive" // Placeholder
		reason := "User flagged via TUI"
		if err := m.flagManager.PersistArchive(m.ctx, version.FileID, destination, &reason); err != nil {
			fmt.Printf("Error flagging as ARCHIVE: %v\n", err)
		}
	case "delete":
		// Check deletion safety
		selectedGroup := m.listPane.GetSelectedGroup()
		if selectedGroup == nil {
			m.lastError = "No group selected"
			return
		}

		safe, reason := m.safetyChecker.CheckDeletionSafety(version, selectedGroup)
		if !safe {
			m.lastError = reason
			return
		}

		// Safe to delete
		deleteReason := "User flagged via TUI"
		if err := m.flagManager.PersistDelete(m.ctx, version.FileID, &deleteReason); err != nil {
			fmt.Printf("Error flagging as DELETE: %v\n", err)
		}
		m.lastError = "" // Clear error on success
	case "skip":
		if err := m.flagManager.PersistSkip(m.ctx, version.FileID); err != nil {
			fmt.Printf("Error clearing flag: %v\n", err)
		}
	}
}

// showHelp displays help
func (m *ConsolidationModel) showHelp() {
	m.modalType = "help"
}
