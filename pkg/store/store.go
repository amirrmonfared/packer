package store

import (
	"context"
	"sync"
)

type PackStore interface {
	UpdatePackSizes(ctx context.Context, sizes []int) error
	GetPackSizes(ctx context.Context) ([]int, error)
}

type MemoryStore struct {
	mu    sync.RWMutex
	packs []int
}

func NewMemoryStore(defaultPacks []int) *MemoryStore {
	cp := make([]int, len(defaultPacks))
	copy(cp, defaultPacks)
	return &MemoryStore{packs: cp}
}

func (m *MemoryStore) UpdatePackSizes(ctx context.Context, sizes []int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.packs = make([]int, len(sizes))
	copy(m.packs, sizes)
	return nil
}

func (m *MemoryStore) GetPackSizes(ctx context.Context) ([]int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cp := make([]int, len(m.packs))
	copy(cp, m.packs)
	return cp, nil
}
