package store

import (
	"sync"
	"testing"
	"time"

	"github.com/colearendt/cli-conn/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStore_AppendAndResults(t *testing.T) {
	t.Parallel()
	s := New()
	r1 := model.Result{Timestamp: time.Now(), Target: "a"}
	r2 := model.Result{Timestamp: time.Now(), Target: "b"}

	s.Append(r1)
	s.Append(r2)

	results := s.Results()
	assert.Len(t, results, 2)
	assert.Equal(t, "a", results[0].Target)
	assert.Equal(t, "b", results[1].Target)
}

func TestStore_Concurrent(t *testing.T) {
	t.Parallel()
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Append(model.Result{Timestamp: time.Now()})
		}()
	}
	wg.Wait()

	assert.Len(t, s.Results(), 100)
}
