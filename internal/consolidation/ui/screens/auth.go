package screens

import (
	"context"
	"fmt"
	"github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/haggishunk/filesprawl/internal/consolidation/users"
)

// AuthScreen handles user authentication (login or register)
type AuthScreen struct {
	app          *gotui.Application
	userRepo     users.Repository
	currentUser  *users.User
	selectedMode string // "login" or "register"
	selectedUser *users.User
	allUsers     []users.User
	error        string
}

// NewAuthScreen creates a new authentication screen
func NewAuthScreen(app *gotui.Application, userRepo users.Repository) *AuthScreen {
	return &AuthScreen{
		app:      app,
		userRepo: userRepo,
		error:    "",
	}
}

// Show displays the authentication screen
func (as *AuthScreen) Show(ctx context.Context) (*users.User, error) {
	// Load existing users
	existingUsers, err := as.userRepo.GetFlagsByUser(ctx, 0) // Dummy call to check connection
	_ = err // Ignore for now

	// Create title
	title := widgets.NewParagraph()
	title.Title = "🔐 CONSOLIDATION TOOL - USER AUTHENTICATION"
	title.Text = "Select an action below:"
	title.SetRect(0, 0, gotui.MaxX(), 5)
	title.TitleAlignment = gotui.AlignCenter
	title.BorderStyle = gotui.NewStyle().Foreground(gotui.ColorCyan)

	// Create mode selection list
	modeList := widgets.NewList()
	modeList.Title = "SELECT MODE"
	modeList.Rows = []string{
		"[L] Login as Existing User",
		"[R] Register New User",
	}
	modeList.SelectedRow = 0
	modeList.SetRect(5, 5, 60, 10)

	// Create user selection list
	userList := widgets.NewList()
	userList.Title = "SELECT USER"
	userList.SetRect(5, 16, 60, 25)

	// Create input field for new user
	input := widgets.NewInput()
	input.Label = "Username:"
	input.SetRect(5, 16, 60, 19)

	// Create info paragraph
	info := widgets.NewParagraph()
	info.Title = "INSTRUCTIONS"
	info.Text = "[L] Login / [R] Register\n[Enter] Select / [Q] Quit"
	info.SetRect(5, 26, 60, 30)
	info.BorderStyle = gotui.NewStyle().Foreground(gotui.ColorGreen)

	// Create error paragraph
	errPara := widgets.NewParagraph()
	errPara.SetRect(0, gotui.MaxY()-3, gotui.MaxX(), gotui.MaxY())
	errPara.BorderStyle = gotui.NewStyle().Foreground(gotui.ColorRed)

	// Grid layout
	grid := gotui.NewGrid()
	grid.Set(
		gotui.NewRow(5, title),
		gotui.NewRow(10, modeList),
		gotui.NewRow(9, userList, input),
		gotui.NewRow(5, info),
		gotui.NewRow(3, errPara),
	)

	// Event handling
	uiEvents := gotui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return nil, fmt.Errorf("user cancelled authentication")
		case "l", "L":
			as.selectedMode = "login"
			// Show user list
			userList.Height = 9
			input.Height = 0
		case "r", "R":
			as.selectedMode = "register"
			// Show input field
			userList.Height = 0
			input.Height = 3
		case "<Enter>":
			if as.selectedMode == "login" && userList.Height > 0 {
				if userList.SelectedRow < len(as.allUsers) {
					user := as.allUsers[userList.SelectedRow]
					return &user, nil
				}
			} else if as.selectedMode == "register" && input.Height > 0 {
				username := input.GetText()
				if username == "" {
					as.error = "Username cannot be empty"
					continue
				}
				// Create new user
				newUser, err := as.userRepo.CreateUser(ctx, username)
				if err != nil {
					as.error = fmt.Sprintf("Error creating user: %v", err)
					continue
				}
				return newUser, nil
			}
		}
	}
}

// GetCurrentUser returns the authenticated user
func (as *AuthScreen) GetCurrentUser() *users.User {
	return as.currentUser
}

// SetCurrentUser sets the current user
func (as *AuthScreen) SetCurrentUser(user *users.User) {
	as.currentUser = user
}
