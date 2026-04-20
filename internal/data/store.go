// Package data handles dataset storage and file parsing.
package data

import (
	"strconv"
	"sync"
	"time"

	"gantt/internal/model"
)

var store = struct {
	sync.RWMutex
	m map[string]model.Dataset
}{m: make(map[string]model.Dataset)}

// Store saves a dataset by its ID.
func Store(dataset model.Dataset) {
	store.Lock()
	store.m[dataset.ID] = dataset
	store.Unlock()
}

// Load retrieves a dataset by ID.
func Load(id string) (model.Dataset, bool) {
	store.RLock()
	v, ok := store.m[id]
	store.RUnlock()
	return v, ok
}

// NewID generates a unique dataset ID based on nanosecond timestamp.
func NewID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
