package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
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
	width         int    // Terminal width
	height        int    // Terminal height

	// Settings
	version    string
	reviewMode bool
	remotes    []string

	// Data
	duplicateGroups []browsing.DuplicateGroup
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
	return &ConsolidationModel{
		ctx:             ctx,
		sessionManager:  sessionManager,
		flagManager:     flagManager,
		duplicateGroups: duplicateGroups,
		selectedIndex:   0,
		isHashView:      true,
		modalType:       "",
		width:           80,
		height:          24,
		version:         version,
		reviewMode:      reviewMode,
		remotes:         remotes,
	}
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

	// For now, return a placeholder that shows the model is working
	view := fmt.Sprintf("📎 Filesprawl v%s | Groups: %d | Mode: ",
		m.version, len(m.duplicateGroups))
	if m.reviewMode {
		view += "REVIEW"
	} else {
		view += "EXECUTE"
	}
	view += " | Remotes: " + fmt.Sprint(m.remotes) + "\n"
	view += "Keyboard: [j]down [k]up [v]view [K]keep [a]archive [d]delete [s]skip [?]help [q]quit\n"

	return view
}

// handleKey processes keyboard input
func (m *ConsolidationModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j":
		m.handleScroll("down")
	case "k":
		m.handleScroll("up")
	case "v":
		m.isHashView = !m.isHashView
	case "K":
		m.handleAction("keep")
	case "a":
		m.handleAction("archive")
	case "d":
		m.handleAction("delete")
	case "s":
		m.handleAction("skip")
	case "?":
		m.showHelp()
	case "q":
		return m, tea.Quit
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

// handleScroll moves selection up/down
func (m *ConsolidationModel) handleScroll(direction string) {
	if direction == "down" && m.selectedIndex < len(m.duplicateGroups)-1 {
		m.selectedIndex++
	} else if direction == "up" && m.selectedIndex > 0 {
		m.selectedIndex--
	}
}

// handleAction processes user actions
func (m *ConsolidationModel) handleAction(action string) {
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.duplicateGroups) {
		return
	}
	// TODO: Implement action handlers
	fmt.Printf("Action: %s on index %d\n", action, m.selectedIndex)
}

// showHelp displays help
func (m *ConsolidationModel) showHelp() {
	m.modalType = "help"
}
