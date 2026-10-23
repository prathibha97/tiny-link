package shortener

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestURLStoreSaveAndGet(t *testing.T) {
	store := NewURLStore()

	want := URLMapping{
		Code:        "go123",
		OriginalURL: "https://go.dev",
	}

	if err := store.Save(context.Background(), want); err != nil {
		t.Fatal(err)
	}

	mapping, err := store.Get(context.Background(), "go123")
	if err != nil {
		t.Fatal(err)
	}

	if mapping != want {
		t.Fatalf("got %v, want %v", mapping, want)
	}
}

func TestURLStoreGetMissing(t *testing.T) {
	store := NewURLStore()

	_, err := store.Get(context.Background(), "missing")

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
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

			store.Save(context.Background(), URLMapping{
				Code:        code,
				OriginalURL: "https://go.dev",
			})
			store.Get(context.Background(), code)
		}(i)
	}
	wg.Wait()
}
