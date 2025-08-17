// components/registry.go
package components

import (
	"context"

	"github.com/a-h/templ"
)

// Handler interfaces for the registry
type AppHandler interface {
	RenderComponent(ctx context.Context) templ.Component
}

type BoardHandler interface {
	RenderComponent(ctx context.Context) templ.Component
}

// Registry holds references to all component handlers
type Registry struct {
	AppHandler   AppHandler
	BoardHandler BoardHandler
}

// Global registry instance (set during app initialization)
var globalRegistry *Registry

// SetRegistry initializes the global component registry
func SetRegistry(app AppHandler, board BoardHandler) {
	globalRegistry = &Registry{
		AppHandler:   app,
		BoardHandler: board,
	}
}

// Template functions that can be used in any templ file
func RenderApp(ctx context.Context) templ.Component {
	if globalRegistry != nil && globalRegistry.AppHandler != nil {
		return globalRegistry.AppHandler.RenderComponent(ctx)
	}
	return nil
}

func RenderBoard(ctx context.Context) templ.Component {
	if globalRegistry != nil && globalRegistry.BoardHandler != nil {
		return globalRegistry.BoardHandler.RenderComponent(ctx)
	}
	return nil
}

// Alternative: Context-based approach (no globals)
type contextKey string

const registryKey contextKey = "components.registry"

// WithRegistry adds the registry to context
func WithRegistry(ctx context.Context, registry *Registry) context.Context {
	return context.WithValue(ctx, registryKey, registry)
}

// GetRegistry retrieves the registry from context
func GetRegistry(ctx context.Context) *Registry {
	if registry, ok := ctx.Value(registryKey).(*Registry); ok {
		return registry
	}
	return nil
}

// Context-based template functions (preferred approach)
func RenderAppFromContext(ctx context.Context) templ.Component {
	if registry := GetRegistry(ctx); registry != nil && registry.AppHandler != nil {
		return registry.AppHandler.RenderComponent(ctx)
	}
	return nil
}

func RenderBoardFromContext(ctx context.Context) templ.Component {
	if registry := GetRegistry(ctx); registry != nil && registry.BoardHandler != nil {
		return registry.BoardHandler.RenderComponent(ctx)
	}
	return nil
}
