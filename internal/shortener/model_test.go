package shortener

import (
	"fmt"
	"sync"
	"testing"
)

func TestURLStoreSaveAndGet(t *testing.T) {
	store := NewURLStore()

	want := URLMapping{
		short_code:   "go123",
		original_url: "https://go.dev",
	}

	store.Save(want)

	got, ok := store.Get("go123")
	if !ok {
		t.Fatal("expected mapping to exist")
	}

	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestURLStoreGetMissing(t *testing.T) {
	store := NewURLStore()

	_, ok := store.Get("missing")
	if ok {
		t.Fatal("expected missing mapping to return false")
	}
}

func TestURLStoreConcurrentAccess(t *testing.T) {
	store := NewURLStore()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			code := fmt.Sprintf("code-%d", i)

			store.Save(URLMapping{
				short_code:   code,
				original_url: "https://go.dev",
			})
			store.Get(code)
		}(i)
	}
	wg.Wait()
}
