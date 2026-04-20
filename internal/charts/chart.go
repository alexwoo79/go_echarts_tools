// Package charts defines the ChartBuilder interface and a global registry.
// New chart types register themselves via Register().
package charts

import "gantt/internal/model"

// ChartBuilder is the contract every chart type must implement.
type ChartBuilder interface {
	// ID returns the unique key for this chart type (e.g. "gantt").
	ID() string
	// Name returns the human-readable display name.
	Name() string
	// InferDefaults suggests a MappingConfig based on header names.
	InferDefaults(headers []string) model.MappingConfig
	// DefaultOptions returns sensible render options for this chart type.
	DefaultOptions() model.ChartOptions
	// Build transforms a Dataset + config into JSON-serialisable chart data.
	Build(dataset model.Dataset, cfg model.MappingConfig, opts model.ChartOptions) (interface{}, error)
}

var registry = map[string]ChartBuilder{}

// Register adds a ChartBuilder to the global registry.
func Register(b ChartBuilder) {
	registry[b.ID()] = b
}

// Get retrieves a ChartBuilder by ID.
func Get(id string) (ChartBuilder, bool) {
	v, ok := registry[id]
	return v, ok
}

// All returns all registered chart builders.
func All() []ChartBuilder {
	out := make([]ChartBuilder, 0, len(registry))
	for _, v := range registry {
		out = append(out, v)
	}
	return out
}
