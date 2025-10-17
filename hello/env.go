// Package hello is a simple package that prints "Hello, World!"
package hello

import (
	"os"
	"sync"

	"github.com/aide-family/magicbox/strutil"
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
		name:    "Moon",
		version: "latest",
		metadata: map[string]string{
			"author": "Aide Family",
			"email":  "1058165620@qq.com",
		},
		env: "PREVIEW",
		id:  id,
	}
)

type Option func(*env)

func WithName(name string) Option {
	return func(e *env) {
		if strutil.IsNotEmpty(name) {
			e.name = name
		}
	}
}

func WithVersion(version string) Option {
	return func(e *env) {
		if strutil.IsNotEmpty(version) {
			e.version = version
		}
	}
}

func WithMetadata(metadata map[string]string) Option {
	return func(e *env) {
		if len(metadata) > 0 {
			e.metadata = metadata
		}
	}
}

func WithEnv(envType string) Option {
	return func(e *env) {
		if strutil.IsNotEmpty(envType) {
			e.env = envType
		}
	}
}

func WithID(id string) Option {
	return func(e *env) {
		if strutil.IsNotEmpty(id) {
			e.id = id
		}
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
