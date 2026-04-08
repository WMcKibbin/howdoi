package provider

import "fmt"

var registry = map[string]Provider{}

// Register adds a provider to the global registry.
func Register(p Provider) {
	registry[p.Name()] = p
}

// Get returns a provider by name.
func Get(name string) (Provider, error) {
	p, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %q (available: %v)", name, Names())
	}
	return p, nil
}

// Names returns the names of all registered providers.
func Names() []string {
	names := make([]string, 0, len(registry))
	for n := range registry {
		names = append(names, n)
	}
	return names
}
