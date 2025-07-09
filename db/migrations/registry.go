package migrations

import (
	"sort"

	"gorm.io/gorm"
)

// Migration struct
type Migration struct {
	Version string
	Up      func(*gorm.DB) error
}

// Migration registry
var registry = make(map[string]Migration)

// Register migration
func Register(m Migration) {
	registry[m.Version] = m
}

// List migrations
func List() []Migration {
	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := make([]Migration, 0, len(keys))
	for _, k := range keys {
		result = append(result, registry[k])
	}

	return result
}
