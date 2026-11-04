package ui

import (
	"context"
	"fmt"

	"github.com/haggishunk/filesprawl/internal/consolidation/browsing"
	"github.com/haggishunk/filesprawl/internal/consolidation/ui/widgets"
)

// KeyEventHandler is a function type for handling keyboard input
type KeyEventHandler func(ctx context.Context, widget widgets.BrowsingWidget, session *browsing.SessionManager) error

// KeyBindingRegistry maps key runes to event handlers
type KeyBindingRegistry struct {
	handlers map[rune]KeyEventHandler
}

// NewKeyBindingRegistry creates a registry with default keybindings
func NewKeyBindingRegistry() *KeyBindingRegistry {
	return &KeyBindingRegistry{
		handlers: make(map[rune]KeyEventHandler),
	}
}

// Register binds a key to a handler
func (kr *KeyBindingRegistry) Register(key rune, handler KeyEventHandler) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil for key %c", key)
	}
	kr.handlers[key] = handler
	return nil
}

// Handle processes a key event by looking up the registered handler
func (kr *KeyBindingRegistry) Handle(key rune, ctx context.Context, widget widgets.BrowsingWidget, session *browsing.SessionManager) error {
	if kr == nil {
		return fmt.Errorf("key registry not initialized")
	}

	if handler, ok := kr.handlers[key]; ok {
		return handler(ctx, widget, session)
	}

	// Unknown key: silently ignore
	return nil
}

// GetRegisteredKeys returns list of all bound keys
func (kr *KeyBindingRegistry) GetRegisteredKeys() []rune {
	if kr == nil {
		return []rune{}
	}
	keys := make([]rune, 0, len(kr.handlers))
	for k := range kr.handlers {
		keys = append(keys, k)
	}
	return keys
}

// HasBinding checks if a key has a registered handler
func (kr *KeyBindingRegistry) HasBinding(key rune) bool {
	if kr == nil {
		return false
	}
	_, ok := kr.handlers[key]
	return ok
}
