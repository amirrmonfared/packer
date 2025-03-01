package store

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStoreConcurrency(t *testing.T) {
	s := NewMemoryStore([]int{250, 500})
	ctx := context.Background()

	var wg sync.WaitGroup
	errCh := make(chan error, 100)

	// 10 goroutines updating pack sizes concurrently
	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			newSizes := []int{250 + i, 500 + i}
			err := s.UpdatePackSizes(ctx, newSizes)
			if err != nil {
				errCh <- err
			}
		}(i)
	}

	// 10 goroutines reading pack sizes concurrently
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.GetPackSizes(ctx)
			if err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for e := range errCh {
		t.Errorf("unexpected error: %v", e)
	}

	// we just ensure we don't panic or race.
	packs, err := s.GetPackSizes(ctx)
	assert.NoError(t, err)
	assert.Len(t, packs, 2)
}
