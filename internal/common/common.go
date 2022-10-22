package common

import (
	"log"
	"sync"
)

func LogFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type SyncMap[K comparable, V any] struct {
	mut *sync.Mutex
	m   map[K]V
}

func NewSyncMap[K comparable, V any](mutex *sync.Mutex) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		mut: mutex,
		m:   make(map[K]V),
	}
}

func (m *SyncMap[K, V]) Get(key K) (V, bool) {
	m.mut.Lock()
	defer m.mut.Unlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *SyncMap[K, V]) MustGet(key K) V {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.m[key]
}

func (m *SyncMap[K, V]) Set(key K, val V) {
	m.mut.Lock()
	defer m.mut.Unlock()
	m.m[key] = val
}
