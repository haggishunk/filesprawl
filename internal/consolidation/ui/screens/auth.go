package screens

import (
	"context"
	"fmt"
	"strings"

	"github.com/haggishunk/filesprawl/internal/consolidation/users"
)

// AuthScreen handles user authentication (login or register)
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

// Show displays the authentication screen and handles user auth flow
func (as *AuthScreen) Show(ctx context.Context) (*users.User, error) {
	if as.userRepo == nil {
		return nil, fmt.Errorf("user repository not initialized")
	}

	// TODO: Phase 4 - Implement full gotui modal authentication
	// For now, use console-based input as placeholder
	// This allows the system to work without full gotui modal integration yet

	// Create a simple console-based authentication input
	fmt.Print("Enter your username: ")
	var username string
	_, err := fmt.Scanln(&username)
	if err != nil {
		return nil, fmt.Errorf("failed to read username: %w", err)
	}

	username = strings.TrimSpace(username)

	// Validate username
	if err := as.validateUsername(username); err != nil {
		return nil, fmt.Errorf("username validation failed: %w", err)
	}

	// Get or create user
	user, err := as.getUserOrCreate(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	as.currentUser = user
	return user, nil
}

// validateUsername validates the username format and constraints
func (as *AuthScreen) validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(username) > 32 {
		return fmt.Errorf("username must be 32 characters or less (got %d)", len(username))
	}

	return nil
}

// getUserOrCreate retrieves an existing user or creates a new one
func (as *AuthScreen) getUserOrCreate(ctx context.Context, username string) (*users.User, error) {
	if as.userRepo == nil {
		return nil, fmt.Errorf("user repository not initialized")
	}

	// Try to find existing user
	user, err := as.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		// User found, update last_active
		if updateErr := as.userRepo.UpdateLastActive(ctx, user.ID); updateErr != nil {
			// Log but don't fail - continue with existing user
			fmt.Printf("warning: failed to update last_active: %v\n", updateErr)
		}
		return user, nil
	}

	// User not found, create new one
	newUser, err := as.userRepo.CreateUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	return newUser, nil
}

// GetCurrentUser returns the authenticated user
func (as *AuthScreen) GetCurrentUser() *users.User {
	return as.currentUser
}

// SetCurrentUser sets the current user
func (as *AuthScreen) SetCurrentUser(user *users.User) {
	as.currentUser = user
}
