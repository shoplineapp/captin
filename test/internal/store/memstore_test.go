package store_test

import (
	"fmt"
	"testing"
	"time"

	store "github.com/shoplineapp/captin/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestStoreWithTTL(t *testing.T) {
	ms := store.NewMemoryStore()
	for i := 0; i < 10000; i++ {
		k, v := fmt.Sprint("key", i), fmt.Sprint("value", i)
		result, err := ms.Set(k, v, 200*time.Millisecond)
		assert.Nil(t, err)
		assert.True(t, result)
	}
	assert.Equal(t, 10000, ms.Len())
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, 0, ms.Len())
}
