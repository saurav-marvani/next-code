package linters

import (
	"context"
	"fmt"
	"sync"
)

// Registry manages all available linters
type Registry struct {
	linters map[string]Linter
	mu      sync.RWMutex
}

// NewRegistry creates a new linter registry
func NewRegistry() *Registry {
	return &Registry{
		linters: make(map[string]Linter),
	}
}

// Register adds a linter to the registry
func (r *Registry) Register(linter Linter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := linter.Name()
	if _, exists := r.linters[name]; exists {
		return fmt.Errorf("linter %s already registered", name)
	}

	r.linters[name] = linter
	return nil
}

// Get retrieves a linter by name
func (r *Registry) Get(name string) (Linter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	linter, exists := r.linters[name]
	if !exists {
		return nil, fmt.Errorf("linter %s not found", name)
	}
	return linter, nil
}

// GetAvailable returns all available (installed) linters
func (r *Registry) GetAvailable(ctx context.Context) ([]Linter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []Linter
	for _, linter := range r.linters {
		if linter.IsAvailable(ctx) {
			available = append(available, linter)
		}
	}
	return available, nil
}

// GetForLanguage returns all available linters that support a given language
func (r *Registry) GetForLanguage(ctx context.Context, language string) ([]Linter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matching []Linter
	for _, linter := range r.linters {
		if !linter.IsAvailable(ctx) {
			continue
		}

		for _, lang := range linter.SupportedLanguages() {
			if lang == language {
				matching = append(matching, linter)
				break
			}
		}
	}
	return matching, nil
}

// List returns all registered linters with their availability status
func (r *Registry) List(ctx context.Context) map[string]bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]bool)
	for name, linter := range r.linters {
		result[name] = linter.IsAvailable(ctx)
	}
	return result
}

// ListAll returns all registered linters
func (r *Registry) ListAll() []Linter {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []Linter
	for _, linter := range r.linters {
		all = append(all, linter)
	}
	return all
}
