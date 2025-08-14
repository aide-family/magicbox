// Package hello is a simple package that prints "Hello, World!"
package hello

import (
	"os"
	"sync"
)

type env struct {
	name     string
	version  string
	metadata map[string]string
	id       string
	env      string
}

var (
	once  sync.Once
	id, _ = os.Hostname()
	e     = &env{
		name:     "Moon",
		version:  "latest",
		metadata: make(map[string]string),
		env:      "dev",
		id:       id,
	}
)

type Option func(*env)

func WithName(name string) Option {
	return func(e *env) {
		e.name = name
	}
}

func WithVersion(version string) Option {
	return func(e *env) {
		e.version = version
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(e *env) {
		e.metadata = metadata
	}
}

func WithEnv(envType string) Option {
	return func(e *env) {
		e.env = envType
	}
}

func WithID(id string) Option {
	return func(e *env) {
		e.id = id
	}
}

func SetEnvWithOption(opts ...Option) {
	once.Do(func() {
		for _, opt := range opts {
			opt(e)
		}
	})
}

func Env() string {
	return e.env
}

func Name() string {
	return e.name
}

func Version() string {
	return e.version
}

func Metadata() map[string]string {
	return e.metadata
}

func ID() string {
	return e.id
}
