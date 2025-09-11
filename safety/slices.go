package safety

import (
	"encoding"
	"slices"
	"sync"

	"github.com/aide-family/magicbox/serialize"
)

var _ encoding.BinaryMarshaler = (*Slice[any])(nil)
var _ encoding.BinaryUnmarshaler = (*Slice[any])(nil)

type Slice[T any] struct {
	mu sync.RWMutex
	s  []T
}

func NewSlice[T any](s []T) *Slice[T] {
	return &Slice[T]{s: slices.Clone(s)}
}

func (s *Slice[T]) Get(i int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s[i]
}

func (s *Slice[T]) Set(i int, v T) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s[i] = v
	return s
}

func (s *Slice[T]) Append(v T) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = append(s.s, v)
	return s
}

func (s *Slice[T]) AppendSlice(ss ...[]T) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	total := 0
	for _, s2 := range ss {
		total += len(s2)
	}
	s.s = slices.Grow(s.s, total)

	for _, s2 := range ss {
		s.s = append(s.s, s2...)
	}
	return s
}

func (s *Slice[T]) Delete(i int) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = slices.Delete(s.s, i, 1)
	return s
}

func (s *Slice[T]) DeleteFunc(f func(v T) bool) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = slices.DeleteFunc(s.s, f)
	return s
}

func (s *Slice[T]) Range(f func(v T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.s {
		if !f(v) {
			break
		}
	}
}

func (s *Slice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.s)
}

func (s *Slice[T]) Clone() *Slice[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return NewSlice(slices.Clone(s.s))
}

func (s *Slice[T]) Clear() *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = nil
	return s
}

func (s *Slice[T]) List() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return slices.Clone(s.s)
}

func (s *Slice[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bs, _ := serialize.JSONMarshal(s.s)
	return string(bs)
}

func (s *Slice[T]) MarshalBinary() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return serialize.JSONMarshal(s.s)
}

func (s *Slice[T]) UnmarshalBinary(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return serialize.JSONUnmarshal(data, &s.s)
}
