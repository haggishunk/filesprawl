package screens

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/users"
)

// AuthScreen handles user authentication (login or register)
// NOTE: This is a stub implementation. Full gotui integration is in progress.
type AuthScreen struct {
	userRepo    users.Repository
	currentUser *users.User
}

// NewAuthScreen creates a new authentication screen
func NewAuthScreen(userRepo users.Repository) *AuthScreen {
	return &AuthScreen{
		userRepo: userRepo,
	}
}

// Show displays the authentication screen (stub implementation)
func (as *AuthScreen) Show(ctx context.Context) (*users.User, error) {
	// TODO: Implement full gotui-based authentication screen
	// For now, return a stub user for testing
	user, err := as.userRepo.CreateUser(ctx, "default-user")
	if err != nil {
		return nil, fmt.Errorf("failed to get/create default user: %w", err)
	}
	as.currentUser = user
	return user, nil
}

// GetCurrentUser returns the authenticated user
func (as *AuthScreen) GetCurrentUser() *users.User {
	return as.currentUser
}

// SetCurrentUser sets the current user
func (as *AuthScreen) SetCurrentUser(user *users.User) {
	as.currentUser = user
}
