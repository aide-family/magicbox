package cache

import (
	"fmt"
	"strings"
)

// K is a key for the cache.
type K string

// NewKey returns a new key for the cache.
func NewKey(s ...any) K {
	ks := make([]string, 0, len(s))
	for _, a := range s {
		ks = append(ks, fmt.Sprintf("%v", a))
	}
	return K(strings.Join(ks, separator))
}

// Key returns a key for the cache.
func (k K) Key(s ...any) string {
	ks := make([]string, 0, len(s))
	for _, a := range s {
		ks = append(ks, fmt.Sprintf("%v", a))
	}
	return strings.Join(append([]string{prefix, string(k)}, ks...), separator)
}

var (
	separator = ":"
	prefix    = ""
)

// SetSeparator sets the separator for the cache.
func SetSeparator(s string) {
	separator = s
}

// SetPrefix sets the prefix for the cache.
func SetPrefix(p string) {
	prefix = p
}
